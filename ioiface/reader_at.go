package ioiface

// ReadAt mimics the io.ReaderAt interface for mocking.
type ReadAt interface {
	ReadAt([]byte, int64) (int, error)
}
