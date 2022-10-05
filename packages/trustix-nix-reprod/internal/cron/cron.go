package cron

import (
	"time"
)

type CronJob struct {
	stopChan chan struct{}
	ticker   *time.Ticker
}

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

func (j *CronJob) Stop() {
	j.stopChan <- struct{}{}

	j.ticker.Stop()
}
