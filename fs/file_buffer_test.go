package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

func TestSortHashes(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "when there are two out of order hashes",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 10,
						endB:   20,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
				}

				sortHashes(hashes)

				assert.Equal(st, hashes[0].startB, int64(0), "first hash is first")
				assert.Equal(st, hashes[1].startB, int64(10), "second hash is second")

				return nil
			},
		},
		{
			description: "when there are a bunch of out of order hashes",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 20,
						endB:   30,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 10,
						endB:   20,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 30,
						endB:   40,
					},
				}

				sortHashes(hashes)

				assert.Equal(st, hashes[0].startB, int64(0), "first hash is first")
				assert.Equal(st, hashes[3].startB, int64(30), "fourth hash is fourth")

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

func TestIsCompleteFile(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "if there are no hashes",
			test: func(st *testing.T) error {
				hashes := []FileHash{}

				assert.False(st, isCompleteFile(int64(10), hashes))
				return nil
			},
		},
		{
			description: "when a file has one complete hash",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
				}

				assert.True(st, isCompleteFile(int64(10), hashes))
				return nil
			},
		},
		{
			description: "when a file has one incomplete hash",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
				}

				assert.False(st, isCompleteFile(int64(11), hashes))
				return nil
			},
		},
		{
			description: "when a file has two complete hashes in order",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 10,
						endB:   20,
					},
				}

				assert.True(st, isCompleteFile(int64(20), hashes))
				return nil
			},
		},
		{
			description: "when a file has two complete hashes out of order",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 10,
						endB:   20,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
				}

				assert.False(st, isCompleteFile(int64(20), hashes))
				return nil
			},
		},
		{
			description: "when a file has three complete hashes in order",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 10,
						endB:   20,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 20,
						endB:   27,
					},
				}

				assert.True(st, isCompleteFile(int64(27), hashes))
				return nil
			},
		},
		{
			description: "when a file is missing the last hash",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 10,
						endB:   20,
					},
				}

				assert.False(st, isCompleteFile(int64(27), hashes))
				return nil
			},
		},
		{
			description: "when a file is missing the first hash",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 10,
						endB:   20,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 20,
						endB:   27,
					},
				}

				assert.False(st, isCompleteFile(int64(27), hashes))
				return nil
			},
		},
		{
			description: "when a file has three complete hashes out of order",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 20,
						endB:   27,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 10,
						endB:   20,
					},
				}

				assert.False(st, isCompleteFile(int64(27), hashes))
				return nil
			},
		},
		{
			description: "when a file is missing a hash in the middle",
			test: func(st *testing.T) error {
				hashes := []FileHash{
					{
						sha256: randstr.RandomString(8),
						startB: 0,
						endB:   10,
					},
					{
						sha256: randstr.RandomString(8),
						startB: 20,
						endB:   27,
					},
				}

				assert.False(st, isCompleteFile(int64(27), hashes))
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
