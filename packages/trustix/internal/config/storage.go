// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

import (
	"fmt"
)

type NativeStorage struct {
}

type Storage struct {
	Type   string         `toml:"type" json:"type"`
	Native *NativeStorage `toml:"native" json:"native"`
}

func (s *Storage) Validate() error {
	if s == nil {
		return fmt.Errorf("No storage type configured")
	}

	switch s.Type {
	case "native":
		return nil
	case "memory":
		return nil
	default:
		return fmt.Errorf("Unhandled storage type: '%s'", s.Type)
	}
}
