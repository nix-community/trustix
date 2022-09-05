// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package lib

import (
	"io"
	"sync"

	"go.uber.org/multierr"
)

type MultiCloser struct {
	mux     *sync.Mutex
	closers []io.Closer
}

func NewMultiCloser() *MultiCloser {
	return &MultiCloser{
		closers: []io.Closer{},
		mux:     &sync.Mutex{},
	}
}

func (m *MultiCloser) Add(closer io.Closer) {
	m.mux.Lock()
	m.closers = append(m.closers, closer)
	m.mux.Unlock()
}

func (m *MultiCloser) Close() error {
	m.mux.Lock()
	defer m.mux.Unlock()

	wg := new(sync.WaitGroup)

	errors := []error{}
	mux := &sync.Mutex{}

	for _, s := range m.closers {
		s := s

		wg.Add(1)
		go func() {
			err := s.Close()
			if err != nil {
				mux.Lock()
				errors = append(errors, err)
				mux.Unlock()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	if len(errors) == 0 {
		return nil
	}

	return multierr.Combine(errors...)
}
