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

	"github.com/tweag/trustix/packages/trustix/api"
)

type TrustixCombinedRPCServerMap struct {
	logMap map[string]api.TrustixLogAPI
	mux    sync.Mutex
}

func NewTrustixCombinedRPCServerMap() *TrustixCombinedRPCServerMap {
	return &TrustixCombinedRPCServerMap{
		logMap: make(map[string]api.TrustixLogAPI),
	}
}

func (m *TrustixCombinedRPCServerMap) Add(name string, log api.TrustixLogAPI) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.logMap[name] = log
}

func (m *TrustixCombinedRPCServerMap) Get(name string) (api.TrustixLogAPI, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	log, ok := m.logMap[name]
	if !ok {
		return nil, fmt.Errorf("Missing log '%s'", name)
	}

	return log, nil
}

func (m *TrustixCombinedRPCServerMap) Names() []string {
	m.mux.Lock()
	defer m.mux.Unlock()

	keys := make([]string, len(m.logMap))
	i := 0
	for k := range m.logMap {
		keys[i] = k
		i++
	}
	return keys
}

func (m *TrustixCombinedRPCServerMap) Map() map[string]api.TrustixLogAPI {
	m.mux.Lock()
	defer m.mux.Unlock()

	logMap := make(map[string]api.TrustixLogAPI)
	for name, log := range m.logMap {
		logMap[name] = log
	}

	return logMap
}
