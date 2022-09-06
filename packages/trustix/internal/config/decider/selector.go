// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

import (
	"fmt"
)

type Decider struct {
	Engine     string             `toml:"engine" json:"engine"`
	JS         *JSDecider         `toml:"javascript" json:"javascript"`
	LogID      *LogIDDecider      `toml:"logid" json:"logid"`
	Percentage *PercentageDecider `toml:"percentage" json:"percentage"`
}

func (s *Decider) Validate() error {
	switch s.Engine {
	case "javascript":
		return s.JS.Validate()
	case "logid":
		return s.LogID.Validate()
	case "percentage":
		return s.Percentage.Validate()
	default:
		return fmt.Errorf("Unhandled decider engine: '%s'", s.Engine)
	}
}
