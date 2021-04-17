// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package config

import (
	"encoding/base64"
)

type PublicKey struct {
	Type string `toml:"type"`
	Pub  string `toml:"pub"`
}

func (p *PublicKey) Validate() error {
	if _, err := p.Decode(); err != nil {
		return err
	}
	return nil
}

func (p *PublicKey) Decode() ([]byte, error) {
	return base64.StdEncoding.DecodeString(p.Pub)
}
