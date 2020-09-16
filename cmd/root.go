package cmd

import (
	"crypto/sha256"
	"fmt"
	"github.com/lazyledger/smt"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/signer"
	"github.com/tweag/trustix/sth"
	"github.com/tweag/trustix/storage"
	"os"
	"sync"
)

var once sync.Once
var configPath string

var rootCmd = &cobra.Command{
	Use:   "trustix",
	Short: "Trustix",
	Long:  `Trustix`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if configPath == "" {
			return fmt.Errorf("Missing config flag")
		}

		config, err := config.NewConfigFromFile(configPath)
		if err != nil {
			Error(err)
		}

		for _, logConfig := range config.Logs {
			hasher := sha256.New()

			sig, err := signer.FromConfig(logConfig.Signer)
			if err != nil {
				Error(err)
			}

			if !sig.CanSign() {
				Error(fmt.Errorf("Cannot sign using the current configuration, aborting."))
			}

			store, err := storage.FromConfig(logConfig.Storage)
			if err != nil {
				Error(err)
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

				fmt.Println(i)

				err = store.Update(func(txn storage.Transaction) error {
					mapStore.setTxn(txn)
					defer mapStore.unsetTxn()

					a := []byte(fmt.Sprintf("lolboll%d", i))
					b := []byte(fmt.Sprintf("testhest%d", i))

					tree.Update(a, b)

					sth, err := sthManager.Sign()
					if err != nil {
						return err
					}

					return mapStore.Set([]byte("HEAD"), sth)
				})

			}
		}

		return nil
	},
}

func initCommands() {
	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to config.toml")

	rootCmd.AddCommand(generateKeyCmd)
	initGenerate()
}

func Execute() {
	once.Do(initCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
