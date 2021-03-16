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

package sthmanager

import (
	"fmt"
	"sync"

	"github.com/tweag/trustix/packages/trustix-proto/schema"
)

type STHManager struct {
	logs map[string]STHCache
}

func NewSTHManager() *STHManager {
	return &STHManager{
		logs: make(map[string]STHCache),
	}
}

func (m *STHManager) Add(logName string, c STHCache) {
	m.logs[logName] = c
}

func (m *STHManager) Get(logName string) (*schema.STH, error) {
	cache, ok := m.logs[logName]
	if !ok {
		return nil, fmt.Errorf("Missing log '%s'", logName)
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
