// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package rpc

import (
	"fmt"
	"sync"

	pb "github.com/tweag/trustix/packages/trustix-proto/rpc"
)

type SignerMetaMap struct {
	m   map[string]*pb.LogSigner
	mux *sync.RWMutex
}

func NewSignerMetaMap() *SignerMetaMap {
	return &SignerMetaMap{
		m:   make(map[string]*pb.LogSigner),
		mux: &sync.RWMutex{},
	}
}

func (s *SignerMetaMap) Set(logID string, keyType string, pub string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	keyTypeValue, ok := pb.LogSigner_KeyTypes_value[keyType]
	if !ok {
		return fmt.Errorf("Invalid enum value: %s", keyType)
	}

	s.m[logID] = &pb.LogSigner{
		KeyType: pb.LogSigner_KeyTypes(keyTypeValue).Enum(),
		Public:  &pub,
	}

	return nil
}

func (s *SignerMetaMap) Get(logID string) (*pb.LogSigner, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	sig, ok := s.m[logID]
	if !ok {
		return nil, fmt.Errorf("Missing logID: %s", logID)
	}

	return sig, nil
}
