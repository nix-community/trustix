package sth

import (
	"crypto"
	"crypto/rand"
	"encoding/json"
	"github.com/lazyledger/smt"
	vlog "github.com/tweag/trustix/log"
)

// STHManager - A manager that keeps a reference to a signer
// This is analogous to Certificate Transparency's concept of a Signed Tree Head
type STHManager struct {
	signer crypto.Signer
	smTree *smt.SparseMerkleTree
	vLog   *vlog.VerifiableLog
}

func NewSTHManager(smTree *smt.SparseMerkleTree, vLog *vlog.VerifiableLog, signer crypto.Signer) *STHManager {
	return &STHManager{
		signer: signer,
		smTree: smTree,
		vLog:   vLog,
	}
}

// Sign - Sign/Marshal the current state of the smTree
func (sth *STHManager) Sign() ([]byte, error) {

	opts := crypto.SignerOpts(crypto.Hash(0))

	smTreeRoot := sth.smTree.Root()
	smTreeSig, err := sth.signer.Sign(rand.Reader, smTreeRoot, opts)
	if err != nil {
		return nil, err
	}

	vLogRoot := sth.vLog.Root()
	sthSig, err := sth.signer.Sign(rand.Reader, vLogRoot, opts)
	if err != nil {
		return nil, err
	}

	s := newSMH(sth.vLog.Size(), vLogRoot, sthSig, smTreeRoot, smTreeSig)
	return json.Marshal(s)
}
