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
	Engine     string             `toml:"engine"`
	Lua        *LuaDecider        `toml:"lua"`
	LogID      *LogIDDecider      `toml:"logid"`
	Percentage *PercentageDecider `toml:"percentage"`
}

func (s *Decider) Validate() error {
	switch s.Engine {
	case "lua":
		return s.Lua.Validate()
	case "logid":
		return s.LogID.Validate()
	case "percentage":
		return s.Percentage.Validate()
	default:
		return fmt.Errorf("Unhandled decider engine: '%s'", s.Engine)
	}
}
