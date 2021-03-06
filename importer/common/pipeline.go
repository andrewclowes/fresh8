package common

import (
	"sync"
)

// StepRunner exposes functions that take in a chan interface{} and outputs
// to a chan interface{}
type StepRunner interface {
	Run(in <-chan interface{}, errc chan<- error) chan interface{}
}

// Step takes one input channel and one output channel
type Step func(in <-chan interface{}, out chan interface{}, errc chan<- error)

// Run takes an input channel, and a series of operators, and uses the output
// of each successive operator as the input for the next
func (o Step) Run(in <-chan interface{}, errc chan<- error) chan interface{} {
	out := make(chan interface{})
	go func() {
		o(in, out, errc)
		close(out)
	}()
	return out
}

// StepHandler is a function that handles the step logic
type StepHandler func(in interface{}, errc chan<- error) interface{}

// ConcurrentStep is a step that processes each item concurrently
type ConcurrentStep struct {
	handle StepHandler
}

// NewConcurrentStep creates a pipeline step in which items
// are processed concurrently
func NewConcurrentStep(handler StepHandler) StepRunner {
	step := &ConcurrentStep{
		handle: handler,
	}
	return step
}

// Run takes an input channel, and a series of operators, and uses the output
// of each successive operator as the input for the next
func (s *ConcurrentStep) Run(in <-chan interface{}, errc chan<- error) chan interface{} {
	out := make(chan interface{})
	go func() {
		var wg sync.WaitGroup
		for m := range in {
			wg.Add(1)
			go func(n interface{}) {
				defer wg.Done()
				o := s.handle(n, errc)
				if o != nil {
					out <- o
				}
			}(m)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

// StepLimiterFactory creates new instances of a limiter
type StepLimiterFactory func() Limiter

// RateLimitedStep is a step that processes each item concurrently but is
// rate limited
type RateLimitedStep struct {
	handle        StepHandler
	createLimiter StepLimiterFactory
}

// NewRateLimitedStep creates a new rate limited pipeline step
func NewRateLimitedStep(handler StepHandler, limit int) StepRunner {
	limiterFactory := func() Limiter {
		limiter := NewLimiter(limit)
		return limiter
	}
	step := &RateLimitedStep{
		handle:        handler,
		createLimiter: limiterFactory,
	}
	return step
}

// Run takes an input channel, and a series of operators, and uses the output
// of each successive operator as the input for the next
func (s *RateLimitedStep) Run(in <-chan interface{}, errc chan<- error) chan interface{} {
	limiter := s.createLimiter()
	out := make(chan interface{})
	go func() {
		for m := range in {
			n := m
			limiter.Execute(func() {
				o := s.handle(n, errc)
				if o != nil {
					out <- o
				}
			})
		}
		limiter.Wait()
		close(out)
	}()
	return out
}

// Steps is a slice of Steps that can be applied in sequence
type Steps []StepRunner

// NewSteps create a new slice of steps
func NewSteps(steps ...StepRunner) Steps {
	return Steps(steps)
}

// Run takes an input channel and runs the operators in the slice in order
func (s Steps) Run(in chan interface{}) (chan interface{}, chan error) {
	errc := make(chan error)
	for _, m := range s {
		in = m.Run(in, errc)
	}
	return in, errc
}

// PipelineJob is a job that runs a chainable series of steps
type PipelineJob struct {
	steps  *Steps
	logger Logger
}

// NewPipelineJob creates a new pipeline job for a series of steps
func NewPipelineJob(steps *Steps, logger Logger) Job {
	job := &PipelineJob{
		steps:  steps,
		logger: logger,
	}
	return job
}

// Run starts the job and will wait until all steps to finish
// before returning
func (j *PipelineJob) Run() {
	out, errc := j.steps.Run(nil)
loop:
	for {
		select {
		case err := <-errc:
			j.logger.Error(err.Error())
		case _, ok := <-out:
			if !ok {
				break loop
			}
		}
	}
}
