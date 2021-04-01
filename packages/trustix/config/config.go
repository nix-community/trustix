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

type ED25519SignerConfig struct {
	PrivateKeyPath string `toml:"private-key-path"`
}

type SignerConfig struct {
	Type      string               `toml:"type"`
	KeyType   string               `toml:"key-type"`
	PublicKey string               `toml:"public-key"`
	ED25519   *ED25519SignerConfig `toml:"ed25519"`
}

type LogConfig struct {
	// Name    string            `toml:"name"`
	Mode    string            `toml:"mode"`
	Storage *StorageConfig    `toml:"storage"`
	Signer  *SignerConfig     `toml:"signer"`
	Meta    map[string]string `toml:"meta"`
}

type Config struct {
	Deciders    []*DeciderConfig    `toml:"decider"`
	Logs        []*LogConfig        `toml:"log"`
	Subscribers []*SubscriberConfig `toml:"subscriber"`
}

func NewConfigFromFile(path string) (*Config, error) {
	conf := &Config{}

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}
