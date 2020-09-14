package main

import (
	"github.com/BurntSushi/toml"
)

type GitStorageConfig struct {
	Remote   string `toml:"remote"`
	Path     string `toml:"path"`
	Commiter string `toml:"commiter"`
	Email    string `toml:"email"`
}

type StorageConfig struct {
	Type string            `toml:"type"`
	Git  *GitStorageConfig `toml:"git"`
}

type ED25519SignerConfig struct {
	PrivateKeyPath string `toml:"private-key-path"`
}

type SignerConfig struct {
	Type      string `toml:"type"`
	KeyType   string `toml:"key-type"`
	PublicKey string `toml:"public-key"`
}

type LogConfig struct {
	Mode    string         `toml:"mode"`
	Storage *StorageConfig `toml:"storage"`
}

type Config struct {
	Logs []*LogConfig `toml:"log"`
}

func newConfig(path string) (*Config, error) {
	conf := &Config{}

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}
