// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

import "fmt"

type HydraJobset struct {
	BaseURL string `toml:"base_url" json:"base_url"`
	Project string `toml:"project" json:"project"`
	Jobset  string `toml:"jobset" json:"jobset"`
}

func (c *HydraJobset) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("missing baseURL")
	}

	if c.Project == "" {
		return fmt.Errorf("missing project")
	}

	if c.Jobset == "" {
		return fmt.Errorf("missing jobset")
	}

	return nil
}

type Channels struct {
	Hydra map[string]*HydraJobset `toml:"hydra" json:"hydra"`
}

func (c *Channels) init() {
	if c.Hydra == nil {
		c.Hydra = make(map[string]*HydraJobset)
	}
}

func (c *Channels) Validate() error {
	for name, jobset := range c.Hydra {
		err := jobset.Validate()
		if err != nil {
			return fmt.Errorf("error validating jobset '%s': %w", name, err)
		}
	}

	return nil
}
