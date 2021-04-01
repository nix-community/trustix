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

func TestLuaScript(t *testing.T) {

	script := `
      function(inputs)
        t = {}
        t["LogNames"] = {"DummyLogName"}
        t["OutputHash"] = "DummyReturn"
        return t
      end
    `

	assert := assert.New(t)

	inputs := []*LogDeciderInput{
		&LogDeciderInput{
			LogName:    "test1",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&LogDeciderInput{
			LogName:    "test2",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&LogDeciderInput{
			LogName:    "test3",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
		&LogDeciderInput{
			LogName:    "test4",
			OutputHash: "26c499a911e8376c52940e050cecc7fc1b9699e759d18856323391c82a2210aa",
		},
	}

	decider, err := NewLuaDecider(script)
	assert.Nil(err)

	output, err := decider.Decide(inputs)
	assert.Nil(err)
	assert.NotNil(output)

}
