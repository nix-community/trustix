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

type Config struct {
	Logs []*LogConfig `toml:"log"`
}

func NewConfigFromFile(path string) (*Config, error) {
	conf := &Config{}

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}
