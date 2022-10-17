// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

import "fmt"

type HydraJobset struct {
	BaseURL      string `toml:"base_url" json:"base_url"`
	Project      string `toml:"project" json:"project"`
	Jobset       string `toml:"jobset" json:"jobset"`
	PollInterval int64  `toml:"interval" json:"interval"`
}

func (j *HydraJobset) Validate() error {
	if j.BaseURL == "" {
		return fmt.Errorf("missing baseURL")
	}

	if j.Project == "" {
		return fmt.Errorf("missing project")
	}

	if j.Jobset == "" {
		return fmt.Errorf("missing jobset")
	}

	return nil
}

func (j *HydraJobset) init() {
	// Every hour by default
	if j.PollInterval == 0 {
		j.PollInterval = 60 * 60
	}
}

type Channels struct {
	Hydra map[string]*HydraJobset `toml:"hydra" json:"hydra"`
}

func (c *Channels) init() {
	if c.Hydra == nil {
		c.Hydra = make(map[string]*HydraJobset)
	}

	for _, hydraJobset := range c.Hydra {
		hydraJobset.init()
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
