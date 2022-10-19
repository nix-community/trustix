// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package auth

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

type PublicToken struct {
	Name   string
	Verify func(message, sig []byte) bool
}

type PrivateToken struct {
	Name string
	Sign func(message []byte) ([]byte, error)
}

// Reads a public key and returns a public token instance that can verify API calls
func NewPublicTokenFromPub(token string) (*PublicToken, error) {
	n := strings.IndexRune(token, ':')
	if n < 1 {
		return nil, fmt.Errorf("bad format: missing name/key separator")
	}

	name := string(token[:n])

	pubBytes, err := base64.StdEncoding.DecodeString(token[n+1:])
	if err != nil {
		return nil, fmt.Errorf("couldnt decode public key: %w", err)
	}

	pub := ed25519.PublicKey(pubBytes)

	return &PublicToken{
		Name: name,
		Verify: func(message, sig []byte) bool {
			return ed25519.Verify(pub, message, sig)
		},
	}, nil
}

// Reads a private key and returns a public token instance that can verify API calls
func NewPublicTokenFromPriv(r io.Reader) (*PublicToken, error) {
	tokenBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading token private key: %w", err)
	}

	n := bytes.IndexRune(tokenBytes, ':')
	if n < 1 {
		return nil, fmt.Errorf("bad format: missing name/key separator")
	}

	name := string(tokenBytes[:n])

	privBytesB64 := tokenBytes[n+1:]

	privBytes := make([]byte, base64.StdEncoding.DecodedLen(len(privBytesB64)))
	_, err = base64.StdEncoding.Decode(privBytes, privBytesB64)
	if err != nil {
		return nil, fmt.Errorf("couldnt decode private key: %w", err)
	}

	priv := ed25519.PrivateKey(privBytes[:64])
	pub := priv.Public()

	ed25519Pub := pub.(ed25519.PublicKey)

	return &PublicToken{
		Name: name,
		Verify: func(message, sig []byte) bool {
			return ed25519.Verify(ed25519Pub, message, sig)
		},
	}, nil
}

// Reads a private key and returns a private token instance that can sign API calls
func NewPrivateToken(r io.Reader) (*PrivateToken, error) {
	tokenBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading token private key: %w", err)
	}

	n := bytes.IndexRune(tokenBytes, ':')
	if n < 1 {
		return nil, fmt.Errorf("bad format: missing name/key separator")
	}

	name := string(tokenBytes[:n])

	privBytesB64 := tokenBytes[n+1:]

	privBytes := make([]byte, base64.StdEncoding.DecodedLen(len(privBytesB64)))
	_, err = base64.StdEncoding.Decode(privBytes, privBytesB64)
	if err != nil {
		return nil, fmt.Errorf("couldnt decode private key: %w", err)
	}

	priv := ed25519.PrivateKey(privBytes[:64])

	return &PrivateToken{
		Name: name,
		Sign: func(message []byte) ([]byte, error) {
			return priv.Sign(rand.Reader, message, crypto.SignerOpts(crypto.Hash(0)))
		},
	}, nil
}
