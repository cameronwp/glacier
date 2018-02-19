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
			d.Echan <- err
		}
		return
	}

	// Next doesn't perform a lock, so it's possible another Drain goroutine
	// picked up the same job. Double check that it isn't running.
	if j.Status.State == jobqueue.InProgress {
		go d.Drain(q)
		return
	}

	j.Lock()
	defer j.Unlock()

	j.Status.State = jobqueue.InProgress
	j.Status.StartedAt = time.Now()

	emitStatus(d.Schan, j)

	complete := attemptJob(d.action, j)

	if !complete {
		go d.Drain(q)
		return
	}

	emitStatus(d.Schan, j)

	_, err = q.Complete(j)
	if err != nil {
		d.Echan <- err
		return
	}

	go d.Drain(q)
}

func attemptJob(a Actor, j *jobqueue.Job) bool {
	j.IncrAttempts()
	err := a(j.Status.Chunk)
	if err != nil {
		if !j.AtMaxAttempts() {
			// TODO: better logging
			// TODO: report restarting?
			fmt.Printf("%s failed, retrying | %s\n", j.Status.Chunk.Path, err)
			j.Status.State = jobqueue.Waiting
			return false
		}

		j.Status.State = jobqueue.Erred
	} else {
		j.Status.State = jobqueue.Completed
	}

	j.Status.CompletedAt = time.Now()
	return true
}

func emitStatus(schan chan jobqueue.Status, j *jobqueue.Job) {
	schan <- j.Status
}
