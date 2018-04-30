package common

// StepRunner exposes functions that take in a chan interface{} and outputs
// to a chan interface{}
type StepRunner interface {
	Run(in <-chan interface{}) chan interface{}
}

// Step takes one input channel and one output channel
type Step func(in <-chan interface{}, out chan interface{})

// Run takes an input channel, and a series of operators, and uses the output
// of each successive operator as the input for the next
func (o Step) Run(in <-chan interface{}) chan interface{} {
	out := make(chan interface{})
	go func() {
		o(in, out)
		close(out)
	}()
	return out
}

// StepHandler is a function that handles the step logic
type StepHandler func(in interface{}) (interface{}, error)

// ConcurrentStep is a step that processes each item concurrently
type ConcurrentStep struct {
	handle StepHandler
}

// Run takes an input channel, and a series of operators, and uses the output
// of each successive operator as the input for the next
func (s ConcurrentStep) Run(in <-chan interface{}) chan interface{} {
	out := make(chan interface{})
	go func() {
		for m := range in {
			go func(n interface{}) {
				o, err := s.handle(n)
				if err == nil {
					out <- o
				}
			}(m)
		}
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
func (s *RateLimitedStep) Run(in <-chan interface{}) chan interface{} {
	limiter := s.createLimiter()
	out := make(chan interface{})
	go func() {
		for m := range in {
			n := m
			limiter.Execute(func() {
				o, err := s.handle(n)
				if err == nil {
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
func (s Steps) Run(in chan interface{}) chan interface{} {
	for _, m := range s {
		in = m.Run(in)
	}
	return in
}
