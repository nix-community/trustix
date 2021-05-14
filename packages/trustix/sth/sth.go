// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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

func SignHead(vLog *vlog.VerifiableLog, smTree *smt.SparseMerkleTree, vMapLog *vlog.VerifiableLog, signer crypto.Signer) (*schema.LogHead, error) {
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

	return &schema.LogHead{
		LogRoot:    vLogRoot,
		TreeSize:   &vLogSize,
		MapRoot:    smTreeRoot,
		MHRoot:     vMapLogRoot,
		MHTreeSize: &vMapLogSize,
		Signature:  sig,
	}, nil
}

func VerifyLogHeadSig(verifier signer.TrustixVerifier, head *schema.LogHead) bool {

	h := sha256.New()
	h.Write(head.LogRoot)
	h.Write(uint64ToBytes(*head.TreeSize))
	h.Write(head.MapRoot)
	h.Write(head.MHRoot)
	h.Write(uint64ToBytes(*head.MHTreeSize))
	sum := h.Sum(nil)

	return verifier.Verify(sum, head.Signature)
}
