package main // import "github.com/tweag/trustix"

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"github.com/lazyledger/smt"
	"github.com/tweag/trustix/sth"
)

func main() {
	repoPath := "repo"

	_, priv, _ := ed25519.GenerateKey(nil)

	store, err := newGitKVStore(repoPath, "commiter", "commiter@example.com")
	if err != nil {
		panic(err)
	}

	hasher := sha256.New()

	var tree *smt.SparseMerkleTree
	oldHead, err := store.GetRaw([]string{"HEAD"})
	if err != nil {
		// No STH yet, new tree
		if err == ObjectNotFoundError {
			tree = smt.NewSparseMerkleTree(store, hasher)
		} else {
			panic(err)
		}
	} else {
		oldSTH := &sth.STH{}
		err = oldSTH.FromJSON(oldHead)
		if err != nil {
			panic(err)
		}

		rootBytes, err := oldSTH.UnmarshalRoot()
		if err != nil {
			panic(err)
		}

		tree = smt.ImportSparseMerkleTree(store, hasher, rootBytes)
	}

	sthManager := sth.NewSTHManager(tree, priv)

	for i := 0; i < (10); i++ {
		fmt.Println(i)

		a := []byte(fmt.Sprintf("lolboll%d", i))
		b := []byte(fmt.Sprintf("testhest%d", i))

		tree.Update(a, b)

		sth, err := sthManager.Sign()
		if err != nil {
			panic(err)
		}

		store.SetRaw([]string{"HEAD"}, sth)

		err = store.createCommit(fmt.Sprintf("Set key"))
		if err != nil {
			panic(err)
		}
	}

}
