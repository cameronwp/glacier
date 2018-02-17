package pool

import "testing"

func TestDrain(t *testing.T) {
	t.Parallel()

	t.Run("when jobs are active", func(st *testing.T) {
		st.Parallel()

		testCases := []struct {
			description string
			max         int
			chunk       Chunk
			uploadError error
		}{
			{
				description: "attempts to upload when it's below the max connections",
			},
		}

		for _, tc := range testCases {
			uploader := func(c Chunk) error {
				return tc.uploadError
			}
			testPool := New(uploader, tc.max)
			testPool.run = false
		}
	})

	t.Run("when no jobs are active", func(st *testing.T) {
		//
	})
}
