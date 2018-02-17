package pool

import (
	"fmt"
	"time"
)

// Actor is the action performed against a chunk.
type Actor func(*Chunk) error

// Drain is capable of emptying a job queue.
type Drain struct {
	action Actor
	Schan  chan Status
	Echan  chan error
}

// Drainer asynchronously performs actions against a job queue.
type Drainer interface {
	Drain(*JobQueue)
}

var _ Drainer = (*Drain)(nil)

// NewDrain creates a new drain with an Actor.
func NewDrain(a Actor) *Drain {
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

	emitStatus(d.Schan, j)

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

func emitStatus(schan chan Status, j *Job) {
	schan <- j.Status
}
