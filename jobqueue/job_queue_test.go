package jobqueue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

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

func TestIncrAttempts(t *testing.T) {
	j := Job{}

	assert.Equal(t, 0, j.numAttempts, "starting number of attempts")

	j.IncrAttempts()

	assert.Equal(t, 1, j.numAttempts, "ending number of attempts")
}

func TestAtMaxAttempts(t *testing.T) {
	j := Job{}

	for i := 0; i < (MaxJobAttempts - 1); i++ {
		j.IncrAttempts()
	}

	assert.False(t, j.AtMaxAttempts(), "max attempts under limit")

	j.IncrAttempts()

	assert.True(t, j.AtMaxAttempts(), "max attempts at limit")

	j.IncrAttempts()

	assert.True(t, j.AtMaxAttempts(), "max attempts over limit")
}

func TestAddJob(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "barfs when a chunk doesn't have an ID",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{}
				c1 := Chunk{}
				_, err := testJobQueue.AddJob(c1)

				return err
			},
			expectedError: ErrInvalidChunk,
		},
		{
			description: "creates a waiting job from a chunk",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{}
				c1 := Chunk{ID: "asdf1234"}
				_, err := testJobQueue.AddJob(c1)
				if err != nil {
					return err
				}

				assert.Equal(st, 1, len(testJobQueue.waitingJobs), "# waiting jobs")
				assert.Equal(st, c1.ID, testJobQueue.waitingJobs[0].Status.Chunk.ID, "chunk ID in queue")
				assert.Equal(st, Waiting, testJobQueue.waitingJobs[0].Status.State, "job state")

				return nil
			},
		},
		{
			description: "returns the right number of waiting jobs",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{}
				c1 := Chunk{ID: "asdf1234"}
				c2 := Chunk{ID: "1234asdf"}
				_, err := testJobQueue.AddJob(c1)
				if err != nil {
					return err
				}

				numJobs, err := testJobQueue.AddJob(c2)
				if err != nil {
					return err
				}

				assert.Equal(st, 2, numJobs, "# waiting jobs")

				return nil
			},
		},
		{
			description: "keeps adding new waiting jobs",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{}
				c1 := Chunk{ID: "asdf1234"}
				c2 := Chunk{ID: "1234asdf"}
				_, err := testJobQueue.AddJob(c1)
				if err != nil {
					return err
				}
				_, err = testJobQueue.AddJob(c2)
				if err != nil {
					return err
				}

				assert.Equal(st, 2, len(testJobQueue.waitingJobs), "# waiting jobs")
				assert.Equal(st, c1.ID, testJobQueue.waitingJobs[0].Status.Chunk.ID, "chunk ID in queue")
				assert.Equal(st, c2.ID, testJobQueue.waitingJobs[1].Status.Chunk.ID, "chunk ID in queue")

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

func TestActivateOldestWaitingJob(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "moves a waiting job to active when under the max",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{
					MaxJobs: 1,
				}
				j := Job{}
				testJobQueue.waitingJobs = append(testJobQueue.waitingJobs, &j)

				assert.Equal(st, 1, len(testJobQueue.waitingJobs), "# waiting jobs")

				_, err := testJobQueue.ActivateOldestWaitingJob()
				if err != nil {
					return err
				}

				assert.Equal(st, 0, len(testJobQueue.waitingJobs), "# waiting jobs")
				assert.Equal(st, 1, len(testJobQueue.activeJobs), "# active jobs")

				return nil
			},
		},
		{
			description: "moves the oldest waiting job to active",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{
					MaxJobs: 1,
				}
				c1 := Chunk{ID: "asdf1234"}
				c2 := Chunk{ID: "1234asdf"}
				_, err := testJobQueue.AddJob(c1)
				if err != nil {
					return err
				}
				_, err = testJobQueue.AddJob(c2)
				if err != nil {
					return err
				}

				_, err = testJobQueue.ActivateOldestWaitingJob()
				if err != nil {
					return err
				}

				assert.Equal(st, c1.ID, testJobQueue.activeJobs[0].Status.Chunk.ID, "chunk IDs")

				return nil
			},
		},
		{
			description: "errs when at the max # of jobs",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{
					MaxJobs: 1,
				}
				j := Job{}
				testJobQueue.waitingJobs = append(testJobQueue.waitingJobs, &j)
				testJobQueue.waitingJobs = append(testJobQueue.waitingJobs, &j)

				_, err := testJobQueue.ActivateOldestWaitingJob()
				if err != nil {
					return err
				}

				assert.Equal(st, 1, len(testJobQueue.activeJobs), "# active jobs")

				_, err = testJobQueue.ActivateOldestWaitingJob()
				return err
			},
			expectedError: ErrMaxActiveJobs,
		},
		{
			description: "errs when over the max # of jobs",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{
					MaxJobs: 1,
				}
				j := Job{}
				testJobQueue.waitingJobs = append(testJobQueue.waitingJobs, &j)
				testJobQueue.activeJobs = append(testJobQueue.activeJobs, &j)
				testJobQueue.activeJobs = append(testJobQueue.activeJobs, &j)

				_, err := testJobQueue.ActivateOldestWaitingJob()
				if err != nil {
					return err
				}

				assert.Equal(st, 1, len(testJobQueue.activeJobs), "# active jobs")

				_, err = testJobQueue.ActivateOldestWaitingJob()
				return err
			},
			expectedError: ErrMaxActiveJobs,
		},
		{
			description: "errs if no jobs are waiting",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{
					MaxJobs: 1,
				}
				_, err := testJobQueue.ActivateOldestWaitingJob()

				return err
			},
			expectedError: ErrNoWaitingJobs,
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

func TestNext(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "returns the oldest waiting job from the active queue",
			test: func(st *testing.T) error {
				testJobQueue, err := fillActiveQueue(2)
				if err != nil {
					return err
				}

				j, err := testJobQueue.Next()
				if err != nil {
					return err
				}

				assert.Equal(st, j.Status.Chunk.ID, testJobQueue.activeJobs[0].Status.Chunk.ID, "chunk IDs match")

				return nil
			},
		},
		{
			description: "skips over in progress jobs in the active queue",
			test: func(st *testing.T) error {
				testJobQueue, err := fillActiveQueue(2)
				if err != nil {
					return err
				}

				testJobQueue.activeJobs[0].Status.State = InProgress

				j, err := testJobQueue.Next()
				if err != nil {
					return err
				}

				assert.Equal(st, j.Status.Chunk.ID, testJobQueue.activeJobs[1].Status.Chunk.ID, "chunk IDs match")

				return nil
			},
		},
		{
			description: "errs if all jobs are in progress",
			test: func(st *testing.T) error {
				testJobQueue, err := fillActiveQueue(2)
				if err != nil {
					return err
				}

				testJobQueue.activeJobs[0].Status.State = InProgress
				testJobQueue.activeJobs[1].Status.State = InProgress

				j, err := testJobQueue.Next()

				assert.Nil(st, j)
				return err
			},
			expectedError: ErrAllActiveJobsInProgress,
		},
		{
			description: "errs if no jobs are active",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{}

				_, err := testJobQueue.Next()
				return err
			},
			expectedError: ErrNoActiveJobs,
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

func TestCompleteJob(t *testing.T) {
	testCases := []struct {
		description   string
		test          func(*testing.T) error
		expectedError error
	}{
		{
			description: "moves an active job to completed jobs",
			test: func(st *testing.T) error {
				testJobQueue, err := fillActiveQueue(2)
				if err != nil {
					return err
				}

				j, err := testJobQueue.Next()
				if err != nil {
					return err
				}

				_, err = testJobQueue.CompleteJob(j)

				assert.Equal(st, 1, len(testJobQueue.completedJobs), "# completed jobs")

				return err
			},
		},
		{
			description: "moves the right active job to completed jobs when multiple are in progress",
			test: func(st *testing.T) error {
				testJobQueue, err := fillActiveQueue(2)
				if err != nil {
					return err
				}

				j1, err := testJobQueue.Next()
				if err != nil {
					return err
				}

				j1.Status.State = InProgress

				j2, err := testJobQueue.Next()
				if err != nil {
					return err
				}

				j2.Status.State = InProgress

				_, err = testJobQueue.CompleteJob(j2)

				assert.Equal(st, j2.Status.Chunk.ID, testJobQueue.completedJobs[0].Status.Chunk.ID, "completed job ID")

				return err
			},
		},
		{
			description: "reports an accurate number of completed jobs",
			test: func(st *testing.T) error {
				testJobQueue, err := fillActiveQueue(2)
				if err != nil {
					return err
				}

				j, err := testJobQueue.Next()
				if err != nil {
					return err
				}

				numCompleted, err := testJobQueue.CompleteJob(j)

				assert.Equal(st, 1, numCompleted, "# completed jobs")

				return err
			},
		},
		{
			description: "barfs if the job can't be found",
			test: func(st *testing.T) error {
				testJobQueue := JobQueue{}

				_, err := testJobQueue.CompleteJob(&Job{
					Status: Status{
						Chunk: &Chunk{
							ID: "qwertyuiop",
						},
					},
				})

				return err
			},
			expectedError: fmt.Errorf("job with chunk ID 'qwertyuiop' is not active"),
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
