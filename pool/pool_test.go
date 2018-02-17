package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func nilUploader(Chunk) error {
	return nil
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

				assert.Equal(st, 1, len(testJobQueue.waitingJobs), "wrong # waiting jobs")
				assert.Equal(st, c1.ID, testJobQueue.waitingJobs[0].Status.Chunk.ID, "wrong chunk ID in queue")

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

				assert.Equal(st, 2, numJobs, "wrong # waiting jobs")

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

				assert.Equal(st, 2, len(testJobQueue.waitingJobs), "wrong # waiting jobs")
				assert.Equal(st, c1.ID, testJobQueue.waitingJobs[0].Status.Chunk.ID, "wrong chunk ID in queue")
				assert.Equal(st, c2.ID, testJobQueue.waitingJobs[1].Status.Chunk.ID, "wrong chunk ID in queue")

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
			description: "",
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

// 	t.Parallel()
// 	for _, tc := range testCases {
// 		t.Run(tc.description, tc.test)
// 	}
// }

// func TestCycle(t *testing.T) {
// 	testCases := []struct {
// 		description string
// 		test        func(*testing.T)
// 	}{
// 		{
// 			description: "automatically cycles waiting to active if under the max",
// 			test: func(st *testing.T) {
// 				uploader := func(c Chunk) error {
// 					return nil
// 				}
// 				testPool := New(uploader, 1)
// 				testPool.stopDrain = true
// 				testPool.Pool(Chunk{})

// 				if len(testPool.waitingJobs) != 0 {
// 					st.Errorf("expected 0 waiting jobs, found %d", len(testPool.waitingJobs))
// 				}

// 				if len(testPool.activeJobs) != 1 {
// 					st.Errorf("expected 1 active jobs, found %d", len(testPool.activeJobs))
// 				}
// 			},
// 		},
// 		{
// 			description: "leaves jobs in waiting if at the max",
// 			test: func(st *testing.T) {
// 				uploader := func(c Chunk) error {
// 					return nil
// 				}
// 				testPool := New(uploader, 1)
// 				testPool.stopDrain = true
// 				testPool.Pool(Chunk{})
// 				testPool.Pool(Chunk{})

// 				if len(testPool.waitingJobs) != 1 {
// 					st.Errorf("expected 1 waiting jobs, found %d", len(testPool.waitingJobs))
// 				}

// 				if len(testPool.activeJobs) != 1 {
// 					st.Errorf("expected 1 active jobs, found %d", len(testPool.activeJobs))
// 				}
// 			},
// 		},
// 		{
// 			description: "definitely doesn't add jobs if over the max",
// 			test: func(st *testing.T) {
// 				uploader := func(c Chunk) error {
// 					return nil
// 				}
// 				testPool := New(uploader, 1)
// 				testPool.stopDrain = true
// 				testPool.activeJobs = []*job{&job{}, &job{}}
// 				testPool.Pool(Chunk{})
// 				testPool.Pool(Chunk{})

// 				if len(testPool.waitingJobs) != 2 {
// 					st.Errorf("expected 2 waiting jobs, found %d", len(testPool.waitingJobs))
// 				}
// 			},
// 		},
// 	}

// 	t.Parallel()
// 	for _, tc := range testCases {
// 		t.Run(tc.description, tc.test)
// 	}
// }

// func TestDrain(t *testing.T) {
// 	testCases := []struct {
// 		description string
// 		test        func(*testing.T)
// 	}{
// 		{
// 			description: "updates status to inProgress",
// 			test: func(st *testing.T) {
// 				u := func(c Chunk) error {
// 					return nil
// 				}
// 				j := job{}
// 				drain(&j, u)
// 			},
// 		},
// 	}

// 	t.Parallel()
// 	for _, tc := range testCases {
// 		t.Run(tc.description, tc.test)
// 	}
// }

// // res := <-*schan
// // if res.state != completed {
// // 	t.Errorf("expect test to succeed, found %d", res.state)
// // }
