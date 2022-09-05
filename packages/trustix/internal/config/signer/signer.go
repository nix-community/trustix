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

type ED25519 struct {
	PrivateKeyPath string `toml:"private-key-path" json:"private-key-path"`
}

func (s *ED25519) Validate() error {
	if s.PrivateKeyPath == "" {
		return fmt.Errorf("Empty private key path")
	}
	return nil
}

type Signer struct {
	Type    string   `toml:"type" json:"type"`
	ED25519 *ED25519 `toml:"ed25519" json:"ed25519"`
}

func (s *Signer) Validate() error {
	switch s.Type {
	case "ed25519":
		return s.ED25519.Validate()
	default:
		return fmt.Errorf("Unhandled signer type: %s", s.Type)
	}
}
