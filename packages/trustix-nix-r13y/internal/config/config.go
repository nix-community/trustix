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
)

type Config struct {
	Channels *Channels           `toml:"channels" json:"channels"`
	Cron     *Cron               `toml:"cron" json:"cron"`
	Lognames map[string]string   `toml:"lognames" json:"lognames"`
	Attrs    map[string][]string `toml:"attrs" json:"attrs"`
}

func (c *Config) init() {
	if c.Cron == nil {
		c.Cron = &Cron{}
	}
	c.Cron.init()

	if c.Channels == nil {
		c.Channels = &Channels{}
	}
	c.Channels.init()

	if c.Lognames == nil {
		c.Lognames = make(map[string]string)
	}

	if c.Attrs == nil {
		c.Attrs = make(map[string][]string)
	}
}

func (c *Config) Validate() error {
	err := c.Channels.Validate()
	if err != nil {
		return fmt.Errorf("error validating channels: %w", err)
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

	conf.init()

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return conf, nil
}
