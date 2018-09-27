package cron

import (
	"sync"
	"time"

	"github.com/Hatch1fy/errors"
	"github.com/Path94/atoms"
)

// newJob is the intended and safe way to create a job
func newJob(query string, fn func()) (jp *Job, err error) {
	var j Job
	// Parse the inbound query and set the descriptor, reference, value, and locations
	if j.descriptor, j.reference, j.value, j.loc, err = parseQuery(query); err != nil {
		// Error parsing query, return
		return
	}

	// Set the function
	j.fn = fn
	// Set jp as a pointer to our job
	jp = &j
	return
}

// Job represents a job entry
type Job struct {
	mu sync.Mutex

	// Descriptor, currently can be "@" or "every"
	descriptor string
	// Reference, currently can be "time" or a unit of time shorthand (e.g. ms, s, m, h)
	reference string
	// Time location (time zone)
	loc *time.Location
	// Value of the job, the meaning of the value varies between time and non-time references
	value int64

	// Function to be called when job is ran
	fn func()

	// Closed state
	closed atoms.Bool
}

func (j *Job) getSleepDuration() (dur time.Duration) {
	j.mu.Lock()
	defer j.mu.Unlock()

	// Set duration to the job value
	dur = time.Duration(j.value)

	switch j.reference {
	case "ms":
		// Convert the duration to milliseconds
		dur *= time.Millisecond
	case "s":
		// Convert the duration to seconds
		dur *= time.Second
	case "m":
		// Convert the duration to minute
		dur *= time.Minute
	case "h":
		// Convert the duration to hour
		dur *= time.Hour
	case "time":
		now := time.Now()
		start := GetStartOfDay(now)
		minutes := time.Minute * time.Duration(j.value)
		target := start.Add(minutes)

		if target.Before(now) {
			// Target already occurred today, set target for tomorrow
			target = target.AddDate(0, 0, 1)
		}

		dur = time.Second * time.Duration(target.Unix()-now.Unix())

	case "":
		return 0

	default:
		// TODO: Replace this with proper error handling, though this should NEVER happen.
		panic("invalid reference, how was this called? " + j.reference)
	}

	return
}

// sleep will sleep for the Job's intended sleep amount
func (j *Job) sleep() {
	dur := j.getSleepDuration()
	// Sleep for duration
	time.Sleep(dur)
}

func (j *Job) action() {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.closed.Get() {
		// Job has closed, return
		return
	}

	// Run job function
	j.fn()
}

func (j *Job) isNotRepeating() (ok bool) {
	// If job has an "@" descriptor rather than an "every" descriptor, it is not repeating
	return j.descriptor == "@"
}

// Run will run a job
func (j *Job) Run() {
	for {
		// Sleep for job delay time
		j.sleep()
		// Perform job action
		j.action()
		// Check to see if job is repeating
		if j.isNotRepeating() {
			return
		}
	}
}

// Close will close a job
func (j *Job) Close() (err error) {
	// Attempt to close job atomically
	if !j.closed.Set(true) {
		return errors.ErrIsClosed
	}

	// Acquire lock
	j.mu.Lock()
	// Defer the release of lock
	defer j.mu.Unlock()

	// Set reference to empty
	j.reference = ""
	// Set func to nil
	j.fn = nil
	return
}
