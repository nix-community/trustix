package main

import (
// "crypto/rand"
// "crypto/sha256"
// "fmt"
// "github.com/lazyledger/smt"
)

func main() {

	store, err := newGitKVStore("repo", "commiter", "commiter@example.com")
	if err != nil {
		panic(err)
	}

	store.Set([]byte("lol"), []byte("boll"))

	// tree := smt.NewSparseMerkleTree(store, sha256.New())

	// for i := 0; i < 1; i++ {
	// 	fmt.Println(i)
	// 	// a := make([]byte, 32)
	// 	// b := make([]byte, 32)

	// 	// rand.Read(a)
	// 	// rand.Read(b)

	// 	a := []byte(fmt.Sprintf("lolboll%d", i))
	// 	b := []byte(fmt.Sprintf("testhest%d", i))

	// 	tree.Update(a, b)

	// 	fmt.Println("Proofing")
	// 	proof, _ := tree.Prove(a)
	// 	root := tree.Root() // We also need the current tree root for the proof
	// 	fmt.Println("Done proofing")

	// 	// Verify the Merkle proof for foo=bar
	// 	if smt.VerifyProof(proof, root, a, b, sha256.New()) {
	// 		fmt.Println("Proof verification succeeded.")
	// 	} else {
	// 		fmt.Println("Proof verification failed.")
	// 	}

	// }

	// Update the key "foo" with the value "bar"
	// tree.Update([]byte("foo"), []byte("bar"))

	// Generate a Merkle proof for foo=bar

	// fmt.Println("Proofing")
	// proof, _ := tree.Prove([]byte("foo"))
	// root := tree.Root() // We also need the current tree root for the proof

	// // Verify the Merkle proof for foo=bar
	// if smt.VerifyProof(proof, root, []byte("foo"), []byte("bar"), sha256.New()) {
	// 	fmt.Println("Proof verification succeeded.")
	// } else {
	// 	fmt.Println("Proof verification failed.")
	// }
}
