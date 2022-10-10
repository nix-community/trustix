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
	Channels map[string]*Channel `toml:"channels" json:"channels"`
	Cron     *Cron               `toml:"cron" json:"cron"`
}

func (c *Config) init() {
	if c.Cron == nil {
		c.Cron = &Cron{}
	}
	c.Cron.init()
}

func (c *Config) Validate() error {
	for channel, channelConf := range c.Channels {
		err := channelConf.Validate()
		if err != nil {
			return fmt.Errorf("error validating channel config for '%s': %w", channel, err)
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

	conf.init()

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return conf, nil
}
