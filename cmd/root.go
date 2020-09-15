package cmd

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"github.com/lazyledger/smt"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/sth"
	"github.com/tweag/trustix/store"
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
