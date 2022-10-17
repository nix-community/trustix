// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cron

import (
	"context"
	"math/rand"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type CronJob struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type CronFunc = func(context.Context)

// Run fn at an interval
func NewCronJob(name string, d time.Duration, fn CronFunc) *CronJob {
	ctx, cancel := context.WithCancel(context.Background())

	j := &CronJob{
		wg:     sync.WaitGroup{},
		cancel: cancel,
	}

	l := log.WithFields(log.Fields{
		"job":      "cron." + name,
		"interval": d,
	})

	run := func() {
		j.wg.Add(1)
		defer j.wg.Done()

		l.Info("starting job")
		defer l.Info("job done")

		fn(ctx)
	}

	j.wg.Add(1)

	// on the initial run of the cron job add a random sleep within the interval
	// to prevent all concurrent jobs triggering at the same time
	duration := time.Microsecond * time.Duration(rand.Int63n(d.Microseconds()))

	l.WithFields(log.Fields{
		"initial": duration,
	}).Info("initialized job")

	go func() {
		defer j.wg.Done()

		stopChan := ctx.Done()

		for {
			select {
			case <-stopChan:
				l.Info("stopping")
				cancel()
				break
			case <-time.After(duration):
				go run()
			}

			if duration != d {
				duration = d
			}
		}
	}()

	return j
}

func NewSingletonCronJob(name string, d time.Duration, fn CronFunc) *CronJob {
	var mux sync.Mutex
	running := false

	return NewCronJob(name, d, func(ctx context.Context) {
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
	j.cancel()
	j.wg.Wait()

	return nil
}
