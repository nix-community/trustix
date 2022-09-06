// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

import (
	"fmt"

	signer "github.com/nix-community/trustix/packages/trustix/internal/config/signer"
)

type Publisher struct {
	Protocol  string            `toml:"protocol" json:"protocol"`
	Signer    string            `toml:"signer" json:"signer"`
	PublicKey *PublicKey        `toml:"publicKey" json:"publicKey"`
	Meta      map[string]string `toml:"meta" json:"meta"`
}

func (p *Publisher) Validate(signers map[string]*signer.Signer) error {
	if p.Signer == "" {
		return missingField("signer")
	}

	if p.Protocol == "" {
		return missingField("protocol")
	}

	if err := p.PublicKey.Validate(); err != nil {
		return err
	}

	_, ok := signers[p.Signer]
	if !ok {
		return fmt.Errorf("Signer '%s' referenced but does not exist", p.Signer)
	}

	return nil
}

func (p *Publisher) GetMeta() map[string]string {
	if p.Meta != nil {
		return p.Meta
	}
	return make(map[string]string)
}
