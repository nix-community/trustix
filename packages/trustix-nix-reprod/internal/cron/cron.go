package cron

import (
	"sync/atomic"
	"time"
)

type CronJob struct {
	stopChan chan struct{}
	ticker   *time.Ticker
}

// Run fn at an interval
func NewCronJob(d time.Duration, fn func()) *CronJob {
	j := &CronJob{
		ticker:   time.NewTicker(d),
		stopChan: make(chan struct{}),
	}

	go func() {
		go fn()

		for {
			select {
			case <-j.stopChan:
				break
			case <-j.ticker.C:
				go fn()
			}
		}
	}()

	return j
}

// A variant of a cronjob which only ever runs one instance of a function
// if another instance is already running, the tick is dropped.
func NewSingletonCronJob(d time.Duration, fn func()) *CronJob {
	j := &CronJob{
		ticker:   time.NewTicker(d),
		stopChan: make(chan struct{}),
	}

	var running atomic.Bool

	run := func() {
		if running.Load() {
			return
		}

		running.Store(true)
		defer func() {
			running.Store(false)
		}()

		fn()
	}

	go func() {
		go run()

		for {
			select {
			case <-j.stopChan:
				break
			case <-j.ticker.C:
				go run()
			}
		}
	}()

	return j
}

func (j *CronJob) Stop() {
	j.stopChan <- struct{}{}

	j.ticker.Stop()
}
