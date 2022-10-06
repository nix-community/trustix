// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package future

import (
	"sync"
)

// future - holds the promise of a future computation
type Future[T any] struct {
	result T
	done   bool
	err    error
	cond   *sync.Cond
}

func NewFuture[T any]() *Future[T] {
	mux := sync.Mutex{}
	cond := sync.NewCond(&mux)

	return &Future[T]{
		cond: cond,
	}
}

func (f *Future[T]) Result() (T, error) {
	f.cond.L.Lock()
	for !f.done {
		f.cond.Wait()
	}
	f.cond.L.Unlock()

	return f.result, f.err
}

func (f *Future[T]) set(result T, err error) {
	f.result = result
	f.err = err
	f.done = true

	f.cond.L.Lock()
	f.cond.Broadcast()
	f.cond.L.Unlock()
}

// a future wrapper that only returns one future for a giving key
type KeyedFutures[T any] struct {
	mux     sync.Mutex
	futures map[string]*Future[T]
}

func NewKeyedFutures[T any]() *KeyedFutures[T] {
	return &KeyedFutures[T]{
		mux:     sync.Mutex{},
		futures: make(map[string]*Future[T]),
	}
}

func (kf *KeyedFutures[T]) Run(key string, fn func() (T, error)) *Future[T] {
	kf.mux.Lock()
	defer kf.mux.Unlock()

	fut, ok := kf.futures[key]
	if ok {
		return fut
	}

	fut = NewFuture[T]()
	kf.futures[key] = fut

	go func() {
		result, err := fn()

		kf.mux.Lock()
		delete(kf.futures, key)
		kf.mux.Unlock()

		fut.set(result, err)
	}()

	return fut
}
