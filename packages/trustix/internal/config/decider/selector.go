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

type Decider struct {
	Engine     string             `toml:"engine" json:"engine"`
	JS         *JSDecider         `toml:"javascript" json:"javascript"`
	LogID      *LogIDDecider      `toml:"logid" json:"logid"`
	Percentage *PercentageDecider `toml:"percentage" json:"percentage"`
}

func (s *Decider) Validate() error {
	switch s.Engine {
	case "javascript":
		return s.JS.Validate()
	case "logid":
		return s.LogID.Validate()
	case "percentage":
		return s.Percentage.Validate()
	default:
		return fmt.Errorf("Unhandled decider engine: '%s'", s.Engine)
	}
}
