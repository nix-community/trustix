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

type logIDDecider struct {
	logID string
}

// NewLogIDDecider - Matches entry from a single log id
func NewLogIDDecider(logID string) (LogDecider, error) {
	return &logIDDecider{
		logID: logID,
	}, nil
}

func (d *logIDDecider) Name() string {
	return "logID"
}

func (d *logIDDecider) Decide(inputs []*LogDeciderInput) (*LogDeciderOutput, error) {
	for i := range inputs {
		input := inputs[i]
		if input.LogID == d.logID {
			return &LogDeciderOutput{
				OutputHash: input.OutputHash,
				Confidence: 100,
			}, nil
		}
	}

	return nil, fmt.Errorf("Could not find any match for log name %s", d.logID)
}
