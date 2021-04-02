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
)

type NativeStorageConfig struct {
}

type StorageConfig struct {
	Type   string               `toml:"type"`
	Native *NativeStorageConfig `toml:"native"`
}

type Config struct {
	Deciders    []*DeciderConfig    `toml:"decider"`
	Publishers  []*PublisherConfig  `toml:"publisher"`
	Subscribers []*SubscriberConfig `toml:"subscriber"`
	Storage     *StorageConfig      `toml:"storage"`
}

func NewConfigFromFile(path string) (*Config, error) {
	conf := &Config{}

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}
