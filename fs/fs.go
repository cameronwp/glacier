package fs

import (
	"os"
	"path/filepath"
)

const startingPartSize = int64(1 << 20) // 1MB

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

	// recurse over the directory until files are found
	if maybeFileOrDirectory.IsDir() {
		err := filepath.Walk(fp, func(path string, info os.FileInfo, err error) error {
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
	} else {
		aggregator = append(aggregator, fp)
	}

	return aggregator, err
}
