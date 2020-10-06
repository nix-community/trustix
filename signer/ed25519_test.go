// MIT License
//
// Copyright (c) 2020 Tweag IO
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package signer

import (
	"crypto"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"github.com/tweag/trustix/config"
	"os"
	"path"
	"testing"
)

func fixturePath(name string) string {
	wd, _ := os.Getwd()
	return path.Join(wd, "fixtures", name)
}

func mkSnakeOil(t *testing.T) TrustixSigner {
	config := &config.SignerConfig{
		Type:    "ed25519",
		KeyType: "ed25519",
		ED25519: &config.ED25519SignerConfig{
			PrivateKeyPath: fixturePath("priv"),
		},
	}
	signer, err := genED25519Signer(config)
	assert.Nil(t, err)
	return signer
}

func TestCanSign(t *testing.T) {
	signer := mkSnakeOil(t)
	assert.Equal(t, signer.CanSign(), true, "Signer can sign")
}

func TestSign(t *testing.T) {
	signer := mkSnakeOil(t)

	message := []byte("somepayload")
	opts := crypto.SignerOpts(crypto.Hash(0))

	sig, err := signer.Sign(rand.Reader, message, opts)
	assert.Nil(t, err)

	valid := signer.Verify(message, sig)
	assert.Equal(t, valid, true, "Signature valid")
}
