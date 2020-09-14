package main // import "github.com/tweag/trustix"

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"github.com/lazyledger/smt"
	"github.com/tweag/trustix/sth"
	"github.com/tweag/trustix/store"
)

func main() {

	config, err := newConfig("./config.toml")
	if err != nil {
		panic(err)
	}

	for _, logConfig := range config.Logs {

		_, priv, _ := ed25519.GenerateKey(nil)

		if logConfig.Storage.Type != "git" {
			panic("Only git implemented at this time")
		}

		kvStore, err := store.NewGitKVStore(logConfig.Storage.Git.Path, logConfig.Storage.Git.Commiter, logConfig.Storage.Git.Email)
		if err != nil {
			panic(err)
		}

		hasher := sha256.New()

		var tree *smt.SparseMerkleTree
		oldHead, err := kvStore.GetRaw([]string{"HEAD"})
		if err != nil {
			// No STH yet, new tree
			if err == store.ObjectNotFoundError {
				tree = smt.NewSparseMerkleTree(kvStore, hasher)
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

			tree = smt.ImportSparseMerkleTree(kvStore, hasher, rootBytes)
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

			kvStore.SetRaw([]string{"HEAD"}, sth)

			err = kvStore.CreateCommit(fmt.Sprintf("Set key"))
			if err != nil {
				panic(err)
			}
		}

	}

}
