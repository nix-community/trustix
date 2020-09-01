package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/lazyledger/smt"
)

func main() {
	// Initialise a new key-value store to store the nodes of the tree
	store := NewSimpleMap()
	// Initialise the tree
	tree := smt.NewSparseMerkleTree(store, sha256.New())

	// Update the key "foo" with the value "bar"
	tree.Update([]byte("foo"), []byte("bar"))

	// Generate a Merkle proof for foo=bar
	proof, _ := tree.Prove([]byte("foo"))
	root := tree.Root() // We also need the current tree root for the proof

	// Verify the Merkle proof for foo=bar
	if smt.VerifyProof(proof, root, []byte("foo"), []byte("bar"), sha256.New()) {
		fmt.Println("Proof verification succeeded.")
	} else {
		fmt.Println("Proof verification failed.")
	}
}
