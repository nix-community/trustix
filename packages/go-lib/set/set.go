package set

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"sync"
)

type Set[T constraints.Ordered] struct {
	values map[T]struct{}
	mux    *sync.RWMutex
}

// NewSet returns a new set (thread unsafe)
func NewSet[T constraints.Ordered]() *Set[T] {
	return &Set[T]{
		values: make(map[T]struct{}),
		mux:    nil,
	}
}

// NewSafeSet returns a new set (thread safe)
func NewSafeSet[T constraints.Ordered]() *Set[T] {
	return &Set[T]{
		values: make(map[T]struct{}),
		mux:    &sync.RWMutex{},
	}
}

// Return all values of the set (in sorted order)
func (s *Set[T]) Values() []T {
	if s.mux != nil {
		s.mux.RLock()
		defer s.mux.RUnlock()
	}

	values := make([]T, len(s.values))

	i := 0
	for v, _ := range s.values {
		values[i] = v
		i++
	}

	slices.Sort(values)

	return values
}

// Check if a set has member of value.
func (s *Set[T]) Has(value T) bool {
	if s.mux != nil {
		s.mux.RLock()
		defer s.mux.RUnlock()
	}

	_, ok := s.values[value]
	return ok
}

// Add a member.
func (s *Set[T]) Add(value T) {
	if s.mux != nil {
		s.mux.Lock()
		defer s.mux.Unlock()
	}

	s.values[value] = struct{}{}
}

// Remove a member.
func (s *Set[T]) Remove(value T) {
	if s.mux != nil {
		s.mux.Lock()
		defer s.mux.Unlock()
	}

	delete(s.values, value)
}

// Return the union of sets as a new set.
func (s *Set[T]) Union(set *Set[T]) *Set[T] {
	if s.mux != nil {
		s.mux.RLock()
		defer s.mux.RUnlock()
	}

	us := &Set[T]{
		// Note: Size is the minimum possible size of the new set
		values: make(map[T]struct{}, len(s.values)),
	}

	for v, _ := range s.values {
		us.Add(v)
	}

	for v, _ := range set.values {
		us.Add(v)
	}

	return us
}

// Return a shallow copy of the set.
func (s *Set[T]) Copy() *Set[T] {
	if s.mux != nil {
		s.mux.RLock()
		defer s.mux.RUnlock()
	}

	copy := &Set[T]{
		values: make(map[T]struct{}, len(s.values)),
	}

	for v, _ := range s.values {
		copy.Add(v)
	}

	return copy
}

// Returns the difference between two sets.
func (s *Set[T]) Diff(set *Set[T]) *Set[T] {
	if s.mux != nil {
		s.mux.RLock()
		defer s.mux.RUnlock()
	}

	diff := &Set[T]{
		values: make(map[T]struct{}),
	}

	for v, _ := range s.values {
		if !set.Has(v) {
			diff.Add(v)
		}
	}

	return diff
}

// Update a set with the union of itself and set.
func (s *Set[T]) Update(set *Set[T]) {
	if s.mux != nil {
		s.mux.Lock()
		defer s.mux.Unlock()
	}

	for v, _ := range set.values {
		s.Add(v)
	}
}
