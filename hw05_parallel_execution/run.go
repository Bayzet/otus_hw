package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrErrorsNoGoroutines  = errors.New("not goroutines for work")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if n <= 0 {
		return ErrErrorsNoGoroutines
	}

	tasksChan := make(chan Task)
	var errCount int32
	wg := sync.WaitGroup{}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for task := range tasksChan {
				err := task()
				if err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errCount) >= int32(m) {
			break
		}
		tasksChan <- task
	}
	close(tasksChan)
	wg.Wait()

	if errCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
