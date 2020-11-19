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
