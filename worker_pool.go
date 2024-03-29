package workerpool

import (
	"errors"
	"log"
	"sync"
)

type Task func()

type IWorkerPool interface {
	Start()
	Submit(Task) error
	Shutdown()
}

type workerPool struct {
	name       string
	maxWorkers int

	workers   []*worker
	taskQueue chan Task
	workerWg  *sync.WaitGroup

	quit bool
}

func New(name string, maxWorkers int, queueSize int) IWorkerPool {
	taskQueue := make(chan Task, queueSize)
	var wg sync.WaitGroup

	var workers []*worker
	for i := 0; i < maxWorkers; i++ {
		w := newWorker(i, taskQueue)
		workers = append(workers, w)
	}

	return &workerPool{
		name:       name,
		maxWorkers: maxWorkers,
		workers:    workers,
		taskQueue:  taskQueue,
		workerWg:   &wg,
	}
}

func (wp *workerPool) Start() {
	log.Println("worker pool execution started...")

	for _, w := range wp.workers {
		wp.workerWg.Add(1)
		go w.run(wp.workerWg)
	}
}

func (wp *workerPool) Submit(job Task) error {
	if !wp.quit {
		wp.taskQueue <- job
		return nil
	}
	return errors.New("shutting down, cannot submit the job")
}

func (wp *workerPool) Shutdown() {
	wp.quit = true
	close(wp.taskQueue)
	wp.workerWg.Wait()
	log.Println("worker pool shutdown ...")
}
