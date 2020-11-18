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

type aggDecider struct {
	deciders []LogDecider
}

// Create a new logdecider that steps over a slice of deciders one by one,
// returning the first match, and returns an aggregated error if no decision could be made
func NewAggDecider(deciders ...LogDecider) LogDecider {
	return &aggDecider{
		deciders: deciders,
	}
}

func (d *aggDecider) Name() string {
	return "aggregate"
}

func (d *aggDecider) Decide(inputs []*LogDeciderInput) (*LogDeciderOutput, error) {
	errors := make([]error, len(d.deciders))
	for i, decider := range d.deciders {
		decision, err := decider.Decide(inputs)
		if err != nil {
			errors[i] = err
			continue
		}
		return decision, nil
	}

	errorS := "Encountered errors while deciding:\n"
	for i, err := range errors {
		if err == nil {
			continue
		}

		decider := d.deciders[i]

		errorS = errorS + decider.Name() + ": " + err.Error() + "\n"
	}
	return nil, fmt.Errorf(errorS)
}
