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
