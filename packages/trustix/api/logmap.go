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
	"sync"
)

type TrustixLogMap struct {
	logMap map[string]TrustixLogAPI
	mux    *sync.RWMutex
}

func NewTrustixLogMap() *TrustixLogMap {
	return &TrustixLogMap{
		logMap: make(map[string]TrustixLogAPI),
		mux:    &sync.RWMutex{},
	}
}

func (m *TrustixLogMap) Add(logID string, log TrustixLogAPI) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.logMap[logID] = log
}

func (m *TrustixLogMap) Get(logID string) (TrustixLogAPI, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	log, ok := m.logMap[logID]
	if !ok {
		return nil, fmt.Errorf("Missing log '%s'", logID)
	}

	return log, nil
}

func (m *TrustixLogMap) IDs() []string {
	m.mux.RLock()
	defer m.mux.RUnlock()

	keys := make([]string, len(m.logMap))
	i := 0
	for k := range m.logMap {
		keys[i] = k
		i++
	}
	return keys
}

func (m *TrustixLogMap) Map() map[string]TrustixLogAPI {
	m.mux.RLock()
	defer m.mux.RUnlock()

	logMap := make(map[string]TrustixLogAPI)
	for logID, log := range m.logMap {
		logMap[logID] = log
	}

	return logMap
}
