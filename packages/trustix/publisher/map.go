// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package publisher

import (
	"fmt"
	"sync"
)

type PublisherMap struct {
	m   map[string]*Publisher
	mux *sync.RWMutex
}

func NewPublisherMap() *PublisherMap {
	return &PublisherMap{
		m:   make(map[string]*Publisher),
		mux: &sync.RWMutex{},
	}
}

func (pm *PublisherMap) Set(logID string, pub *Publisher) error {
	pm.mux.Lock()
	defer pm.mux.Unlock()

	_, exists := pm.m[logID]
	if exists {
		return fmt.Errorf("Publisher already exists")
	}

	pm.m[logID] = pub

	return nil
}

func (pm *PublisherMap) Get(logID string) (*Publisher, error) {
	pm.mux.RLock()
	defer pm.mux.RUnlock()

	pub, exists := pm.m[logID]
	if !exists {
		return nil, fmt.Errorf("Publisher doesn't exist")
	}

	return pub, nil
}

func (pm *PublisherMap) Close() {
	pm.mux.RLock()
	defer pm.mux.RUnlock()

	for _, pub := range pm.m {
		pub.Close()
	}
}
