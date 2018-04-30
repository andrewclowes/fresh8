package common

// Job is the basic unit of work
type Job interface {
	Run()
}

// JobRegistry enables new jobs to be registered
type JobRegistry interface {
	Register(job Job)
}

// JobRunner allows registered
type JobRunner interface {
	RunJobs() error
}
