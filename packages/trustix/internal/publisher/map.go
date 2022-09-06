// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

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
		return nil, fmt.Errorf("Publisher with log id '%s' doesn't exist", logID)
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
