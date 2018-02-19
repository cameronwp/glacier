package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFilepaths(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "succeeds if passed a path to a file",
			test: func(st *testing.T) error {
				paths, err := GetFilepaths("./testdata/f0")

				assert.Len(st, paths, 1, "1 file")
				assert.Regexp(st, "f0", paths[0], "path to file")
				return err
			},
		},
		{
			description: "succeeds if passed a path to a dir",
			test: func(st *testing.T) error {
				paths, err := GetFilepaths("./testdata")

				assert.Len(st, paths, 4, "4 files")
				return err
			},
		},
		{
			description: "returns an empty slice if passed a non-existent path",
			test: func(st *testing.T) error {
				paths, err := GetFilepaths("./notarealpath/")

				assert.Len(st, paths, 0, "0 files")
				return err
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.description, func(st *testing.T) {
			err := tc.test(st)
			if tc.expectedError != nil {
				assert.Equal(st, err, tc.expectedError)
			} else {
				assert.NoError(st, err)
			}
		})
	}
}

func TestCreateChunks(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "succeeds when a file is larger than the partsize",
			test: func(st *testing.T) error {
				chunker := &OSChunker{}
				chunks, err := chunker.CreateChunks("asdf", "asdf", int64(10), int64(12))

				assert.Len(st, chunks, 2, "2 chunks")

				assert.Equal(st, int64(0), chunks[0].StartB, "first starting byte")
				assert.Equal(st, int64(10), chunks[0].EndB, "first ending byte")
				assert.Equal(st, int64(10), chunks[1].StartB, "second starting byte")
				assert.Equal(st, int64(12), chunks[1].EndB, "second ending byte")
				return err
			},
		},
		{
			description: "succeeds when a file is smaller than the partsize",
			test: func(st *testing.T) error {
				chunker := &OSChunker{}
				chunks, err := chunker.CreateChunks("asdf", "asdf", int64(12), int64(10))

				assert.Len(st, chunks, 1, "1 chunk")

				assert.Equal(st, int64(0), chunks[0].StartB, "first starting byte")
				assert.Equal(st, int64(10), chunks[0].EndB, "first ending byte")
				return err
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.description, func(st *testing.T) {
			err := tc.test(st)
			if tc.expectedError != nil {
				assert.Equal(st, err, tc.expectedError)
			} else {
				assert.NoError(st, err)
			}
		})
	}
}

func TestGetFilesize(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "succeeds when a file exists",
			test: func(st *testing.T) error {
				chunker := &OSChunker{}
				size, err := chunker.GetFilesize("./testdata/f0")

				assert.Equal(st, int64(32), size, "known filesize reported")

				return err
			},
		},
		{
			description: "errs when a file does not exist",
			test: func(st *testing.T) error {
				chunker := &OSChunker{}
				size, err := chunker.GetFilesize("./notarealfile")

				assert.Equal(st, int64(0), size, "0 size")

				pathError := &os.PathError{}
				assert.IsType(st, pathError, err, "PathError is returned")

				return nil
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.description, func(st *testing.T) {
			err := tc.test(st)
			if tc.expectedError != nil {
				assert.Equal(st, err, tc.expectedError)
			} else {
				assert.NoError(st, err)
			}
		})
	}
}
