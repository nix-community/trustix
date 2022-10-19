// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	decider "github.com/nix-community/trustix/packages/trustix/internal/config/decider"
	signer "github.com/nix-community/trustix/packages/trustix/internal/config/signer"
)

type Config struct {
	Deciders    map[string][]*decider.Decider `toml:"decider" json:"decider"`
	Publishers  []*Publisher                  `toml:"publishers" json:"publishers"`
	Subscribers []*Subscriber                 `toml:"subscribers" json:"subscribers"`
	Signers     map[string]*signer.Signer     `toml:"signers" json:"signers"`
	Storage     *Storage                      `toml:"storage" json:"storage"`
	Remotes     []string                      `toml:"remotes" json:"remotes"`
	WriteTokens []string                      `toml:"write_tokens" json:"write_tokens"`
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

	for protocol, deciders := range c.Deciders {

		if protocol == "" {
			return fmt.Errorf("Empty protocol not allowed")
		}

		for _, decider := range deciders {
			if err := decider.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func NewConfigFromFile(path string) (*Config, error) {
	conf := &Config{}

	switch filepath.Ext(path) {

	case ".toml":
		if _, err := toml.DecodeFile(path, &conf); err != nil {
			return nil, err
		}

	case ".json":
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, &conf)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Unhandled config file extension: '%s'", filepath.Ext(path))
	}

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return conf, nil
}
