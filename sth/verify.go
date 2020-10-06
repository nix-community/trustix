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

package sth

import (
	"fmt"
	"github.com/tweag/trustix/signer"
)

func verifySMHSig(signer signer.TrustixSigner, smh *SMH) error {

	rootBytes, err := smh.UnmarshalSMHRoot()
	if err != nil {
		return err
	}

	sigBytes, err := smh.UnmarshalSMHSignature()
	if err != nil {
		return err
	}

	if !signer.Verify(rootBytes, sigBytes) {
		return fmt.Errorf("SMH signature verification failed")
	}

	return nil
}

func verifySTHSig(signer signer.TrustixSigner, smh *SMH) error {

	rootBytes, err := smh.UnmarshalSTHRoot()
	if err != nil {
		return err
	}

	sigBytes, err := smh.UnmarshalSTHSignature()
	if err != nil {
		return err
	}

	if !signer.Verify(rootBytes, sigBytes) {
		return fmt.Errorf("SMH signature verification failed")
	}

	return nil
}
