package cron

import (
	"time"
)

// New creates a new job
func New(fn func()) (jp *Job) {
	var j Job
	// Set the function
	j.fn = fn
	return &j
}

// Job represents a job entry
type Job struct {
	// Function to be called when job is ran
	fn func()
}

// runAfter will run a function after waiting for a given duration
func (j *Job) runAfter(duration time.Duration) {
	time.Sleep(duration)
	j.fn()
}

// runEvery will run a function continuously with the given duration as a delay
func (j *Job) runEvery(duration time.Duration) {
	for {
		j.runAfter(duration)
	}
}

// runAt will run a function after waiting until a target time
func (j *Job) runAt(target time.Time) {
	j.runAfter(getDelay(target))
}

// runEveryAt will run a function continuously with the target time every day
func (j *Job) runEveryAt(target time.Time) {
	for {
		j.runAt(target)
	}
}

// After will run a function after waiting for a given duration
func (j *Job) After(duration time.Duration) {
	go j.runAfter(duration)
}

// Every will run a function continuously with the given duration as a delay
func (j *Job) Every(duration time.Duration) {
	go j.runEvery(duration)
}

// At will run a function after waiting until a target time
func (j *Job) At(target time.Time) {
	go j.runAt(target)
}

// EveryAt will run a function continuously with the target time every day
func (j *Job) EveryAt(duration time.Duration) {
	go j.runEvery(duration)
}
