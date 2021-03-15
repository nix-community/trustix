// MIT License
//
// Copyright (c) 2020 Tweag IO
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

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
