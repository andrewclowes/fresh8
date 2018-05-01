package common

// Job is the basic unit of work
type Job interface {
	Run()
}

// JobRunner runs a slice of jobs
type JobRunner interface {
	Run(jobs ...Job)
}

// SequentialJobRunner runs a slice of jobs in sequence
type SequentialJobRunner struct{}

// Run will execute the jobs sequentially
func (r *SequentialJobRunner) Run(jobs ...Job) {
	for _, j := range jobs {
		j.Run()
	}
}

// NewJobRunner creates a new job runner
func NewJobRunner() JobRunner {
	return &SequentialJobRunner{}
}
