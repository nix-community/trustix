package core

import (
	"crypto/sha256"
	"fmt"
	"github.com/lazyledger/smt"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/signer"
	"github.com/tweag/trustix/sth"
	"github.com/tweag/trustix/storage"
)

type TrustixCore struct{}

func CoreFromConfig(conf *config.LogConfig) (*TrustixCore, error) {

	hasher := sha256.New()

	sig, err := signer.FromConfig(conf.Signer)
	if err != nil {
		return nil, err
	}

	if !sig.CanSign() {
		return nil, fmt.Errorf("Cannot sign using the current configuration, aborting.")
	}

	store, err := storage.FromConfig(conf.Storage)
	if err != nil {
		return nil, err
	}

	var tree *smt.SparseMerkleTree

	mapStore := newMapStore()

	err = store.View(func(txn storage.Transaction) error {
		mapStore.setTxn(txn)
		defer mapStore.unsetTxn()

		oldHead, err := txn.Get([]byte("HEAD"))
		if err != nil {
			// No STH yet, new tree
			if err == storage.ObjectNotFoundError {
				tree = smt.NewSparseMerkleTree(mapStore, hasher)
			} else {
				return err
			}
		} else {
			oldSTH := &sth.STH{}
			err = oldSTH.FromJSON(oldHead)
			if err != nil {
				return err
			}

			rootBytes, err := oldSTH.UnmarshalRoot()
			if err != nil {
				return err
			}

			tree = smt.ImportSparseMerkleTree(mapStore, hasher, rootBytes)
		}

		return nil
	})

	sthManager := sth.NewSTHManager(tree, sig)

	for i := 0; i < (10); i++ {

		err = store.Update(func(txn storage.Transaction) error {
			fmt.Println(i)
			mapStore.setTxn(txn)
			defer mapStore.unsetTxn()

			a := []byte(fmt.Sprintf("lolboll%d", i))
			b := []byte(fmt.Sprintf("testhest%d", i))

			tree.Update(a, b)

			sth, err := sthManager.Sign()
			if err != nil {
				return err
			}

			// Generate a Merkle proof for foo=bar
			proof, _ := tree.Prove(a)
			root := tree.Root() // We also need the current tree root for the proof

			// Verify the Merkle proof for foo=bar
			if !smt.VerifyProof(proof, root, a, b, hasher) {
				return fmt.Errorf("Proof verification failed.")
			}

			return mapStore.Set([]byte("HEAD"), sth)
		})

		if err != nil {
			return nil, err
		}

	}
	return nil, nil
}
