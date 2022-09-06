// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

import (
	"fmt"
)

type aggDecider struct {
	deciders []LogDecider
}

// Create a new logdecider that steps over a slice of deciders one by one,
// returning the first match, and returns an aggregated error if no decision could be made
func NewAggDecider(deciders ...LogDecider) LogDecider {
	return &aggDecider{
		deciders: deciders,
	}
}

func (d *aggDecider) Name() string {
	return "aggregate"
}

func (d *aggDecider) Decide(inputs []*DeciderInput) (*DeciderOutput, error) {
	if len(d.deciders) == 0 {
		return nil, fmt.Errorf("No decision making engines configured")
	}

	errors := make([]error, len(d.deciders))
	for i, decider := range d.deciders {
		decision, err := decider.Decide(inputs)
		if err != nil {
			errors[i] = err
			continue
		}
		return decision, nil
	}

	errorS := "Encountered errors while deciding:\n"
	for i, err := range errors {
		if err == nil {
			continue
		}

		decider := d.deciders[i]

		errorS = errorS + decider.Name() + ": " + err.Error() + "\n"
	}
	return nil, fmt.Errorf(errorS)
}
