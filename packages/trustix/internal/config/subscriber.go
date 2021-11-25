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

type Subscriber struct {
	Protocol  string            `toml:"protocol" json:"protocol"`
	PublicKey *PublicKey        `toml:"publicKey" json:"publicKey"`
	SyncMode  string            `toml:"syncmode" json:"syncmode"`
	Meta      map[string]string `toml:"meta" json:"meta"`
}

func (s *Subscriber) Validate() error {
	if s.Protocol == "" {
		return missingField("protocol")
	}
	if s.SyncMode != "" && s.SyncMode != "light" {
		return fmt.Errorf("Unknown sync mode: '%s'", s.SyncMode)
	}

	return nil
}

func (s *Subscriber) GetMeta() map[string]string {
	if s.Meta != nil {
		return s.Meta
	}
	return make(map[string]string)
}
