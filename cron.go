package cron

import (
	"regexp"
	"sync"

	"github.com/Hatch1fy/errors"
	"github.com/PathDNA/atoms"
)

var (
	millisecondRegEx = regexp.MustCompile(`^[\d]+[ms]{2}$`)
	secondRegEx      = regexp.MustCompile(`^[\d]+[s]{1}$`)
	minuteRegEx      = regexp.MustCompile(`^[\d]+[m]{1}$`)
	hourRegEx        = regexp.MustCompile(`^[\d]+[h]{1}$`)
)

// New will return a new instance of Cron
func New() *Cron {
	var c Cron
	return &c
}

// Cron manages cron-jobs
type Cron struct {
	mu sync.RWMutex
	// Slice of associated jobs
	jobs []*Job
	// Closed state of Cron
	closed atoms.Bool
}

// Set will set a cron-job
func (c *Cron) Set(query string, fn func()) (err error) {
	// Acquire lock
	c.mu.Lock()
	// Defer the release of the lock
	defer c.mu.Unlock()

	// Ensure cron is not closed
	if c.closed.Get() {
		// Cron is closed, return ErrIsClosed
		return errors.ErrIsClosed
	}

	var j *Job
	// Create a new job with the given name, query, and function
	if j, err = newJob(query, fn); err != nil {
		return
	}

	// Append the job to our Cron's job list
	c.jobs = append(c.jobs, j)
	// Run the job in a goroutine
	go j.Run()
	return
}

// Close will close a Cron and it's child jobs
func (c *Cron) Close() (err error) {
	if !c.closed.Set(true) {
		return errors.ErrIsClosed
	}

	// Acquire lock so any Set can finish
	c.mu.Lock()
	// Defer the release of the lock
	defer c.mu.Unlock()

	// Iterate through all jobs
	for _, j := range c.jobs {
		// Close job
		j.Close()
	}

	// Set job slice to nil
	c.jobs = nil
	return
}
