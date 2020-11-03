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
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/config"
)

func FromConfig(signerConfig *config.SignerConfig) (crypto.Signer, error) {

	if signerConfig.Type == "" {
		fmt.Errorf("Missing signer config field 'type'.")
	}

	log.WithField("type", signerConfig.Type).Info("Creating signer")
	switch signerConfig.Type {
	case "ed25519":
		return newED25519Signer(signerConfig)
	}

	return nil, fmt.Errorf("Signer type '%s' is not supported.", signerConfig.Type)

}
