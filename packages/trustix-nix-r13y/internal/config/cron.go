// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

type Cron struct {
	LogInterval int64 `toml:"log_index" json:"log_index"`
}

func (c *Cron) init() {
	// Every 15 minutes by default
	if c.LogInterval == 0 {
		c.LogInterval = 15 * 60
	}
}
