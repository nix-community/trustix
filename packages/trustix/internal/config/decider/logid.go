// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

import (
	"fmt"
)

type LogIDDecider struct {
	ID string `toml:"id" json:"id"`
}

func (s *LogIDDecider) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("Empty log ids are invalid")
	}
	return nil
}
