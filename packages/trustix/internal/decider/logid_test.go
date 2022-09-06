// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogIDSimple(t *testing.T) {

	assert := assert.New(t)

	inputs := []*DeciderInput{
		{
			LogID: "test1",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210ab",
		},
		{
			LogID: "test2",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210af",
		},
		{
			LogID: "test3",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210a7",
		},
		{
			LogID: "test4",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
	}

	decider, err := NewLogIDDecider("test3")
	assert.Nil(err)

	output, err := decider.Decide(inputs)
	assert.Nil(err)

	assert.Equal(output.Value, "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210a7", "The correct match is returned")

}

func TestLogIDNonMatch(t *testing.T) {

	assert := assert.New(t)

	inputs := []*DeciderInput{
		{
			LogID: "test1",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210ab",
		},
		{
			LogID: "test2",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210af",
		},
		{
			LogID: "test3",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210a7",
		},
		{
			LogID: "test4",
			Value: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
	}

	decider, err := NewLogIDDecider("test5")
	assert.Nil(err)

	output, err := decider.Decide(inputs)
	assert.Nil(output)
	assert.NotNil(err)

}
