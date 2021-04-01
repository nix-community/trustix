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
)

type lognameDecider struct {
	logName string
}

// NewLogNameDecider - Matches a single name of a log
// This mode is meant to be used as a fallback for e.g. cache.nixos.org
func NewLogNameDecider(logName string) (LogDecider, error) {
	return &lognameDecider{
		logName: logName,
	}, nil
}

func (l *lognameDecider) Name() string {
	return "logname"
}

func (l *lognameDecider) Decide(inputs []*LogDeciderInput) (*LogDeciderOutput, error) {
	for i := range inputs {
		input := inputs[i]
		if input.LogName == l.logName {
			return &LogDeciderOutput{
				OutputHash: input.OutputHash,
				Confidence: 100,
			}, nil
		}
	}

	return nil, fmt.Errorf("Could not find any match for log name %s", l.logName)
}
