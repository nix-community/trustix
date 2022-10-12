// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: MIT

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
	for v := range s.values {
		values[i] = v
		i++
	}

	slices.Sort(values)

	return values
}

// Check if a set has member of value.
func (s *Set[T]) has(value T) bool {
	_, ok := s.values[value]
	return ok
}

// Check if a set has member of value.
func (s *Set[T]) Has(value T) bool {
	if s.mux != nil {
		s.mux.RLock()
		defer s.mux.RUnlock()
	}

	return s.has(value)
}

func (s *Set[T]) add(value T) {
	s.values[value] = struct{}{}
}

// Add a member.
func (s *Set[T]) Add(value T) (added bool) {
	if s.mux != nil {
		s.mux.Lock()
		defer s.mux.Unlock()
	}

	if s.has(value) {
		return false
	}

	s.add(value)

	return true
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
	if set.mux != nil {
		set.mux.RLock()
		defer set.mux.RUnlock()
	}

	us := &Set[T]{
		// Note: Size is the minimum possible size of the new set
		values: make(map[T]struct{}, len(s.values)),
	}

	for v := range s.values {
		us.add(v)
	}

	for v := range set.values {
		us.add(v)
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

	for v := range s.values {
		copy.add(v)
	}

	return copy
}

// Returns the difference between two sets.
func (s *Set[T]) Diff(set *Set[T]) *Set[T] {
	if s.mux != nil {
		s.mux.RLock()
		defer s.mux.RUnlock()
	}
	if set.mux != nil {
		set.mux.RLock()
		defer set.mux.RUnlock()
	}

	diff := &Set[T]{
		values: make(map[T]struct{}),
	}

	for v := range s.values {
		if !set.has(v) {
			diff.add(v)
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
	if set.mux != nil {
		set.mux.RLock()
		defer set.mux.RUnlock()
	}

	for v := range set.values {
		s.add(v)
	}
}
