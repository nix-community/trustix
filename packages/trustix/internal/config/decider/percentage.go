// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

import (
	"fmt"
)

type PercentageDecider struct {
	Minimum int `toml:"minimum" json:"minimum"`
}

func (s *PercentageDecider) Validate() error {
	if s.Minimum < 0 || s.Minimum > 100 {
		return fmt.Errorf("Minimum percentage decider out of bounds (%d)", s.Minimum)
	}
	return nil
}
