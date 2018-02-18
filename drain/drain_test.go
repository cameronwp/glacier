package drain

import (
	"fmt"
	"testing"

	"github.com/cameronwp/glacier/jobqueue"
	"github.com/cameronwp/glacier/jobqueue/jobqueuemocks"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

func nilUploader(*jobqueue.Chunk) error {
	return nil
}

func failUploader(*jobqueue.Chunk) error {
	return fmt.Errorf("fail")
}

func randomJob() *jobqueue.Job {
	return &jobqueue.Job{
		Status: jobqueue.Status{
			Chunk: &jobqueue.Chunk{
				ID: randstr.Hex(8),
			},
		},
	}
}

func TestAttemptJob(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "when a job succeeds",
			test: func(st *testing.T) error {
				j := randomJob()

				attemptJob(nilUploader, j)

				assert.Equal(st, jobqueue.Completed, j.Status.State, "job state")

				return nil
			},
		},
		{
			description: "when a job fails the first time it goes back to waiting",
			test: func(st *testing.T) error {
				j := randomJob()

				attemptJob(failUploader, j)

				assert.Equal(st, jobqueue.Waiting, j.Status.State, "job state")

				return nil
			},
		},
		{
			description: "when a job fails the last time it gets marked as erred",
			test: func(st *testing.T) error {
				j := randomJob()

				for i := 0; i <= jobqueue.MaxJobAttempts; i++ {
					attemptJob(failUploader, j)
				}

				assert.Equal(st, jobqueue.Erred, j.Status.State, "job state")

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
					description: "first reports when the job is starting",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()

						mockJQ.On("Next").
							Return(j, nil)

						mockJQ.On("CompleteJob", j).
							Return(1, nil)

						mockJQ.On("ActivateOldestWaitingJob").
							Return(0, jobqueue.ErrNoWaitingJobs)

						d := NewDrain(nilUploader)
						go d.Drain(&mockJQ)

						val := <-d.Schan
						assert.Equal(st, jobqueue.InProgress, val.State, "chunk state")

						return nil
					},
				},
				{
					description: "finally reports when the job ended",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()

						mockJQ.On("Next").
							Return(j, nil)

						mockJQ.On("CompleteJob", j).
							Return(1, nil)

						mockJQ.On("ActivateOldestWaitingJob").
							Return(0, jobqueue.ErrNoWaitingJobs)

						d := NewDrain(nilUploader)
						go d.Drain(&mockJQ)

						<-d.Schan
						val := <-d.Schan
						assert.Equal(st, jobqueue.Completed, val.State, "chunk state")

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
