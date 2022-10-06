package config

type HydraChannel struct {
	BaseURL string
	Project string
	Jobset  string
}

type Channel struct {
	Type  string        `toml:"type" json:"type"`
	Hydra *HydraChannel `toml:"hydra" json:"hydra"`
}
