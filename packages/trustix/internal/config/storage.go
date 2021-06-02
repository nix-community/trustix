// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package config

import (
	"fmt"
)

type NativeStorage struct {
}

type Storage struct {
	Type   string         `toml:"type" json:"type"`
	Native *NativeStorage `toml:"native" json:"native"`
}

func (s *Storage) Validate() error {
	if s == nil {
		return fmt.Errorf("No storage type configured")
	}

	switch s.Type {
	case "native":
		return nil
	case "memory":
		return nil
	default:
		return fmt.Errorf("Unhandled storage type: '%s'", s.Type)
	}
}
