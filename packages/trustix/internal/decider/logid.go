// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

import (
	"fmt"
)

type logIDDecider struct {
	logID string
}

// NewLogIDDecider - Matches entry from a single log id
func NewLogIDDecider(logID string) (LogDecider, error) {
	return &logIDDecider{
		logID: logID,
	}, nil
}

func (d *logIDDecider) Name() string {
	return "logID"
}

func (d *logIDDecider) Decide(inputs []*DeciderInput) (*DeciderOutput, error) {
	for i := range inputs {
		input := inputs[i]
		if input.LogID == d.logID {
			return &DeciderOutput{
				Value:      input.Value,
				Confidence: 100,
			}, nil
		}
	}

	return nil, fmt.Errorf("Could not find any match for log name %s", d.logID)
}
