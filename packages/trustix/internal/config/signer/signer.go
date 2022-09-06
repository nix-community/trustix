// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

import (
	"fmt"
)

type ED25519 struct {
	PrivateKeyPath string `toml:"private-key-path" json:"private-key-path"`
}

func (s *ED25519) Validate() error {
	if s.PrivateKeyPath == "" {
		return fmt.Errorf("Empty private key path")
	}
	return nil
}

type Signer struct {
	Type    string   `toml:"type" json:"type"`
	ED25519 *ED25519 `toml:"ed25519" json:"ed25519"`
}

func (s *Signer) Validate() error {
	switch s.Type {
	case "ed25519":
		return s.ED25519.Validate()
	default:
		return fmt.Errorf("Unhandled signer type: %s", s.Type)
	}
}
