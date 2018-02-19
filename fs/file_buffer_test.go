package fs

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/cameronwp/glacier/ioiface/ioifacemocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thanhpk/randstr"
)

func TestFetchBuffer(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "succeeds when reader covers the range of bytes",
			test: func(st *testing.T) error {
				filepath := "path/to/file"
				startB := int64(0)
				endB := int64(10)
				intLen := int(endB - startB)
				randomBuffer := randstr.RandomBytes(intLen)
				randomBufferHash := sha256.Sum256(randomBuffer)

				osb := OSBuffer{}
				var readerMock ioifacemocks.ReadAt

				readerMock.On("ReadAt",
					mock.Anything,
					startB,
				).Run(func(args mock.Arguments) {
					buf := args.Get(0).([]byte)
					copy(buf, randomBuffer)
				}).Return(intLen, nil)

				assert.Empty(st, osb, "no hashes to start")

				fc, err := osb.FetchAndHash(&readerMock, filepath, startB, endB)

				assert.Len(st, osb, 1, "has hashes for 1 file")
				assert.Equal(st, randomBuffer, fc.buf, "buffer is right")
				assert.Equal(st, startB, osb[filepath][0].startB, "right start bytes")
				assert.Equal(st, endB, osb[filepath][0].endB, "right end bytes")
				assert.Equal(st, randomBufferHash[:], fc.sha256, "hash is right in FileChunk")
				assert.Equal(st, randomBufferHash[:], osb[filepath][0].sha256, "hash is right in FileHash")
				return err
			},
		},
		{
			description: "errs with errors reading",
			test: func(st *testing.T) error {
				filepath := "path/to/file"
				startB := int64(0)
				endB := int64(10)
				intLen := int(endB - startB)

				osb := OSBuffer{}
				var readerMock ioifacemocks.ReadAt

				readerMock.On("ReadAt",
					mock.Anything,
					startB,
				).Return(intLen, fmt.Errorf("error"))

				_, err := osb.FetchAndHash(&readerMock, filepath, startB, endB)

				return err
			},
			expectedError: fmt.Errorf("error"),
		},
		{
			description: "errs if the wrong number of bytes were read",
			test: func(st *testing.T) error {
				filepath := "path/to/file"
				startB := int64(0)
				endB := int64(10)
				intLen := int(endB - startB)

				osb := OSBuffer{}
				var readerMock ioifacemocks.ReadAt

				readerMock.On("ReadAt",
					mock.Anything,
					startB,
				).Return(intLen-1, nil)

				_, err := osb.FetchAndHash(&readerMock, filepath, startB, endB)

				return err
			},
			expectedError: ErrIncompleteBuffer,
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

func TestSortHashes(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "when there are two out of order hashes",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 10,
						endB:   20,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
				}

				sortHashes(fileHashes)

				assert.Equal(st, fileHashes[0].startB, int64(0), "first hash is first")
				assert.Equal(st, fileHashes[1].startB, int64(10), "second hash is second")

				return nil
			},
		},
		{
			description: "when there are a bunch of out of order hashes",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 20,
						endB:   30,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 10,
						endB:   20,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 30,
						endB:   40,
					},
				}

				sortHashes(fileHashes)

				assert.Equal(st, fileHashes[0].startB, int64(0), "first hash is first")
				assert.Equal(st, fileHashes[3].startB, int64(30), "fourth hash is fourth")

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

func TestGetFileHashes(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "if there are no hashes",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{}

				_, ok := getFileHashes(int64(10), fileHashes)
				assert.False(st, ok)
				return nil
			},
		},
		{
			description: "when a file has one complete hash",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
				}

				hashes, ok := getFileHashes(int64(10), fileHashes)
				assert.True(st, ok)
				assert.Len(st, hashes, 1, "one hash")
				assert.Equal(st, hashes[0], fileHashes[0].sha256, "hashes are the same")
				return nil
			},
		},
		{
			description: "when a file has one incomplete hash",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
				}

				_, ok := getFileHashes(int64(11), fileHashes)
				assert.False(st, ok)
				return nil
			},
		},
		{
			description: "when a file has two complete hashes in order",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 10,
						endB:   20,
					},
				}

				hashes, ok := getFileHashes(int64(20), fileHashes)
				assert.Len(st, hashes, 2, "2 hashes")
				assert.Equal(st, hashes[0], fileHashes[0].sha256, "first hash is the same")
				assert.Equal(st, hashes[1], fileHashes[1].sha256, "second hash is the same")
				assert.True(st, ok)
				return nil
			},
		},
		{
			description: "when a file has two complete hashes out of order",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 10,
						endB:   20,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
				}

				_, ok := getFileHashes(int64(20), fileHashes)
				assert.False(st, ok)
				return nil
			},
		},
		{
			description: "when a file has three complete hashes in order",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 10,
						endB:   20,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 20,
						endB:   27,
					},
				}

				hashes, ok := getFileHashes(int64(27), fileHashes)
				assert.Len(st, hashes, 3, "3 hashes")
				assert.True(st, ok)
				return nil
			},
		},
		{
			description: "when a file is missing the last hash",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 10,
						endB:   20,
					},
				}

				_, ok := getFileHashes(int64(27), fileHashes)
				assert.False(st, ok)
				return nil
			},
		},
		{
			description: "when a file is missing the first hash",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 10,
						endB:   20,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 20,
						endB:   27,
					},
				}

				_, ok := getFileHashes(int64(27), fileHashes)
				assert.False(st, ok)
				return nil
			},
		},
		{
			description: "when a file has three complete hashes out of order",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 20,
						endB:   27,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 10,
						endB:   20,
					},
				}

				_, ok := getFileHashes(int64(27), fileHashes)
				assert.False(st, ok)
				return nil
			},
		},
		{
			description: "when a file is missing a hash in the middle",
			test: func(st *testing.T) error {
				fileHashes := []FileHash{
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 0,
						endB:   10,
					},
					{
						sha256: []byte(randstr.RandomString(8)),
						startB: 20,
						endB:   27,
					},
				}

				_, ok := getFileHashes(int64(27), fileHashes)
				assert.False(st, ok)
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
