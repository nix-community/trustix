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

func TestAggSimple(t *testing.T) {

	assert := assert.New(t)

	inputs := []*LogDeciderInput{
		{
			LogID:      "test1",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210ab",
		},
		{
			LogID:      "test2",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210af",
		},
		{
			LogID:      "test3",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210a7",
		},
		{
			LogID:      "test4",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
	}

	_decider, err := NewLogIDDecider("test3")
	assert.Nil(err)

	decider := NewAggDecider(_decider)

	output, err := decider.Decide(inputs)
	assert.Nil(err)

	assert.Equal(output.OutputHash, "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210a7", "The correct match is returned")

}

func TestAggNonMatch(t *testing.T) {

	assert := assert.New(t)

	inputs := []*LogDeciderInput{
		{
			LogID:      "test1",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210ab",
		},
		{
			LogID:      "test2",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210af",
		},
		{
			LogID:      "test3",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210a7",
		},
		{
			LogID:      "test4",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
	}

	_decider, err := NewLogIDDecider("test5")
	assert.Nil(err)

	decider := NewAggDecider(_decider)

	output, err := decider.Decide(inputs)
	assert.Nil(output)
	assert.NotNil(err)

}
