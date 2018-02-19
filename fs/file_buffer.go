package fs

// FileChunk is a literal chunk of a file alongside its sha256 hash.
type FileChunk struct {
	buf  []byte
	hash string
}

// memoize grabbing a buffer from a file
// garbage collection:
//   * after all of a file's chunks are done
//   * kill all of a file's buffers if one chunk fails
//   * if memory usage is too high

// GetFileBuffer returns a buffer of a chunk of a file alongside its 256 hash.
func GetFileBuffer(filepath string, startB int64, endB int64) ([]FileChunk, error) {
	return nil, nil
}
