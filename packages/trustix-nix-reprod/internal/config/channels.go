package config

import "fmt"

type HydraChannel struct {
	BaseURL string `toml:"base_url" json:"base_url"`
	Project string `toml:"project" json:"project"`
	Jobset  string `toml:"jobset" json:"jobset"`
}

func (c *HydraChannel) Validate() error {
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

type Channel struct {
	Type  string
	Expr  string
	Hydra *HydraChannel `toml:"hydra" json:"hydra"`
}

func (c *Channel) Validate() error {
	switch c.Type {
	case "hydra":
		err := c.Hydra.Validate()
		if err != nil {
			return fmt.Errorf("error validating hydra channel config: %w", err)
		}

	default:
		return fmt.Errorf("error validating channel config: invalid type '%s'", c.Type)
	}

	if c.Expr == "" {
		return fmt.Errorf("missing expr")
	}

	return nil
}
