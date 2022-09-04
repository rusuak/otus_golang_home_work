package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksQueue := make(chan Task)
	errorsBus := make(chan error, n)
	errLimitWasReached := make(chan bool)
	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go work(tasksQueue, errorsBus, wg)
	}
	go controlTasks(tasks, tasksQueue, errLimitWasReached)
	go controlErrors(errorsBus, m, errLimitWasReached)

	wg.Wait()
	close(errorsBus)

	select {
	case <-errLimitWasReached:
		return ErrErrorsLimitExceeded
	default:
		return nil
	}
}

func work(queue <-chan Task, errorsBus chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range queue {
		err := task()
		if err != nil {
			errorsBus <- err
		}
	}
}

func controlTasks(tasks []Task, queue chan<- Task, errLimitWasReached <-chan bool) {
	defer close(queue)

	for _, task := range tasks {
		select {
		case queue <- task:
		case <-errLimitWasReached:
			return
		}
	}
}

func controlErrors(errorsBus <-chan error, errorsLimit int, errLimitWasReached chan<- bool) {
	errorsCount := 0
	for {
		_, ok := <-errorsBus
		if !ok {
			return
		}
		errorsCount++
		if (errorsLimit > 0) && (errorsCount >= errorsLimit) {
			close(errLimitWasReached)
			return
		}
	}
}
