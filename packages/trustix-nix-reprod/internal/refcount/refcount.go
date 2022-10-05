package refcount

import (
	"sync/atomic"
)

// a reference counted value where cleanup action is run upon the number of references reaching zero
type RefCountedValue[T any] struct {
	Value T
	count int64
	done  bool
	fn    func() error
}

func NewRefCountedValue[T any](value T, fn func() error) *RefCountedValue[T] {
	return &RefCountedValue[T]{
		count: 0,
		done:  false,
		fn:    fn,
		Value: value,
	}
}

func (rv *RefCountedValue[T]) Incr() {
	if rv.count < 0 {
		panic("reference counter was negative")
	}

	if rv.done {
		panic("already done executing")
	}

	atomic.AddInt64(&rv.count, 1)
}

func (rv *RefCountedValue[T]) Decr() error {
	atomic.AddInt64(&rv.count, -1)

	if rv.count <= 0 && !rv.done {
		rv.done = true
		return rv.fn()
	}

	return nil
}
