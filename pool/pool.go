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

type uploadState int

const (
	waiting uploadState = iota
	inProgress
	completed
	erred
)

// Status describes the completion status of an upload.
type Status struct {
	state       uploadState
	startedAt   time.Time
	completedAt time.Time
}

type job struct {
	chunk       Chunk
	status      Status
	numAttempts int
	schan       chan Status
	echan       chan error
}

type uploader func(Chunk) error

// Pool is a collection of upload jobs.
type Pool struct {
	waitingJobs    []*job
	activeJobs     []*job
	completedJobs  []*job
	mux            sync.Mutex
	uploader       uploader
	maxConnections int
	// for controlling executing during testing
	stopCycle bool
	stopDrain bool
}

// Pooler is capable of collecting file chunks and uploading them asynchronously.
type Pooler interface {
	Pool(Chunk) (*chan Status, *chan error)
	Cycle()
	Drain(*job)
}

var _ Pooler = (*Pool)(nil)

// New creates an empty Pool.
func New(u uploader, max int) Pool {
	return Pool{
		uploader:       u,
		maxConnections: max,
	}
}

// Pool creates a job. There's no guarantee about when it will run.
func (p *Pool) Pool(c Chunk) (*chan Status, *chan error) {
	j := job{
		chunk: c,
		status: Status{
			state: waiting,
		},
		// schan: make(chan Status),
		// echan: make(chan error),
	}

	p.mux.Lock()
	p.waitingJobs = append(p.waitingJobs, &j)
	p.mux.Unlock()

	p.Cycle()

	return &j.schan, &j.echan
}

// Cycle compares the number of jobs running to the max and activates waiting
// jobs if possible.
func (p *Pool) Cycle() {
	// testing only
	if p.stopCycle {
		return
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.activeJobs) >= p.maxConnections {
		// the max # of jobs are already running
		return
	}

	oldestWaitingJob := p.waitingJobs[0]

	// move the first waiting job to active jobs
	p.activeJobs = append(p.activeJobs, oldestWaitingJob)
	p.waitingJobs = append(p.waitingJobs[:0], p.waitingJobs[1:]...)

	go p.Drain(oldestWaitingJob)

	// aggressively look for more jobs
	go p.Cycle()
}

// Drain runs a job. It will try up to 4 times to complete the upload.
func (p *Pool) Drain(j *job) {
	// testing only
	if p.stopDrain {
		return
	}

	// run the job
	drain(j, p.uploader)

	// move the job from active to completed
	p.mux.Lock()
	defer p.mux.Unlock()

	i := 0
	for i < len(p.activeJobs) {
		if p.activeJobs[i].chunk.ID == j.chunk.ID {
			break
		}
		i++
	}

	p.completedJobs = append(p.completedJobs, p.activeJobs[i])
	p.activeJobs = append(p.activeJobs[:i], p.activeJobs[i:]...)

	// look for more waiting jobs
	go p.Cycle()
}

func drain(j *job, u uploader) {
	j.status.state = inProgress
	j.status.startedAt = time.Now()

	// let listeners know this job is starting
	j.schan <- j.status

	j.numAttempts = j.numAttempts + 1
	err := u(j.chunk)

	// try 4 times to upload a chunk
	if err != nil && j.numAttempts < 4 {
		fmt.Printf("attempt #%d for %s failed, retrying | %s\n", j.numAttempts, j.chunk.ID, err)
		drain(j, u)
		return
	}

	// compile status
	j.status.completedAt = time.Now()

	if err != nil {
		j.status.state = erred
		j.echan <- err
	} else {
		j.status.state = completed
	}

	// let listeners know this job finished
	j.schan <- j.status
}
