package core

import (
	"fmt"
	"github.com/tweag/trustix/storage"
)

var noCurrentTransaction = fmt.Errorf("No current transaction")

type smtMapStore struct {
	txn   storage.Transaction
	inTxn bool
}

// Implement MapStore for SMT lib
func newMapStore() *smtMapStore {
	return &smtMapStore{}
}

func (s *smtMapStore) setTxn(txn storage.Transaction) {
	s.txn = txn
	s.inTxn = true
}

func (s *smtMapStore) unsetTxn() {
	s.txn = nil
	s.inTxn = false
}

func (s *smtMapStore) Get(key []byte) ([]byte, error) {
	if !s.inTxn {
		return nil, noCurrentTransaction
	}

	return s.txn.Get(key)
}

func (s *smtMapStore) Set(key []byte, value []byte) error {
	if !s.inTxn {
		return noCurrentTransaction
	}

	return s.txn.Set(key, value)
}

func (s *smtMapStore) Delete(key []byte) error {
	return fmt.Errorf("Delete unsupported")
}
