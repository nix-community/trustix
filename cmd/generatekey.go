package cmd

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var generatePrivateKeyOutput string
var generatePublicKeyOutput string

var generateKeyCmd = &cobra.Command{
	Use:   "generate-key",
	Short: "Generate an ed25519 public/private key pair",
	RunE: func(cmd *cobra.Command, args []string) error {
		if generatePublicKeyOutput == "" {
			return fmt.Errorf("Missing pubkey flag")
		}

		if generatePrivateKeyOutput == "" {
			return fmt.Errorf("Missing privkey flag")
		}

		pub, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			Error(err)
		}

		err = ioutil.WriteFile(generatePublicKeyOutput, []byte(base64.StdEncoding.EncodeToString(pub)), 0644)
		if err != nil {
			Error(err)
		}

		err = ioutil.WriteFile(generatePrivateKeyOutput, []byte(base64.StdEncoding.EncodeToString(priv)), 0644)
		if err != nil {
			Error(err)
		}

		return nil
	},
}

func initGenerate() {
	generateKeyCmd.Flags().StringVar(&generatePrivateKeyOutput, "privkey", "", "Output private key to file")
	generateKeyCmd.Flags().StringVar(&generatePublicKeyOutput, "pubkey", "", "Output public key to file")
}
