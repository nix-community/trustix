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
	"encoding/binary"

	"github.com/lazyledger/smt"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	vlog "github.com/tweag/trustix/packages/trustix/log"
	"github.com/tweag/trustix/packages/trustix/signer"
)

func uint64ToBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b[:]
}

func SignHead(vLog *vlog.VerifiableLog, smTree *smt.SparseMerkleTree, vMapLog *vlog.VerifiableLog, signer crypto.Signer) (*schema.STH, error) {
	opts := crypto.SignerOpts(crypto.Hash(0))

	vLogRoot, err := vLog.Root()
	if err != nil {
		return nil, err
	}
	smTreeRoot := smTree.Root()

	_, err = vMapLog.Append(smTreeRoot)
	if err != nil {
		return nil, err
	}

	vMapLogRoot, err := vMapLog.Root()
	if err != nil {
		return nil, err
	}

	vLogSize := vLog.Size()
	vMapLogSize := vMapLog.Size()

	h := sha256.New()
	h.Write(vLogRoot)
	h.Write(uint64ToBytes(vLogSize))
	h.Write(smTreeRoot)
	h.Write(vMapLogRoot)
	h.Write(uint64ToBytes(vMapLogSize))
	sum := h.Sum(nil)

	sig, err := signer.Sign(rand.Reader, sum, opts)
	if err != nil {
		return nil, err
	}

	return &schema.STH{
		LogRoot:    vLogRoot,
		TreeSize:   &vLogSize,
		MapRoot:    smTreeRoot,
		MHRoot:     vMapLogRoot,
		MHTreeSize: &vMapLogSize,
		Signature:  sig,
	}, nil
}

func VerifySTHSig(verifier signer.TrustixVerifier, sth *schema.STH) bool {

	h := sha256.New()
	h.Write(sth.LogRoot)
	h.Write(uint64ToBytes(*sth.TreeSize))
	h.Write(sth.MapRoot)
	h.Write(sth.MHRoot)
	h.Write(uint64ToBytes(*sth.MHTreeSize))
	sum := h.Sum(nil)

	return verifier.Verify(sum, sth.Signature)
}
