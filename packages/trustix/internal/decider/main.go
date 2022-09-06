// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

// DeciderInput
type DeciderInput struct {
	LogID string
	Value string
}

type DeciderOutput struct {
	// The decided Value
	Value string
	// An arbitrary number conveying the underlying engines confidence in the result
	Confidence int
}

type LogDecider interface {
	// Name - The name of the decision engine
	Name() string

	// Decide - Decide on an output hash from aggregated logs
	Decide([]*DeciderInput) (*DeciderOutput, error)
}
