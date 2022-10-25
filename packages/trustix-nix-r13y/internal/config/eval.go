// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

type Eval struct {
	Workers int `toml:"workers" json:"workers"`
	Jobs    int `toml:"jobs" json:"jobs"`
}
