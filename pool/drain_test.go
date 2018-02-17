package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

func nilUploader(*Chunk) error {
	return nil
}

// add numActive jobs to the active queue, each with a unique chunk ID
func fillActiveQueue(numActive int) (*JobQueue, error) {
	testJobQueue := JobQueue{
		MaxJobs: numActive,
	}

	for i := 0; i < numActive; i++ {
		c := Chunk{
			ID: randstr.Hex(8),
		}
		_, err := testJobQueue.AddJob(c)
		if err != nil {
			return &testJobQueue, err
		}
		_, err = testJobQueue.ActivateOldestWaitingJob()
		if err != nil {
			return &testJobQueue, err
		}
	}

	return &testJobQueue, nil
}

func TestDrain(t *testing.T) {
	type testCase struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}

	suites := []struct {
		description string
		testCases   []testCase
	}{
		{
			description: "when 1 job is active",
			testCases: []testCase{
				{
					description: "reports when the job is starting",
					test: func(st *testing.T) error {
						testJobQueue, err := fillActiveQueue(8)
						if err != nil {
							return err
						}

						d := NewDrain(nilUploader)
						go d.Drain(testJobQueue)
						val := <-d.Schan

						assert.Nil(st, val)
						return nil
					},
				},
			},
		},
	}

	t.Parallel()
	for _, s := range suites {
		t.Run(s.description, func(st *testing.T) {
			for _, tc := range s.testCases {
				err := tc.test(st)
				if tc.expectedError != nil {
					assert.Equal(st, err, tc.expectedError)
				} else {
					assert.NoError(st, err)
				}
			}
		})
	}
}
