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

package correlator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercentages(t *testing.T) {

	assert := assert.New(t)

	inputs := []*LogCorrelatorInput{
		&LogCorrelatorInput{
			LogName:    "test1",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&LogCorrelatorInput{
			LogName:    "test2",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&LogCorrelatorInput{
			LogName:    "test3",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&LogCorrelatorInput{
			LogName:    "test4",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
	}

	// 30% minimum matches
	correlator, err := NewMinimumPercentCorrelator(30)
	assert.Nil(err)

	// 100% match case
	output, err := correlator.Decide(inputs)
	assert.Nil(err)
	assert.Equal(len(output.LogNames), 4, "The correct number of matches returned")
	assert.Equal(output.Confidence, 100, "Confidence is correct")

	// 75% match case
	inputs[0].OutputHash = "somedummyvalue"
	output, err = correlator.Decide(inputs)
	assert.Nil(err)
	assert.Equal(3, len(output.LogNames), "The correct number of matches returned")
	assert.Equal(75, output.Confidence, "Confidence is correct")

	// No match case
	inputs[0].OutputHash = "somedummyvalue"
	inputs[1].OutputHash = "someotherdummyvalue"
	inputs[2].OutputHash = "somethirddummyvalue"
	output, err = correlator.Decide(inputs)
	assert.NotNil(err)
	assert.Nil(output)

}
