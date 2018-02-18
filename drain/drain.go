package drain

import (
	"fmt"
	"time"

	"github.com/cameronwp/glacier/jobqueue"
)

// Actor is the action performed against a chunk.
type Actor func(*jobqueue.Chunk) error

// Drain is capable of emptying a job queue.
type Drain struct {
	action Actor
	Schan  chan jobqueue.Status
	Echan  chan error
}

// Drainer asynchronously performs actions against a job queue.
type Drainer interface {
	Drain(jobqueue.FIFOQueuer)
}

var _ Drainer = (*Drain)(nil)

// NewDrain creates a new drain with an Actor.
func NewDrain(a Actor) *Drain {
	return &Drain{
		action: a,
		Schan:  make(chan jobqueue.Status),
		Echan:  make(chan error),
	}
}

// Drain will run jobs asynchronously until there are no waiting jobs remaining.
func (d *Drain) Drain(q jobqueue.FIFOQueuer) {
	j, err := q.Next()
	if err != nil {
		switch err {
		case jobqueue.ErrNoActiveJobs:
		case jobqueue.ErrAllActiveJobsInProgress:
			break
		default:
			fmt.Println(err)
		}
		return
	}

	// the job is already running, try to get a different job
	if j.Status.State == jobqueue.InProgress {
		go d.Drain(q)
		return
	}

	j.Status.State = jobqueue.InProgress
	j.Status.StartedAt = time.Now()

	emitStatus(d.Schan, j)

	attemptJob(d.action, j)

	emitStatus(d.Schan, j)

	_, err = q.CompleteJob(j)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = q.ActivateOldestWaitingJob()
	if err != nil {
		// no problem if there are no active jobs
		switch err {
		case jobqueue.ErrNoWaitingJobs:
		case jobqueue.ErrMaxActiveJobs:
			break
		default:
			fmt.Println(err)
		}
		return
	}

	go d.Drain(q)
}

func attemptJob(a Actor, j *jobqueue.Job) {
	j.IncrAttempts()
	err := a(j.Status.Chunk)
	if err != nil {
		if !j.AtMaxAttempts() {
			// TODO: better logging
			// TODO: report restarting?
			fmt.Printf("%s failed, retrying | %s\n", j.Status.Chunk.ID, err)
			j.Status.State = jobqueue.Waiting
			return
		}

		j.Status.State = jobqueue.Erred
	} else {
		j.Status.State = jobqueue.Completed
	}

	j.Status.CompletedAt = time.Now()
}

func emitStatus(schan chan jobqueue.Status, j *jobqueue.Job) {
	schan <- j.Status
}
