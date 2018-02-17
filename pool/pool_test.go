package pool

import (
	"testing"
)

func TestPool(t *testing.T) {
	testCases := []struct {
		description string
		test        func(*testing.T)
	}{
		{
			description: "creates a waiting job from a chunk",
			test: func(st *testing.T) {
				uploader := func(c Chunk) error {
					return nil
				}
				testPool := New(uploader, 1)
				testPool.stopCycle = true
				id := "asdf1234"
				testPool.Pool(Chunk{
					ID: id,
				})

				if len(testPool.waitingJobs) != 1 {
					st.Errorf("expected 1 waiting job, found %d", len(testPool.waitingJobs))
				}

				if testPool.waitingJobs[0].chunk.ID != id {
					st.Errorf("expected the job's chunk ID to match %s, found %s", testPool.waitingJobs[0].chunk.ID, id)
				}
			},
		},
		{
			description: "keeps adding new waiting jobs",
			test: func(st *testing.T) {
				uploader := func(c Chunk) error {
					return nil
				}
				testPool := New(uploader, 1)
				testPool.stopCycle = true
				testPool.Pool(Chunk{})
				testPool.Pool(Chunk{})

				if len(testPool.waitingJobs) != 2 {
					st.Errorf("expected 2 waiting jobs, found %d", len(testPool.waitingJobs))
				}
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.description, tc.test)
	}
}

func TestCycle(t *testing.T) {
	testCases := []struct {
		description string
		test        func(*testing.T)
	}{
		{
			description: "automatically cycles waiting to active if under the max",
			test: func(st *testing.T) {
				uploader := func(c Chunk) error {
					return nil
				}
				testPool := New(uploader, 1)
				testPool.stopDrain = true
				testPool.Pool(Chunk{})

				if len(testPool.waitingJobs) != 0 {
					st.Errorf("expected 0 waiting jobs, found %d", len(testPool.waitingJobs))
				}

				if len(testPool.activeJobs) != 1 {
					st.Errorf("expected 1 active jobs, found %d", len(testPool.activeJobs))
				}
			},
		},
		{
			description: "leaves jobs in waiting if at the max",
			test: func(st *testing.T) {
				uploader := func(c Chunk) error {
					return nil
				}
				testPool := New(uploader, 1)
				testPool.stopDrain = true
				testPool.Pool(Chunk{})
				testPool.Pool(Chunk{})

				if len(testPool.waitingJobs) != 1 {
					st.Errorf("expected 1 waiting jobs, found %d", len(testPool.waitingJobs))
				}

				if len(testPool.activeJobs) != 1 {
					st.Errorf("expected 1 active jobs, found %d", len(testPool.activeJobs))
				}
			},
		},
		{
			description: "definitely doesn't add jobs if over the max",
			test: func(st *testing.T) {
				uploader := func(c Chunk) error {
					return nil
				}
				testPool := New(uploader, 1)
				testPool.stopDrain = true
				testPool.activeJobs = []*job{&job{}, &job{}}
				testPool.Pool(Chunk{})
				testPool.Pool(Chunk{})

				if len(testPool.waitingJobs) != 2 {
					st.Errorf("expected 2 waiting jobs, found %d", len(testPool.waitingJobs))
				}
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.description, tc.test)
	}
}

func TestDrain(t *testing.T) {
	testCases := []struct {
		description string
		test        func(*testing.T)
	}{
		{
			description: "updates status to inProgress",
			test: func(st *testing.T) {
				u := func(c Chunk) error {
					return nil
				}
				j := job{}
				drain(&j, u)
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.description, tc.test)
	}
}

// res := <-*schan
// if res.state != completed {
// 	t.Errorf("expect test to succeed, found %d", res.state)
// }
