// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercentages(t *testing.T) {

	assert := assert.New(t)

	inputs := []*DeciderInput{
		&DeciderInput{
			LogID: "test1",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&DeciderInput{
			LogID: "test2",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&DeciderInput{
			LogID: "test3",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&DeciderInput{
			LogID: "test4",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
	}

	// 30% minimum matches
	decider, err := NewMinimumPercentDecider(30)
	assert.Nil(err)

	// 100% match case
	output, err := decider.Decide(inputs)
	assert.Nil(err)
	assert.Equal(output.Confidence, 100, "Confidence is correct")

	// 75% match case
	inputs[0].Value = "somedummyvalue"
	output, err = decider.Decide(inputs)
	assert.Nil(err)
	assert.Equal(75, output.Confidence, "Confidence is correct")

	// No match case
	inputs[0].Value = "somedummyvalue"
	inputs[1].Value = "someotherdummyvalue"
	inputs[2].Value = "somethirddummyvalue"
	output, err = decider.Decide(inputs)
	assert.NotNil(err)
	assert.Nil(output)

}
