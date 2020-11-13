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
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"github.com/lazyledger/smt"
	vlog "github.com/tweag/trustix/log"
	"github.com/tweag/trustix/schema"
	"github.com/tweag/trustix/signer"
)

func SignHead(smTree *smt.SparseMerkleTree, vLog *vlog.VerifiableLog, signer crypto.Signer) (*schema.STH, error) {
	opts := crypto.SignerOpts(crypto.Hash(0))

	vLogRoot, err := vLog.Root()
	if err != nil {
		return nil, err
	}
	smTreeRoot := smTree.Root()

	h := sha256.New()
	h.Write(vLogRoot)
	h.Write(smTreeRoot)
	sum := h.Sum(nil)

	sig, err := signer.Sign(rand.Reader, sum, opts)
	if err != nil {
		return nil, err
	}

	vLogSize := uint64(vLog.Size())
	return &schema.STH{
		LogRoot:   vLogRoot,
		TreeSize:  &vLogSize,
		MapRoot:   smTreeRoot,
		Signature: sig,
	}, nil
}

func VerifySTHSig(verifier signer.TrustixVerifier, sth *schema.STH) bool {
	h := sha256.New()
	h.Write(sth.LogRoot)
	h.Write(sth.MapRoot)
	sum := h.Sum(nil)

	return verifier.Verify(sum, sth.Signature)
}
