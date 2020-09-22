package correlator

import (
	"fmt"
)

type lognameCorrelator struct {
	logName string
}

// NewLogNameCorrelator - Matches a single name of a log
// This mode is meant to be used as a fallback for e.g. cache.nixos.org
func NewLogNameCorrelator(logName string) LogCorrelator {
	return &lognameCorrelator{
		logName: logName,
	}
}

func (l *lognameCorrelator) Decide(inputs []*LogCorrelatorInput) (*LogCorrelatorOutput, error) {
	for i := range inputs {
		input := inputs[i]
		if input.LogName == l.logName {
			return &LogCorrelatorOutput{
				LogNames:   []string{input.LogName},
				OutputHash: input.OutputHash,
			}, nil
		}
	}

	return nil, fmt.Errorf("Could not find any match for log name %s", l.logName)
}
