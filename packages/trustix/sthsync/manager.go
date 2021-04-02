// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package sthsync

import (
	"sync"
)

type STHSyncer interface {
	Close()
}

type SyncManager struct {
	mux   *sync.Mutex
	syncs []STHSyncer
}

func NewSyncManager() *SyncManager {
	return &SyncManager{
		syncs: []STHSyncer{},
		mux:   &sync.Mutex{},
	}
}

func (m *SyncManager) Add(syncer STHSyncer) {
	m.mux.Lock()
	m.syncs = append(m.syncs, syncer)
	m.mux.Unlock()
}

func (m *SyncManager) Close() {
	wg := new(sync.WaitGroup)

	for _, s := range m.syncs {
		wg.Add(1)
		go func() {
			s.Close()
			wg.Done()
		}()
	}

	wg.Wait()
}
