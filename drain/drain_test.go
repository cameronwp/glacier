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

func randomChunk() *jobqueue.Chunk {
	return &jobqueue.Chunk{
		ID: randstr.Hex(8),
	}
}

func randomJob() *jobqueue.Job {
	return &jobqueue.Job{
		Status: jobqueue.Status{
			Chunk: randomChunk(),
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

				assert.True(st, attemptJob(nilUploader, j))

				assert.Equal(st, jobqueue.Completed, j.Status.State, "job state")

				return nil
			},
		},
		{
			description: "when a job fails the first time it goes back to waiting",
			test: func(st *testing.T) error {
				j := randomJob()

				assert.False(st, attemptJob(failUploader, j))

				assert.Equal(st, jobqueue.Waiting, j.Status.State, "job state")

				return nil
			},
		},
		{
			description: "when a job fails the last time it gets marked as erred",
			test: func(st *testing.T) error {
				j := randomJob()

				for i := 0; i < jobqueue.MaxJobAttempts-1; i++ {
					assert.False(st, attemptJob(failUploader, j))
				}

				assert.True(st, attemptJob(failUploader, j))

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

						mockJQ.On("Complete", j).
							Return(1, nil)

						d := NewDrain(nilUploader)
						go d.Drain(&mockJQ)

						<-d.Schan
						val := <-d.Schan
						assert.Equal(st, jobqueue.Completed, val.State, "chunk state")

						return nil
					},
				},
				{
					description: "retries a job if it fails the first time",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()

						mockJQ.On("Next").
							Return(j, nil)

						mockJQ.On("Complete", j).
							Return(1, nil)

						once := true
						d := NewDrain(func(c *jobqueue.Chunk) error {
							if once {
								once = false
								return fmt.Errorf("error")
							}
							return nil
						})
						go d.Drain(&mockJQ)

						val := <-d.Schan
						assert.Equal(st, jobqueue.InProgress, val.State, "chunk state")
						val = <-d.Schan
						assert.Equal(st, jobqueue.InProgress, val.State, "chunk state")
						val = <-d.Schan
						assert.Equal(st, jobqueue.Completed, val.State, "chunk state")

						return nil
					},
				},
				{
					description: "retries a job if it fails fewer than the max times",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()

						mockJQ.On("Next").
							Return(j, nil).
							Times(jobqueue.MaxJobAttempts)

						mockJQ.On("Next").
							Return(nil, jobqueue.ErrNoActiveJobs)

						mockJQ.On("Complete", j).
							Return(1, nil)

						d := NewDrain(func(c *jobqueue.Chunk) error {
							if !j.AtMaxAttempts() {
								return fmt.Errorf("error")
							}
							return nil
						})
						go d.Drain(&mockJQ)

						<-d.Schan
						<-d.Schan
						<-d.Schan
						<-d.Schan
						val := <-d.Schan
						assert.Equal(st, jobqueue.Completed, val.State, "chunk state")

						return nil
					},
				},
				{
					description: "emits an error if something went wrong getting a job",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer

						err := fmt.Errorf("failed")

						mockJQ.On("Next").
							Return(nil, err)

						d := NewDrain(nilUploader)
						go d.Drain(&mockJQ)

						val := <-d.Echan
						assert.Equal(st, err, val, "error from next")

						return nil
					},
				},
				{
					description: "emits an error if something went wrong completing a job",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()

						err := fmt.Errorf("failed")

						mockJQ.On("Next").
							Return(j, nil)

						mockJQ.On("Complete", j).
							Return(0, err)

						d := NewDrain(nilUploader)
						go d.Drain(&mockJQ)

						<-d.Schan
						<-d.Schan
						val := <-d.Echan
						assert.Equal(st, err, val, "error from complete")

						return nil
					},
				},
				{
					description: "exits gracefully if a job is already in progress",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()
						j.Status.State = jobqueue.InProgress

						mockJQ.On("Next").
							Return(j, nil)

						d := NewDrain(nilUploader)
						d.Drain(&mockJQ)

						return nil
					},
				},
			},
		},
		{
			description: "when no jobs are active",
			testCases: []testCase{
				{
					description: "it returns without running",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer

						mockJQ.On("Next").
							Return(nil, jobqueue.ErrNoActiveJobs)

						d := NewDrain(nilUploader)
						d.Drain(&mockJQ)

						return nil
					},
				},
			},
		},
		{
			description: "when all active jobs are running",
			testCases: []testCase{
				{
					description: "it returns without running",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer

						mockJQ.On("Next").
							Return(nil, jobqueue.ErrAllActiveJobsInProgress)

						d := NewDrain(nilUploader)
						d.Drain(&mockJQ)

						return nil
					},
				},
			},
		},
		{
			description: "with many jobs and many goroutines",
			testCases: []testCase{
				{
					description: "nothing explodes with 1 job and 2 goroutines",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()

						mockJQ.On("Next").
							Return(j, nil).
							Once()

						mockJQ.On("Next").
							Return(nil, jobqueue.ErrNoActiveJobs)

						mockJQ.On("Complete", j).
							Return(1, nil)

						d := NewDrain(nilUploader)
						go d.Drain(&mockJQ)
						go d.Drain(&mockJQ)

						val := <-d.Schan
						assert.Equal(st, jobqueue.InProgress, val.State, "chunk state")
						val = <-d.Schan
						assert.Equal(st, jobqueue.Completed, val.State, "chunk state")

						return nil
					},
				},
				{
					description: "perpetually drains the queue when every job succeeds",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()

						times := 12

						mockJQ.On("Next").
							Return(j, nil).
							Times(times)

						mockJQ.On("Next").
							Return(nil, jobqueue.ErrNoActiveJobs)

						mockJQ.On("Complete", j).
							Return(1, nil).
							Times(times)

						d := NewDrain(nilUploader)
						go d.Drain(&mockJQ)

						for i := 0; i < times; i++ {
							val := <-d.Schan
							assert.Equal(st, jobqueue.InProgress, val.State, "chunk state")
							val = <-d.Schan
							assert.Equal(st, jobqueue.Completed, val.State, "chunk state")
						}

						return nil
					},
				},
				{
					description: "keeps draining even if 1 job fails fewer than the max times",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()
						j2 := randomJob()

						times := 12

						mockJQ.On("Next").
							Return(j2, nil).
							Times(jobqueue.MaxJobAttempts)

						mockJQ.On("Next").
							Return(j, nil).
							Times(times - 1)

						mockJQ.On("Next").
							Return(nil, jobqueue.ErrNoActiveJobs)

						mockJQ.On("Complete", j).
							Return(1, nil).
							Times(times - 1)

						mockJQ.On("Complete", j2).
							Return(1, nil).
							Once()

						d := NewDrain(func(c *jobqueue.Chunk) error {
							if c.ID == j2.Status.Chunk.ID {
								if !j2.AtMaxAttempts() {
									return fmt.Errorf("error")
								}
							}
							return nil
						})
						go d.Drain(&mockJQ)

						for i := 0; i < times*2+jobqueue.MaxJobAttempts-1; i++ {
							<-d.Schan
						}

						return nil
					},
				},
				{
					description: "leaves a job failed if it fails the max times",
					test: func(st *testing.T) error {
						var mockJQ jobqueuemocks.FIFOQueuer
						j := randomJob()
						failingJob := randomJob()

						times := 12

						mockJQ.On("Next").
							Return(j, nil).
							Once()

						mockJQ.On("Next").
							Return(failingJob, nil).
							Times(jobqueue.MaxJobAttempts)

						mockJQ.On("Next").
							Return(j, nil).
							Times(times - 2)

						mockJQ.On("Next").
							Return(nil, jobqueue.ErrNoActiveJobs)

						mockJQ.On("Complete", j).
							Return(1, nil).
							Times(times - 1)

						mockJQ.On("Complete", failingJob).
							Return(1, nil).
							Once()

						d := NewDrain(func(c *jobqueue.Chunk) error {
							if c.ID == failingJob.Status.Chunk.ID {
								return fmt.Errorf("error")
							}
							return nil
						})
						go d.Drain(&mockJQ)

						timesSeen := 0
						for i := 0; i < times*2+jobqueue.MaxJobAttempts-1; i++ {
							val := <-d.Schan
							if val.Chunk.ID == failingJob.Status.Chunk.ID {
								timesSeen++
								if timesSeen == jobqueue.MaxJobAttempts+1 {
									// maxjobattempts times it will be in progress, the last will be erred
									assert.Equal(st, jobqueue.Erred, val.State, "fail job failed")
								}
							}
						}

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
