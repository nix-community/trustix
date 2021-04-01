// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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

func NewED25519Verifier(pub []byte) (TrustixVerifier, error) {
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
