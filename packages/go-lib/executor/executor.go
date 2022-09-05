// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package executor

import (
	"sync"
)

// ParallellExecutor - Execute callback functions in parallell
type ParallellExecutor struct {
	errChan chan error
	wg      *sync.WaitGroup
	mux     *sync.Mutex

	// Error returned by Wait(), cached for other Wait() invocations
	err  error
	done bool
}

func NewParallellExecutor() *ParallellExecutor {
	return &ParallellExecutor{
		errChan: make(chan error),
		mux:     new(sync.Mutex),
		wg:      new(sync.WaitGroup),

		err:  nil,
		done: false,
	}
}

func (e *ParallellExecutor) Add(fn func() error) {
	e.wg.Add(1)

	go func() {
		defer e.wg.Done()

		err := fn()
		if err != nil {
			e.errChan <- err
		}
	}()
}

func (e *ParallellExecutor) Wait() error {
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
