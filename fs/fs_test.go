package fs

import (
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
