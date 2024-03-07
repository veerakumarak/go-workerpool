package workerpool

import "sync"

type worker struct {
	id        int
	taskQueue chan Task
}

func newWorker(id int, taskQueue chan Task) *worker {
	return &worker{
		id:        id,
		taskQueue: taskQueue,
	}
}

func (w worker) run(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-w.taskQueue:
			if !ok {
				return
			}
			task()
		}
	}
}
