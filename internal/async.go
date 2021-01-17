package internal

import "sync"

type Job func()

type AsyncJobQueue struct {
	lock sync.Locker
	jobs []Job
}

func NewAsyncJobQueue() AsyncJobQueue {
	return AsyncJobQueue{lock: SpinLock()}
}

func (q *AsyncJobQueue) Push(job Job) (total int) {
	q.lock.Lock()
	q.jobs = append(q.jobs, job)
	total = len(q.jobs)
	q.lock.Unlock()
	return
}

func (q *AsyncJobQueue) Execute(group *sync.WaitGroup) {
	defer group.Done()
	q.lock.Lock()
	jobs := q.jobs
	q.jobs = nil
	q.lock.Unlock()
	for i := range jobs {
		jobs[i]()
	}
	return
}
