// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package decider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSScript(t *testing.T) {

	script := `function(inputs) {
        return "DummyReturn"
    }`

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

	decider, err := NewJSDecider(script)
	if err != nil {
		t.Log(err) // Make error readable
	}
	assert.Nil(err)

	output, err := decider.Decide(inputs)
	assert.Nil(err)
	assert.NotNil(output)

}
