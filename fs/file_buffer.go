package fs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"sort"

	"github.com/aws/aws-sdk-go/service/glacier"
)

var (
	// ErrMissingFileChunks is returned when calling TreeHash against a file that
	// has not been fully buffered.
	ErrMissingFileChunks = fmt.Errorf("not all file chunks have been hashed")

	// ErrIncompleteBuffer is returned when a ReadAt operation does not return the
	// contentlength (endB - startB) of bytes.
	ErrIncompleteBuffer = fmt.Errorf("could not read the whole buffer")
)

// FileChunk is a literal chunk of a file alongside its sha256 hash.
type FileChunk struct {
	buf    []byte
	sha256 []byte
}

// FileHash represents a hash and the part of the file it represents.
type FileHash struct {
	sha256 []byte
	startB int64
	endB   int64
}

// BufferFetcherHasher can buffer and hash parts of 1+ files.
type BufferFetcherHasher interface {
	FetchAndHash(io.ReaderAt, string, int64, int64) (FileChunk, error)
	TreeHash(Chunker, string) ([]byte, error)
}

// OSBuffer can grab buffers from the local filesystem.
type OSBuffer map[string][]FileHash

var _ BufferFetcherHasher = (*OSBuffer)(nil)

// FetchAndHash returns a buffer of a chunk of something that can be read
// alongside its sha256 hash. It should also save the hash for creating a tree
// hash of a whole readable (like a file).
func (osb OSBuffer) FetchAndHash(f io.ReaderAt, filepath string, startB int64, endB int64) (FileChunk, error) {
	buf := make([]byte, endB-startB)
	n, err := f.ReadAt(buf, startB)
	if err != nil {
		return FileChunk{}, err
	}
	if int64(n) != endB-startB {
		return FileChunk{}, ErrIncompleteBuffer
	}

	hash := sha256.Sum256(buf[:n])

	fileHash := FileHash{
		startB: startB,
		endB:   endB,
		sha256: hash[:],
	}
	fileChunk := FileChunk{
		buf:    buf,
		sha256: hash[:],
	}

	if _, ok := osb[filepath]; !ok {
		osb[filepath] = []FileHash{}
	}

	// save the hash for treehashing later
	osb[filepath] = append(osb[filepath], fileHash)

	return fileChunk, nil
}

// TreeHash returns the full hash for a file. Returns ErrMissingFileChunks if
// the whole file has not been buffered.
func (osb OSBuffer) TreeHash(chunker Chunker, filepath string) ([]byte, error) {
	filesize, err := chunker.GetFilesize(filepath)
	if err != nil {
		return nil, err
	}

	sortHashes(osb[filepath])

	if hashes, ok := getFileHashes(filesize, osb[filepath]); ok {
		return glacier.ComputeTreeHash(hashes), nil
	}
	return nil, ErrMissingFileChunks
}

func sortHashes(fileHash []FileHash) {
	sort.Slice(fileHash, func(i, j int) bool {
		return fileHash[i].startB < fileHash[j].startB
	})
}

// hash bytes should look like: [0, 10], [10, 20], [20, 30], etc
// the hashes must be sorted by starting bytes
func getFileHashes(filesize int64, fileHashes []FileHash) ([][]byte, bool) {
	complete := false
	lastEndB := int64(0)

	numHashes := len(fileHashes)

	hashes := [][]byte{}

	for i, hash := range fileHashes {
		if i == 0 {
			if hash.startB != 0 {
				// first hash is missing
				break
			}
		}

		if hash.startB != lastEndB {
			// something is out of order
			complete = false
			break
		}

		// on the last hash
		if i == numHashes-1 {
			if hash.endB != filesize {
				// last hash is missing
				complete = false
				break
			}
		}

		hashes = append(hashes, hash.sha256)
		lastEndB = hash.endB
		complete = true
	}

	return hashes, complete
}
