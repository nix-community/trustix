package executor

import (
	"sync"
)

// LimitedParallellExecutor - Execute callback functions in parallell
type LimitedParallellExecutor struct {
	errChan chan error
	wg      *sync.WaitGroup
	mux     *sync.Mutex
	guard   chan struct{}

	// Error returned by Wait(), cached for other Wait() invocations
	err  error
	done bool
}

func NewLimitedParallellExecutor(maxWorkers int) *LimitedParallellExecutor {
	return &LimitedParallellExecutor{
		errChan: make(chan error),
		mux:     new(sync.Mutex),
		wg:      new(sync.WaitGroup),
		guard:   make(chan struct{}, maxWorkers),

		err:  nil,
		done: false,
	}
}

func (e *LimitedParallellExecutor) Add(fn func() error) {
	e.wg.Add(1)

	// TODO: Return error if a previously executed function has errored out

	e.guard <- struct{}{} // Block until a worker is available

	go func() {
		defer e.wg.Done()
		defer func() {
			<-e.guard
		}()

		err := fn()
		if err != nil {
			e.errChan <- err
		}
	}()
}

func (e *LimitedParallellExecutor) Wait() error {
	e.mux.Lock()
	defer e.mux.Unlock()

	if e.done {
		return e.err
	}

	var err error

	// Ensure channel is closed
	go func() {
		e.wg.Wait()
		close(e.errChan)
	}()

	for err = range e.errChan {
		if err != nil {
			break
		}
	}

	e.done = true
	e.err = err

	return err
}
