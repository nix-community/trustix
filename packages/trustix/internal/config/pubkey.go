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
	"fmt"

	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix/internal/protocols"
)

type PublicKey struct {
	Type string `toml:"type" json:"type"`
	Pub  string `toml:"key" json:"key"`
}

func (p *PublicKey) Validate() error {
	if p.Type == "" {
		return fmt.Errorf("Required field 'type' not set")
	}

	if _, err := p.Decode(); err != nil {
		return err
	}
	return nil
}

func (p *PublicKey) Decode() ([]byte, error) {
	return base64.StdEncoding.DecodeString(p.Pub)
}

func (p *PublicKey) LogID(pd *protocols.ProtocolDescriptor, mode api.Log_LogModes) (string, error) {
	pubBytes, err := p.Decode()
	if err != nil {
		return "", err
	}

	return pd.LogID(p.Type, pubBytes, mode), nil
}

func (p *PublicKey) Signer() (*api.LogSigner, error) {
	keyTypeValue, ok := api.LogSigner_KeyTypes_value[p.Type]
	if !ok {
		return nil, fmt.Errorf("Invalid enum value: %s", p.Type)
	}

	return &api.LogSigner{
		KeyType: api.LogSigner_KeyTypes(keyTypeValue).Enum(),
		Public:  &p.Pub,
	}, nil
}
