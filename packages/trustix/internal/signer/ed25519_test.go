// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package signer

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fixturePath(name string) string {
	wd, _ := os.Getwd()
	return path.Join(wd, "fixtures", name)
}

func mkSnakeOil(t *testing.T) crypto.Signer {
	signer, err := NewED25519Signer(fixturePath("priv"))
	assert.Nil(t, err)
	return signer
}

func TestSign(t *testing.T) {
	signer := mkSnakeOil(t)

	message := []byte("somepayload")
	opts := crypto.SignerOpts(crypto.Hash(0))

	sig, err := signer.Sign(rand.Reader, message, opts)
	assert.Nil(t, err)

	verifier := &ed25519Verifier{
		pub: signer.Public().(ed25519.PublicKey),
	}
	valid := verifier.Verify(message, sig)
	assert.Equal(t, valid, true, "Signature valid")
}
