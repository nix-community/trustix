// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package sthmanager

import (
	"fmt"
	"sync"

	"github.com/tweag/trustix/packages/trustix-proto/schema"
)

type STHManager struct {
	logs map[string]STHCache
	mux  *sync.RWMutex
}

func NewSTHManager() *STHManager {
	return &STHManager{
		logs: make(map[string]STHCache),
		mux:  &sync.RWMutex{},
	}
}

func (m *STHManager) Set(logID string, c STHCache) {
	m.mux.Lock()
	m.logs[logID] = c
	m.mux.Unlock()
}

func (m *STHManager) Get(logID string) (*schema.STH, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	cache, ok := m.logs[logID]
	if !ok {
		return nil, fmt.Errorf("Missing log '%s'", logID)
	}

	return cache.Get()
}

func (m *STHManager) Close() {
	wg := new(sync.WaitGroup)

	for _, c := range m.logs {
		wg.Add(1)
		go func() {
			c.Close()
			wg.Done()
		}()
	}

	wg.Wait()
}
