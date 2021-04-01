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
