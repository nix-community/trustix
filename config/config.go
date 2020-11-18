// MIT License
//
// Copyright (c) 2020 Tweag IO
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package config

import (
	"github.com/BurntSushi/toml"
)

type GitStorageConfig struct {
	Remote   string `toml:"remote"`
	Commiter string `toml:"commiter"`
	Email    string `toml:"email"`
}

type NativeStorageConfig struct {
}

type StorageConfig struct {
	Type   string               `toml:"type"`
	Git    *GitStorageConfig    `toml:"git"`
	Native *NativeStorageConfig `toml:"native"`
}

type GRPCTransportConfig struct {
	Remote string `toml:"remote"`
}

type TransportConfig struct {
	Type string               `toml:"type"`
	GRPC *GRPCTransportConfig `toml:"grpc"`
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
	Name      string           `toml:"name"`
	Mode      string           `toml:"mode"`
	Storage   *StorageConfig   `toml:"storage"`
	Transport *TransportConfig `toml:"transport"`
	Signer    *SignerConfig    `toml:"signer"`
}

type LuaDeciderConfig struct {
	Script string `toml:"script"`
}

type LogNameDeciderConfig struct {
	Name string `toml:"name"`
}

type PercentageDeciderConfig struct {
	Minimum int `toml:"minimum"`
}

type DeciderConfig struct {
	Engine     string                   `toml:"engine"`
	Lua        *LuaDeciderConfig        `toml:"lua"`
	LogName    *LogNameDeciderConfig    `toml:"logname"`
	Percentage *PercentageDeciderConfig `toml:"percentage"`
}

type Config struct {
	Deciders []*DeciderConfig `toml:"decider"`
	Logs     []*LogConfig     `toml:"log"`
}

func NewConfigFromFile(path string) (*Config, error) {
	conf := &Config{}

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}
