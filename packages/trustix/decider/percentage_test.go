// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package decider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercentages(t *testing.T) {

	assert := assert.New(t)

	inputs := []*LogDeciderInput{
		&LogDeciderInput{
			LogID:      "test1",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&LogDeciderInput{
			LogID:      "test2",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&LogDeciderInput{
			LogID:      "test3",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&LogDeciderInput{
			LogID:      "test4",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
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
	inputs[0].OutputHash = "somedummyvalue"
	output, err = decider.Decide(inputs)
	assert.Nil(err)
	assert.Equal(75, output.Confidence, "Confidence is correct")

	// No match case
	inputs[0].OutputHash = "somedummyvalue"
	inputs[1].OutputHash = "someotherdummyvalue"
	inputs[2].OutputHash = "somethirddummyvalue"
	output, err = decider.Decide(inputs)
	assert.NotNil(err)
	assert.Nil(output)

}
