// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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
