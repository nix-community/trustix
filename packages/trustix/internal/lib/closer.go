// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

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
