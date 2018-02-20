package fs

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/cameronwp/glacier/jobqueue"
)

// DefaultPartSize is the default size of upload parts.
// const DefaultPartSize = int64(1 << 23) // 8MB
const DefaultPartSize = int64(1 << 20) // 1MB

// Chunker can create file chunks.
type Chunker interface {
	GetFilesize(string) (int64, error)
	CreateChunks(string, int64, int64) ([]jobqueue.Chunk, error)
}

// OSChunker implements Chunker using os.
type OSChunker struct{}

var _ Chunker = (*OSChunker)(nil)

// rip up files, watch how long they take, change part sizes, # pool connections

// GetFilepaths returns all filepaths within a dir recursively. It is safe to
// pass a path to a file instead of a dir.
func GetFilepaths(fp string) ([]string, error) {
	aggregator := []string{}

	maybeFileOrDirectory, err := os.Stat(fp)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			// no file exists
			return nil, nil
		default:
			return nil, err
		}
	}

	// just return the file if it isn't a directory
	if !maybeFileOrDirectory.IsDir() {
		aggregator = append(aggregator, fp)
		return aggregator, nil
	}

	// recurse over the directory until files are found
	err = filepath.Walk(fp, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// don't run against the root of the dir again
		if fp == path {
			return nil
		}

		if !info.IsDir() {
			aggregator = append(aggregator, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return aggregator, err
}

// ChunkFile creates chunks from a file. You must already have an upload ID and
// a known part size. Chunk sizes may be determined on the fly.
func ChunkFile(fs Chunker, filepath string, partsize int64) ([]jobqueue.Chunk, error) {
	totalSize, err := fs.GetFilesize(filepath)
	if err != nil {
		return nil, err
	}

	return fs.CreateChunks(filepath, partsize, totalSize)
}

// GetFilesize returns the size of a file.
func (*OSChunker) GetFilesize(filepath string) (int64, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	stats, err := f.Stat()
	if err != nil {
		return 0, nil
	}

	return stats.Size(), err
}

// CreateChunks does the math to determine start byte and end byte for each
// chunk.
func (*OSChunker) CreateChunks(filepath string, partsize int64, totalSize int64) ([]jobqueue.Chunk, error) {
	var aggregator []jobqueue.Chunk

	startB := int64(0)

	for {
		contentLength := int(math.Min(float64(partsize), float64(totalSize-startB)))
		if contentLength == 0 {
			// at the end of the file
			break
		}

		endB := startB + int64(contentLength)

		// note that there is no UploadID. it gets added later
		aggregator = append(aggregator, jobqueue.Chunk{
			Path:   filepath,
			StartB: startB,
			EndB:   endB,
		})

		startB = endB
	}

	return aggregator, nil
}
