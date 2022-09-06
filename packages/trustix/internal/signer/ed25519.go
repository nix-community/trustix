// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package signer

import (
	"crypto"
	"crypto/ed25519"
	"fmt"
)

type ed25519Verifier struct {
	pub ed25519.PublicKey
}

func (v *ed25519Verifier) Public() crypto.PublicKey {
	return v.pub
}

func (v *ed25519Verifier) Verify(message, sig []byte) bool {
	return ed25519.Verify(v.pub, message, sig)
}

func NewED25519Verifier(pub []byte) (Verifier, error) {
	if len(pub) != 32 {
		return nil, fmt.Errorf("Wrong key length: %d != 32", len(pub))
	}

	return &ed25519Verifier{
		pub: pub,
	}, nil
}

func NewED25519Signer(privKeyPath string) (crypto.Signer, error) {
	privBytes, err := readKey(privKeyPath)
	if err != nil {
		return nil, err
	}

	return ed25519.NewKeyFromSeed(privBytes[:32]), nil
}
