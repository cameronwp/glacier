package fs

import (
	"fmt"
	"sort"
)

var (
	// ErrMissingFileChunks is returned when calling TreeHash against a file that
	// has not been fully buffered.
	ErrMissingFileChunks = fmt.Errorf("not all file chunks have been hashed")
)

// FileChunk is a literal chunk of a file alongside its sha256 hash.
type FileChunk struct {
	buf    []byte
	sha256 string
}

// FileHash represents a hash and the part of the file it represents.
type FileHash struct {
	sha256 string
	startB int64
	endB   int64
}

// BufferFetcher can buffer part of a file.
type BufferFetcher interface {
	FetchBuffer(string, int64, int64) (FileChunk, error)
	TreeHash(string) (string, error)
}

// OSBuffer can grab buffers from the local filesystem.
type OSBuffer struct {
	hashes map[string][]FileHash
}

var _ BufferFetcher = (*OSBuffer)(nil)

// FetchBuffer returns a buffer of a chunk of a file alongside its 256 hash.
func (*OSBuffer) FetchBuffer(filepath string, startB int64, endB int64) (FileChunk, error) {
	//
	return FileChunk{}, nil
}

// TreeHash returns the full hash for a file. Returns ErrMissingFileChunks if
// the whole file has not been buffered.
func (osb *OSBuffer) TreeHash(filepath string) (string, error) {
	chunker := OSChunker{}
	filesize, err := chunker.GetFilesize(filepath)
	if err != nil {
		return "", err
	}

	sortHashes(osb.hashes[filepath])

	if isCompleteFile(filesize, osb.hashes[filepath]) {
		// treehash!
	}
	return "", nil
}

func sortHashes(fileHash []FileHash) {
	sort.Slice(fileHash, func(i, j int) bool {
		return fileHash[i].startB < fileHash[j].startB
	})
}

// hash bytes should look like: [0, 10], [10, 20], [20, 30], etc
// the hashes must be sorted by starting bytes
func isCompleteFile(filesize int64, fileHashes []FileHash) bool {
	complete := false
	lastEndB := int64(0)

	numHashes := len(fileHashes)

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

		lastEndB = hash.endB
		complete = true
	}

	return complete
}
