package pool

import (
	"fmt"
	"sync"
	"time"
)

// Chunk is a piece of a file to upload.
type Chunk struct {
	ID       string
	StartB   int64
	EndB     int64
	FilePath string
}

type actionState int

const (
	waiting actionState = iota
	inProgress
	completed
	erred
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
	waitingJobs    []*Job
	activeJobs     []*Job
	completedJobs  []*Job
	mux            sync.Mutex
	MaxConnections int
}

// FIFOQueuer is responsible for moving jobs from waiting -> active -> completed
// queues.
type FIFOQueuer interface {
	AddJob(Chunk) (int, error)
	ActivateOldestWaitingJob() (int, error)
	Next() (*Job, error)
	CompleteJob(*Job) (int, error)
}

var _ FIFOQueuer = (*JobQueue)(nil)

// AddJob creates a job from a chunk and adds a job to the waiting queue. It
// returns the number of waiting jobs and an error.
func (q *JobQueue) AddJob(c Chunk) (int, error) {
	if c.ID == "" {
		return len(q.waitingJobs), ErrInvalidChunk
	}

	j := Job{
		Status: Status{
			Chunk: &c,
			State: waiting,
		},
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	q.waitingJobs = append(q.waitingJobs, &j)
	return len(q.waitingJobs), nil
}

// ActivateOldestWaitingJob moves the oldest waiting job to the active queue. It
// returns the number of active jobs and an error.
func (q *JobQueue) ActivateOldestWaitingJob() (int, error) {
	q.mux.Lock()
	defer q.mux.Unlock()

	if len(q.activeJobs) >= q.MaxConnections {
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

	i := 0
	for i < len(q.activeJobs) {
		if q.activeJobs[i].Status.State == waiting {
			return q.activeJobs[i], nil
		}
	}

	return nil, ErrAllActiveJobsInProgress
}

// CompleteJob moves an active job to the completed queue. It returns the number
// of completed jobs and an error.
func (q *JobQueue) CompleteJob(j *Job) (int, error) {
	q.mux.Lock()
	defer q.mux.Unlock()

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

	q.completedJobs = append(q.completedJobs, q.activeJobs[i])
	q.activeJobs = append(q.activeJobs[:i], q.activeJobs[i:]...)

	return len(q.completedJobs), nil
}

type actor func(*Chunk) error

type Drain struct {
	action actor
	schan  chan Status
	echan  chan error
}

// Drainer asynchronously performs actions against a job queue.
type Drainer interface {
	Drain(*JobQueue)
}

var _ Drainer = (*Drain)(nil)

// New creates a new drain with an actor.
func New(a actor) *Drain {
	return &Drain{
		action: a,
	}
}

// Drain will run jobs asynchronously until there are no waiting jobs remaining.
func (d *Drain) Drain(q *JobQueue) {
	j, err := q.Next()
	if err != nil {
		switch err {
		case ErrNoActiveJobs:
		case ErrAllActiveJobsInProgress:
			break
		default:
			fmt.Println(err)
		}
		return
	}

	// the job is already running
	if j.Status.State == inProgress {
		go d.Drain(q)
		return
	}

	j.Status.State = inProgress
	j.numAttempts = j.numAttempts + 1
	j.Status.StartedAt = time.Now()

	emitStatus(d.schan, d.echan, j)

	err = d.action(j.Status.Chunk)
	// try 4 times to perform an action
	if err != nil && j.numAttempts < 4 {
		j.Status.State = waiting
		fmt.Printf("attempt #%d for %s failed, retrying | %s\n", j.numAttempts, j.Status.Chunk.ID, err)
		go d.Drain(q)
		return
	}

	j.Status.State = completed
	j.Status.CompletedAt = time.Now()

	emitStatus(d.schan, d.echan, j)

	_, err = q.CompleteJob(j)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = q.ActivateOldestWaitingJob()
	if err != nil {
		// no problem if there are no active jobs
		switch err {
		case ErrNoWaitingJobs:
		case ErrMaxActiveJobs:
			break
		default:
			fmt.Println(err)
		}
		return
	}

	go d.Drain(q)
}

func emitStatus(schan chan Status, echan chan error, j *Job) {
	//
}

// func drain(j *Job, a actor) {
// 	j.status.state = inProgress
// 	j.status.startedAt = time.Now()

// 	// TODO: clearly using channels incorrectly here

// 	// let listeners know this job is starting
// 	j.schan <- j.status

// 	j.numAttempts = j.numAttempts + 1
// 	err := u(j.chunk)

// 	// try 4 times to upload a chunk
// 	if err != nil && j.numAttempts < 4 {
// 		fmt.Printf("attempt #%d for %s failed, retrying | %s\n", j.numAttempts, j.chunk.ID, err)
// 		drain(j, u)
// 		return
// 	}

// 	// compile status
// 	j.status.completedAt = time.Now()

// 	if err != nil {
// 		j.status.state = erred
// 		j.echan <- err
// 	} else {
// 		j.status.state = completed
// 	}

// 	// let listeners know this job finished
// 	j.schan <- j.status
// }
