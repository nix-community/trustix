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
)

type SMTMapStore struct {
	logID      string
	storageAPI *StorageAPI
}

func newMapStore(storageAPI *StorageAPI, logID string) *SMTMapStore {
	return &SMTMapStore{
		logID:      logID,
		storageAPI: storageAPI,
	}
}

func (s *SMTMapStore) Get(key []byte) ([]byte, error) {
	return s.storageAPI.GetSMTValue(s.logID, key)
}

func (s *SMTMapStore) Set(key []byte, value []byte) error {
	return s.storageAPI.SetSMTValue(s.logID, key, value)
}

func (s *SMTMapStore) Delete(key []byte) error {
	return fmt.Errorf("Delete unsupported")
}
