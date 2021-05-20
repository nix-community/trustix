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

	signer "github.com/tweag/trustix/packages/trustix/internal/config/signer"
)

type Publisher struct {
	Signer    string            `toml:"signer"`
	PublicKey *PublicKey        `toml:"key"`
	Meta      map[string]string `toml:"meta"`
}

func (p *Publisher) Validate(signers map[string]*signer.Signer) error {
	if p.Signer == "" {
		return missingField("signer")
	}

	if err := p.PublicKey.Validate(); err != nil {
		return err
	}

	_, ok := signers[p.Signer]
	if !ok {
		return fmt.Errorf("Signer '%s' referenced but does not exist", p.Signer)
	}

	return nil
}

func (p *Publisher) GetMeta() map[string]string {
	if p.Meta != nil {
		return p.Meta
	}
	return make(map[string]string)
}
