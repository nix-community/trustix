// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cron

import (
	"context"
	"sync"
	"time"
)

type CronJob struct {
	stopChan chan struct{}
	wg       sync.WaitGroup
}

type CronFunc = func(context.Context)

// Run fn at an interval
func NewCronJob(d time.Duration, fn CronFunc) *CronJob {
	j := &CronJob{
		stopChan: make(chan struct{}),
		wg:       sync.WaitGroup{},
	}

	ctx, cancel := context.WithCancel(context.Background())

	run := func() {
		j.wg.Add(1)
		defer j.wg.Done()

		fn(ctx)
	}

	j.wg.Add(1)

	go func() {
		defer j.wg.Done()

		go run()

		for {
			select {
			case <-j.stopChan:
				cancel()
				break
			case <-time.After(d):
				go run()
			}
		}
	}()

	return j
}

func NewSingletonCronJob(d time.Duration, fn CronFunc) *CronJob {
	var mux sync.Mutex
	running := false

	return NewCronJob(d, func(ctx context.Context) {
		// Return if already running
		{
			mux.Lock()

			if running {
				return
			}

			running = true

			mux.Unlock()
		}

		// Unset running
		defer func() {
			mux.Lock()
			running = false
			mux.Unlock()
		}()

		fn(ctx)
	})
}

func (j *CronJob) Close() error {
	j.stopChan <- struct{}{}

	j.wg.Wait()

	return nil
}
