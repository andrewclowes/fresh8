package common

import (
	"sync/atomic"
)

const defaultLimit = 100

// Task must be implemented by types that want to use the limiter
type Task func()

// Limiter limits the amount of concurrency for concurrent tasks
type Limiter interface {
	Execute(task Task)
	GetNumInProgress() int32
	Wait()
}

// SemaphoreLimiter uses a semaphoreto limit concurrency
type SemaphoreLimiter struct {
	sem chan bool
	num int32
}

// NewLimiter creates a new limiter
func NewLimiter(limit int) Limiter {
	if limit <= 0 {
		limit = defaultLimit
	}
	l := &SemaphoreLimiter{
		sem: make(chan bool, limit),
	}
	return l
}

// Execute runs a given task or will wait until a go routine is available
func (l *SemaphoreLimiter) Execute(task Task) {
	l.sem <- true
	atomic.AddInt32(&l.num, 1)
	go func() {
		defer func() {
			<-l.sem
			atomic.AddInt32(&l.num, -1)
		}()
		task()
	}()
}

// GetNumInProgress shows how many go routines are active in progress
func (l *SemaphoreLimiter) GetNumInProgress() int32 {
	return l.num
}

// Wait will wait until all previously executed jobs have finished
func (l *SemaphoreLimiter) Wait() {
	for i := 0; i < cap(l.sem); i++ {
		l.sem <- true
	}
}
