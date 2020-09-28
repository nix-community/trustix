package core

import (
	"fmt"
	"github.com/tweag/trustix/storage"
)

type smtMapStore struct {
	txn storage.Transaction
}

// Implement MapStore for SMT lib
func newMapStore(txn storage.Transaction) *smtMapStore {
	return &smtMapStore{
		txn: txn,
	}
}

func (s *smtMapStore) Get(key []byte) ([]byte, error) {
	return s.txn.Get([]byte("SMT"), key)
}

func (s *smtMapStore) Set(key []byte, value []byte) error {
	return s.txn.Set([]byte("SMT"), key, value)
}

func (s *smtMapStore) Delete(key []byte) error {
	return fmt.Errorf("Delete unsupported")
}
