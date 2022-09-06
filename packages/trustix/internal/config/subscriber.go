// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

import (
	"fmt"
)

type Subscriber struct {
	Protocol  string            `toml:"protocol" json:"protocol"`
	PublicKey *PublicKey        `toml:"publicKey" json:"publicKey"`
	SyncMode  string            `toml:"syncmode" json:"syncmode"`
	Meta      map[string]string `toml:"meta" json:"meta"`
}

func (s *Subscriber) Validate() error {
	if s.Protocol == "" {
		return missingField("protocol")
	}
	if s.SyncMode != "" && s.SyncMode != "light" {
		return fmt.Errorf("Unknown sync mode: '%s'", s.SyncMode)
	}

	return nil
}

func (s *Subscriber) GetMeta() map[string]string {
	if s.Meta != nil {
		return s.Meta
	}
	return make(map[string]string)
}
