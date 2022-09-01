package safemap

import (
	"errors"
	"sync"
)

var ErrNotExist = errors.New("item does not exist")

type SafeMap[K comparable, V any] struct {
	store map[K]V
	mux   sync.RWMutex
}

func NewMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		store: make(map[K]V),
		mux:   sync.RWMutex{},
	}
}

func (m *SafeMap[K, V]) Get(key K) (V, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	value, ok := m.store[key]
	if !ok {
		return value, ErrNotExist
	}

	return value, nil
}

func (m *SafeMap[K, V]) Set(key K, value V) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.store[key] = value
}

func (m *SafeMap[K, V]) Has(key K) bool {
	m.mux.RLock()
	defer m.mux.RUnlock()

	_, ok := m.store[key]
	return ok
}

func (m *SafeMap[K, V]) Remove(key K) {
	m.mux.Lock()
	defer m.mux.Unlock()

	delete(m.store, key)
}
