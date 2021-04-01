// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package decider

// LogDeciderInput
type LogDeciderInput struct {
	LogName    string
	OutputHash string
}

type LogDeciderOutput struct {
	// The decided OutputHash
	OutputHash string
	// An arbitrary number conveying the underlying engines confidence in the result
	Confidence int
}

type LogDecider interface {
	// Name - The name of the decision engine
	Name() string

	// Decide - Decide on an output hash from aggregated logs
	Decide([]*LogDeciderInput) (*LogDeciderOutput, error)
}
