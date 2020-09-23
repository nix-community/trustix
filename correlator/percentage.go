package correlator

import (
	"fmt"
	"sort"
)

type minimumPercent struct {
	minimumPct int
}

func NewMinimumPercentCorrelator(minimumPct int) LogCorrelator {
	return &minimumPercent{
		minimumPct: minimumPct,
	}
}

func (q *minimumPercent) Decide(inputs []*LogCorrelatorInput) (*LogCorrelatorOutput, error) {
	numInputs := len(inputs)
	pctPerEntry := 100 / numInputs

	// Map OutputHash to list of matches
	entries := make(map[string][]*LogCorrelatorInput)
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

	makeReturn := func(m *sortStruct) (*LogCorrelatorOutput, error) {
		ret := &LogCorrelatorOutput{
			OutputHash: m.key,
		}
		for _, v := range entries[m.key] {
			ret.LogNames = append(ret.LogNames, v.LogName)
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
