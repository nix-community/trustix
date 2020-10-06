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
	"crypto/ed25519"
	"github.com/tweag/trustix/config"
)

func newED25519Verifier(pubkey ed25519.PublicKey) func(message, sig []byte) bool {
	return func(message, sig []byte) bool {
		return ed25519.Verify(pubkey, message, sig)
	}
}

func genED25519Signer(signerConfig *config.SignerConfig) (TrustixSigner, error) {
	pubBytes, err := decodeKey(signerConfig.PublicKey)
	if err != nil {
		return nil, err
	}
	pub := ed25519.PublicKey(pubBytes)

	// Public key is derived from private key for ed25519
	if signerConfig.ED25519 != nil && signerConfig.ED25519.PrivateKeyPath != "" {
		privBytes, err := readKey(signerConfig.ED25519.PrivateKeyPath)
		if err != nil {
			return nil, err
		}

		key := ed25519.NewKeyFromSeed(privBytes[:32])

		return &privkeySigner{
			signer:   key,
			verifier: newED25519Verifier(pub),
		}, nil

	} else {
		return &pubkeySigner{
			pub:      pub,
			verifier: newED25519Verifier(pub),
		}, nil
	}
}
