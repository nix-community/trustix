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

	// Filter out any keys with less than minimum percentage match and put in list
	var matchesMinimum []*sortStruct
	for k, v := range entries {
		pct := len(v) * pctPerEntry
		if pct >= q.minimumPct {
			matchesMinimum = append(matchesMinimum, &sortStruct{
				key: k,
				pct: pct,
			})
		}
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

	if len(matchesMinimum) == 0 {
		return nil, fmt.Errorf("Could not reach the minimum %d quorum", q.minimumPct)
	}

	// If > 50% only one match is possible, skip sorting
	if q.minimumPct >= 50 {
		return makeReturn(matchesMinimum[0])
	}

	// Sort list by the highest match percentage
	sort.SliceStable(matchesMinimum, func(i, j int) bool {
		return matchesMinimum[i].pct > matchesMinimum[j].pct
	})

	return makeReturn(matchesMinimum[0])
}
