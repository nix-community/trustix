// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package pool

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"

	api "github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix/client"
	log "github.com/sirupsen/logrus"
)

const (
	LocalPriority = -1
	GRPCPriority  = 50
)

type stringSet = map[string]struct{}

type PoolClient struct {
	client   *client.Client
	tags     stringSet
	priority int
	pool     *ClientPool
	mux      *sync.Mutex
	active   bool
}

func (pc *PoolClient) RemoveTags(tags []string) {
	pc.mux.Lock()
	defer pc.mux.Unlock()

	pc.removeTags(tags)
}

func (pc *PoolClient) removeTags(tags []string) {
	for _, tag := range tags {
		delete(pc.tags, tag)
	}

	pc.pool.Remove(pc, tags)
}

func (pc *PoolClient) AddTags(tags []string) {
	pc.mux.Lock()
	defer pc.mux.Unlock()

	pc.addTags(tags)
}

func (pc *PoolClient) addTags(tags []string) {
	for _, tag := range tags {
		pc.tags[tag] = struct{}{}
	}

	pc.pool.addTags(pc, tags)
}

func (pc *PoolClient) SetTags(tags []string) {
	pc.mux.Lock()
	defer pc.mux.Unlock()

	removedTags := []string{}

	newTags := make(stringSet)
	for _, tag := range tags {
		newTags[tag] = struct{}{}
	}

	for tag := range pc.tags {
		_, ok := newTags[tag]
		if !ok {
			removedTags = append(removedTags, tag)
		}
	}

	pc.addTags(tags)
	pc.removeTags(removedTags)
}

func (pc *PoolClient) Activate() {
	pc.pool.mux.Lock()
	defer pc.pool.mux.Unlock()

	pc.active = true

	for tag := range pc.tags {
		current, ok := pc.pool.active[tag]

		// If a currently active set exists but priority is lower than current client don't mutate active set
		if ok && pc.priority < current.priority {
			continue
		}

		if ok && current.priority == pc.priority {
			current.set[pc] = struct{}{}
			continue
		}

		pc.pool.recalculateActive(tag)
	}
}

func (pc *PoolClient) Deactivate() {
	pc.pool.mux.Lock()
	defer pc.pool.mux.Unlock()

	pc.active = false

	for tag := range pc.tags {
		current, ok := pc.pool.active[tag]
		if !ok {
			continue
		}
		delete(current.set, pc)

		// There are still active connections, don't recalculate
		if len(current.set) >= 1 {
			continue
		}

		if len(current.set) == 0 {
			delete(pc.pool.active, tag)
		}

		pc.pool.recalculateActive(tag)
	}
}

type poolClientSet = map[*PoolClient]struct{}

type activePoolClients struct {
	priority int
	set      poolClientSet
}

type ClientPool struct {
	// Complete set
	all poolClientSet

	// map[tag][priority]{*PoolClient}
	clients map[string]map[int]poolClientSet

	// map[tag]{*PoolClient}
	active map[string]*activePoolClients

	mux *sync.RWMutex
}

func NewClientPool() *ClientPool {
	return &ClientPool{
		all:     make(poolClientSet),
		clients: make(map[string]map[int]poolClientSet),
		active:  make(map[string]*activePoolClients),
		mux:     &sync.RWMutex{},
	}
}

func (p *ClientPool) Add(c *client.Client, tags []string) (*PoolClient, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if tags == nil {
		tags = []string{}
	}

	var priority int
	switch c.Type {
	case client.LocalClientType:
		priority = LocalPriority
	case client.GRPCClientType:
		priority = GRPCPriority
	default:
		return nil, fmt.Errorf("Unhandled client type: %d", c.Type)
	}

	tagMap := make(stringSet)
	for _, tag := range tags {
		tagMap[tag] = struct{}{}
	}

	pc := &PoolClient{
		client:   c,
		tags:     tagMap,
		priority: priority,
		pool:     p,
		mux:      &sync.Mutex{},
		active:   false,
	}

	p.addTags(pc, tags)

	p.all[pc] = struct{}{}

	return pc, nil
}

func (p *ClientPool) Dial(address string) (*PoolClient, error) {

	// TODO: Retry loop (probably managed by PoolClient somehow)
	c, err := client.CreateClient(address)
	if err != nil {
		return nil, err
	}

	pc, err := p.Add(c, nil)
	if err != nil {
		return nil, err
	}

	// TODO: This responsibility belongs elsewhere
	go func() {
		ctx, _ := client.CreateContext(30)
		resp, err := c.NodeAPI.Logs(ctx, &api.LogsRequest{})
		if err != nil {
			log.Errorf("Error while getting logs from '%s': %v", address, err)
		}

		logIDs := make([]string, len(resp.Logs))
		for i, log := range resp.Logs {
			logIDs[i] = *log.LogID
		}

		pc.AddTags(logIDs)
	}()

	return pc, nil
}

func (p *ClientPool) addTags(client *PoolClient, tags []string) {

	for _, tag := range tags {
		priorities, ok := p.clients[tag]
		if !ok {
			priorities = make(map[int]poolClientSet)
			p.clients[tag] = priorities
		}

		clients, ok := priorities[client.priority]
		if !ok {
			clients = make(poolClientSet)
			priorities[client.priority] = clients
		}

		clients[client] = struct{}{}

		if client.active {
			p.recalculateActive(tag)
		}
	}

}

func (p *ClientPool) Get(tag string) (*client.Client, error) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	active, ok := p.active[tag]
	if !ok {
		return nil, fmt.Errorf("Couldn't find client for tag: %s", tag)
	}

	// If there is only one client simply return it
	if len(active.set) == 1 {
		for pc := range active.set {
			return pc.client, nil
		}
		panic("Programming error")
	}

	// Otherwise select a random one out of the set
	clients := make([]*PoolClient, len(active.set))
	{
		i := 0
		for client := range active.set {
			clients[i] = client
			i++
		}
	}

	idx := rand.Intn(len(clients))
	return clients[idx].client, nil
}

func (p *ClientPool) recalculateActive(tag string) {

	// Iterate active (in priority order) until active connections can be found
	priorities, ok := p.clients[tag]
	if !ok {
		delete(p.active, tag)
		return
	}

	prios := make([]int, len(priorities))
	{
		i := 0
		for prio := range priorities {
			prios[i] = prio
			i++
		}
		sort.Ints(prios)
	}

	clients := poolClientSet{}
	for _, prio := range prios {
		for client := range priorities[prio] {
			if client.active {
				clients[client] = struct{}{}
			}
		}

		if len(clients) > 0 {
			p.active[tag] = &activePoolClients{
				priority: prio,
				set:      clients,
			}
			return
		}
	}

	// Nothing reached
	delete(p.active, tag)
}

func (p *ClientPool) Remove(client *PoolClient, tags []string) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	if tags == nil {
		tags = []string{}
		for tag := range client.tags {
			tags = append(tags, tag)
		}

		delete(p.all, client)
	}

	// Clear out from p.clients
	for _, tag := range tags {
		priorities, ok := p.clients[tag]
		if !ok {
			continue
		}

		clients, ok := priorities[client.priority]
		if !ok {
			continue
		}

		delete(clients, client)

		// Clean up orphaned maps

		if len(clients) == 0 {
			delete(clients, client)
		}

		if len(priorities) == 0 {
			delete(p.clients, tag)
		}
	}

	// Clear out p.active
	for _, tag := range tags {
		clients, ok := p.active[tag]
		if !ok {
			continue
		}
		delete(clients.set, client)

		if len(clients.set) == 0 {
			p.recalculateActive(tag)
		}
	}

}

func (p *ClientPool) Close() {
	wg := &sync.WaitGroup{}

	for client := range p.all {
		wg.Add(1)
		client := client

		go func() {
			defer wg.Done()
			err := client.client.Close()
			if err != nil {
				log.Errorf("Error while closing client: %v", err)
			}
		}()

	}

	wg.Wait()
}
