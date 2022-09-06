// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

import (
	"fmt"
)

type JSDecider struct {
	Function string `toml:"function" json:"function"`
}

func (s *JSDecider) Validate() error {
	if s.Function == "" {
		return fmt.Errorf("Empty script")
	}
	return nil
}
