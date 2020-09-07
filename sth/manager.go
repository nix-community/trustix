package sth

import (
	"crypto"
	"crypto/rand"
	"encoding/json"
	"github.com/lazyledger/smt"
)

// STHManager - A manager that keeps a reference to a signer
// This is analogous to Certificate Transparency's concept of a Signed Tree Head
type STHManager struct {
	signer crypto.Signer
	tree   *smt.SparseMerkleTree
}

func NewSTHManager(tree *smt.SparseMerkleTree, signer crypto.Signer) *STHManager {
	return &STHManager{
		signer: signer,
		tree:   tree,
	}
}

// Sign - Sign/Marshal the current state of the tree
func (sth *STHManager) Sign() ([]byte, error) {

	opts := crypto.SignerOpts(crypto.Hash(0))
	root := sth.tree.Root()

	sig, err := sth.signer.Sign(rand.Reader, root, opts)
	if err != nil {
		return nil, err
	}

	s := newSTH(root, sig)
	return json.Marshal(s)
}
