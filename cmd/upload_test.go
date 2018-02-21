package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDescription(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "returns the filename when target and path are the same",
			test: func(st *testing.T) error {
				target := "../fs/testdata/f0"
				path, err := filepath.Abs(target)
				if err != nil {
					return err
				}

				desc, err := GetDescription(target, path)

				assert.Equal(st, "f0", desc, "description should be the filename")
				return err
			},
		},
		{
			description: "returns a relative path from target to filename when target is a dir",
			test: func(st *testing.T) error {
				target := "../fs/testdata/"
				absPath, err := filepath.Abs(target + "/f0")
				if err != nil {
					return err
				}

				desc, err := GetDescription(target, absPath)

				assert.Equal(st, "f0", desc, "description should be the filename")
				return err
			},
		},
		{
			description: "returns a relative path with intermediate dirs when target is a dir",
			test: func(st *testing.T) error {
				target := "../fs/testdata/"
				absPath, err := filepath.Abs(target + "/dir1/subdir1/f3")
				if err != nil {
					return err
				}

				desc, err := GetDescription(target, absPath)

				assert.Equal(st, "dir1/subdir1/f3", desc, "description should be relative path to the filename")
				return err
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.description, func(st *testing.T) {
			err := tc.test(st)
			if tc.expectedError != nil {
				assert.Equal(st, tc.expectedError, err)
			} else {
				assert.NoError(st, err)
			}
		})
	}
}
