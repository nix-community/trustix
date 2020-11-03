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
	"crypto/ed25519"
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

func NewED25519Verifier(pub string) (TrustixVerifier, error) {
	pubkey, err := decodeKey(pub)
	if err != nil {
		return nil, err
	}

	return &ed25519Verifier{
		pub: pubkey,
	}, nil
}

func NewED25519Signer(privKeyPath string) (crypto.Signer, error) {
	privBytes, err := readKey(privKeyPath)
	if err != nil {
		return nil, err
	}

	return ed25519.NewKeyFromSeed(privBytes[:32]), nil
}
