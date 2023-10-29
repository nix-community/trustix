// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package sth

import (
	"crypto"
	"crypto/rand"
	"encoding/binary"

	"github.com/celestiaorg/smt"
	"github.com/nix-community/trustix/packages/trustix-proto/protocols"
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	vlog "github.com/nix-community/trustix/packages/trustix/internal/log"
	"github.com/nix-community/trustix/packages/trustix/internal/signer"
)

func uint64ToBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b[:]
}

func SignHead(vLog *vlog.VerifiableLog, smTree *smt.SparseMerkleTree, vMapLog *vlog.VerifiableLog, signer crypto.Signer, pd *protocols.ProtocolDescriptor) (*schema.LogHead, error) {
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

	h := pd.NewHash()
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

func VerifyLogHeadSig(verifier signer.Verifier, head *schema.LogHead, pd *protocols.ProtocolDescriptor) bool {

	h := pd.NewHash()
	h.Write(head.LogRoot)
	h.Write(uint64ToBytes(*head.TreeSize))
	h.Write(head.MapRoot)
	h.Write(head.MHRoot)
	h.Write(uint64ToBytes(*head.MHTreeSize))
	sum := h.Sum(nil)

	return verifier.Verify(sum, head.Signature)
}
