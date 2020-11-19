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
	"fmt"
	"sort"
)

type minimumPercent struct {
	minimumPct int
}

func NewMinimumPercentDecider(minimumPct int) (LogDecider, error) {
	return &minimumPercent{
		minimumPct: minimumPct,
	}, nil
}

func (q *minimumPercent) Name() string {
	return "percentage"
}

func (q *minimumPercent) Decide(inputs []*LogDeciderInput) (*LogDeciderOutput, error) {
	numInputs := len(inputs)
	pctPerEntry := 100 / numInputs

	// Map OutputHash to list of matches
	entries := make(map[string][]*LogDeciderInput)
	for i := range inputs {
		input := inputs[i]
		l := entries[input.OutputHash]
		l = append(l, input)
		entries[input.OutputHash] = l
	}

	type sortStruct struct {
		key string
		pct int
	}

	makeReturn := func(m *sortStruct) (*LogDeciderOutput, error) {
		ret := &LogDeciderOutput{
			OutputHash: m.key,
			Confidence: m.pct,
		}
		return ret, nil
	}

	// Filter out any keys with less than minimum percentage match and put in list
	var matchesMinimum []*sortStruct
	for k, v := range entries {
		pct := len(v) * pctPerEntry
		if pct >= q.minimumPct {
			m := &sortStruct{
				key: k,
				pct: pct,
			}
			matchesMinimum = append(matchesMinimum, m)
			if pct >= 50 || len(inputs) == 1 || len(entries) == 1 || len(inputs) == len(v) {
				return makeReturn(m)
			}
		}
	}

	switch len(matchesMinimum) {
	case 0:
		return nil, fmt.Errorf("Could not reach the minimum %d quorum", q.minimumPct)
	case 1:
		return makeReturn(matchesMinimum[0])
	default:
		// Sort list by the highest match percentage
		sort.SliceStable(matchesMinimum, func(i, j int) bool {
			return matchesMinimum[i].pct > matchesMinimum[j].pct
		})

		return makeReturn(matchesMinimum[0])
	}
}
