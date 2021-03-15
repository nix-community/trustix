// MIT License
//
// Copyright (c) 2020 Tweag IO
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package decider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggSimple(t *testing.T) {

	assert := assert.New(t)

	inputs := []*LogDeciderInput{
		{
			LogName:    "test1",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210ab",
		},
		{
			LogName:    "test2",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210af",
		},
		{
			LogName:    "test3",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210a7",
		},
		{
			LogName:    "test4",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
	}

	_decider, err := NewLogNameDecider("test3")
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
			LogName:    "test1",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210ab",
		},
		{
			LogName:    "test2",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210af",
		},
		{
			LogName:    "test3",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210a7",
		},
		{
			LogName:    "test4",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
	}

	_decider, err := NewLogNameDecider("test5")
	assert.Nil(err)

	decider := NewAggDecider(_decider)

	output, err := decider.Decide(inputs)
	assert.Nil(output)
	assert.NotNil(err)

}
