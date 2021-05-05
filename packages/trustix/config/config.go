// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package config

import (
	"github.com/BurntSushi/toml"

	decider "github.com/tweag/trustix/packages/trustix/config/decider"
	signer "github.com/tweag/trustix/packages/trustix/config/signer"
)

type Config struct {
	Deciders    []*decider.Decider        `toml:"decider"`
	Publishers  []*Publisher              `toml:"publisher"`
	Subscribers []*Subscriber             `toml:"subscriber"`
	Signers     map[string]*signer.Signer `toml:"signer"`
	Storage     *Storage                  `toml:"storage"`
}

func (c *Config) Validate() error {

	if err := c.Storage.Validate(); err != nil {
		return err
	}

	for _, signer := range c.Signers {
		if err := signer.Validate(); err != nil {
			return err
		}
	}

	for _, publisher := range c.Publishers {
		if err := publisher.Validate(c.Signers); err != nil {
			return err
		}
	}

	for _, subscriber := range c.Subscribers {
		if err := subscriber.Validate(); err != nil {
			return err
		}
	}

	for _, decider := range c.Deciders {
		if err := decider.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func NewConfigFromFile(path string) (*Config, error) {
	conf := &Config{}

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return conf, nil
}
