// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package api

import (
	"fmt"

	"github.com/tweag/trustix/packages/trustix/storage"
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
