// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: MIT

package executor

import (
	"sync"
)

// LimitedParallellExecutor - Execute callback functions in parallell
type LimitedParallellExecutor struct {
	wg    *sync.WaitGroup
	mux   *sync.Mutex
	guard chan struct{}

	// Error returned by Wait(), cached for other Wait() & Add() invocations
	err  error
	done bool
}

func NewLimitedParallellExecutor(maxWorkers int) *LimitedParallellExecutor {
	e := &LimitedParallellExecutor{
		mux:   new(sync.Mutex),
		wg:    new(sync.WaitGroup),
		guard: make(chan struct{}, maxWorkers),

		err:  nil,
		done: false,
	}

	return e
}

func (e *LimitedParallellExecutor) Add(fn func() error) error {
	if e.err != nil {
		return e.err
	}

	e.wg.Add(1)

	e.guard <- struct{}{} // Block until a worker is available

	go func() {
		defer e.wg.Done()
		defer func() {
			<-e.guard
		}()

		err := fn()

		if err != nil && e.err == nil {
			e.mux.Lock()
			e.err = err
			e.mux.Unlock()
		}
	}()

	return nil
}

func (e *LimitedParallellExecutor) Wait() error {
	if e.done {
		return e.err
	}

	e.wg.Wait()

	e.done = true

	return e.err
}
