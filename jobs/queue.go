package jobs

import (
	"fmt"
	"sync"
)

type Queue struct {
	jobs map[string]*Job
	mu   sync.RWMutex
}

var instance *Queue

func Initialize() {
	instance = &Queue{
		jobs: make(map[string]*Job),
	}
}

func GetQueue() *Queue {
	return instance
}

func (q *Queue) AddJob(job *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs[job.ID] = job
}

func (q *Queue) GetJob(id string) (*Job, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	job, exists := q.jobs[id]
	if !exists {
		return nil, fmt.Errorf("job not found")
	}
	return job, nil
}

func (q *Queue) UpdateJob(job *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs[job.ID] = job
}

func (q *Queue) ListJobs() []*Job {
	q.mu.RLock()
	defer q.mu.RUnlock()

	jobs := make([]*Job, 0, len(q.jobs))
	for _, job := range q.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}
