package jobqueue

import (
	"fmt"
	"sync"
	"time"
)

// MaxJobAttempts is the max number of times to attempt a job before considering
// it a failure.
const MaxJobAttempts = 4

// Chunk is a piece of a file to upload.
type Chunk struct {
	ID       string
	StartB   int64
	EndB     int64
	FilePath string
}

type actionState int

const (
	// Waiting describes a chunk that has yet to start.
	Waiting actionState = iota
	// InProgress describes a chunk that is running.
	InProgress
	// Completed describes a chunk that successfully ran.
	Completed
	// Erred describes a chunk that did not successfully run.
	Erred
)

// Status describes the completion status of an upload.
type Status struct {
	State       actionState
	Chunk       *Chunk
	StartedAt   time.Time
	CompletedAt time.Time
}

// Job is a piece of work.
type Job struct {
	Status      Status
	numAttempts int
}

// IncrAttempts bumps the number of attempts by 1.
func (j *Job) IncrAttempts() {
	j.numAttempts++
}

// AtMaxAttempts returns whether or not the Job has been tried the max number of
// times.
func (j *Job) AtMaxAttempts() bool {
	return j.numAttempts >= MaxJobAttempts
}

var (
	// ErrInvalidChunk occurs when attempting to create a job with a chunk without
	// an ID.
	ErrInvalidChunk = fmt.Errorf("missing chunk ID")

	// ErrMaxActiveJobs occurs when attempting to activate a job and connections
	// are already maxed out.
	ErrMaxActiveJobs = fmt.Errorf("active jobs are already maxed")

	// ErrNoWaitingJobs occurs when trying to activate a job and no jobs are
	// waiting to be activated.
	ErrNoWaitingJobs = fmt.Errorf("no jobs are waiting")

	// ErrNoActiveJobs occurs when trying to grab the next active job but none
	// exist.
	ErrNoActiveJobs = fmt.Errorf("no jobs are active")

	// ErrAllActiveJobsInProgress means all active jobs are running.
	ErrAllActiveJobsInProgress = fmt.Errorf("all active jobs are in progress")
)

// JobQueue is a collection of upload jobs.
type JobQueue struct {
	// Open determines whether completing or adding jobs should automatically
	// activate new ones. If not, you will need to manually call
	// ActivateOldestWaitingJob().
	Open bool
	// MaxJobs is the max # of jobs that can be active simultaneously.
	MaxJobs       int
	waitingJobs   []*Job
	activeJobs    []*Job
	completedJobs []*Job
	mux           sync.Mutex
}

// NewJobQueue returns a new job queue. Max represents the max # active jobs.
// Open represents whether or not the queue should automatically try to keep the
// active jobs filled (imagine that 'open' refers to a valve between a pipe of
// waiting jobs into active jobs).
func NewJobQueue(max int, open bool) *JobQueue {
	return &JobQueue{
		Open:    open,
		MaxJobs: max,
	}
}

// FIFOQueuer is responsible for tracking queued jobs and providing jobs that
// are ready to run.
type FIFOQueuer interface {
	Add(Chunk) (int, error)
	Next() (*Job, error)
	Complete(*Job) (int, error)
}

var _ FIFOQueuer = (*JobQueue)(nil)

// Add creates a job from a chunk and adds a job to the waiting queue. It
// returns the number of waiting jobs and an error.
func (q *JobQueue) Add(c Chunk) (int, error) {
	if c.ID == "" {
		return len(q.waitingJobs), ErrInvalidChunk
	}

	j := Job{
		Status: Status{
			Chunk: &c,
			State: Waiting,
		},
	}

	q.mux.Lock()
	q.waitingJobs = append(q.waitingJobs, &j)
	q.mux.Unlock()

	if q.Open {
		_, err := q.ActivateOldestWaitingJob()
		if err != nil {
			// no problem if there are no active jobs
			switch err {
			case ErrNoWaitingJobs:
			case ErrMaxActiveJobs:
				break
			default:
				return len(q.waitingJobs), err
			}
		}
	}

	return len(q.waitingJobs), nil
}

// ActivateOldestWaitingJob moves the oldest waiting job to the active queue. It
// returns the number of active jobs and an error.
func (q *JobQueue) ActivateOldestWaitingJob() (int, error) {
	q.mux.Lock()
	defer q.mux.Unlock()

	if len(q.activeJobs) >= q.MaxJobs {
		return len(q.waitingJobs), ErrMaxActiveJobs
	}

	if len(q.waitingJobs) == 0 {
		return 0, ErrNoWaitingJobs
	}

	oldestWaitingJob := q.waitingJobs[0]

	// move the first waiting job to active jobs
	q.activeJobs = append(q.activeJobs, oldestWaitingJob)
	q.waitingJobs = append(q.waitingJobs[:0], q.waitingJobs[1:]...)

	return len(q.waitingJobs), nil
}

// Next returns a job that is active but is not in progress.
func (q *JobQueue) Next() (*Job, error) {
	if len(q.activeJobs) == 0 {
		return nil, ErrNoActiveJobs
	}

	for i := range q.activeJobs {
		if q.activeJobs[i].Status.State == Waiting {
			return q.activeJobs[i], nil
		}
	}

	return nil, ErrAllActiveJobsInProgress
}

// Complete moves an active job to the completed queue. It returns the number of
// completed jobs and an error.
func (q *JobQueue) Complete(j *Job) (int, error) {
	q.mux.Lock()

	i := 0
	found := false
	for i < len(q.activeJobs) {
		if q.activeJobs[i].Status.Chunk.ID == j.Status.Chunk.ID {
			found = true
			break
		}
		i++
	}

	if !found {
		return len(q.completedJobs), fmt.Errorf("job with chunk ID '%s' is not active", j.Status.Chunk.ID)
	}

	q.completedJobs = append(q.completedJobs, j)
	q.activeJobs = append(q.activeJobs[:i], q.activeJobs[i+1:]...)

	q.mux.Unlock()

	if q.Open {
		_, err := q.ActivateOldestWaitingJob()
		if err != nil {
			// no problem if there are no active jobs
			switch err {
			case ErrNoWaitingJobs:
			case ErrMaxActiveJobs:
				break
			default:
				return len(q.waitingJobs), err
			}
		}
	}

	return len(q.completedJobs), nil
}
