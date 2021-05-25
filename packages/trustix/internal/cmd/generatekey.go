// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

var generatePrivateKeyOutput string
var generatePublicKeyOutput string
var generateKeyType string

var generateKeyCmd = &cobra.Command{
	Use:   "generate-key",
	Short: "Generate a public/private key pair",
	RunE: func(cmd *cobra.Command, args []string) error {
		if generatePublicKeyOutput == "" {
			return fmt.Errorf("Missing pubkey flag")
		}

		if generatePrivateKeyOutput == "" {
			return fmt.Errorf("Missing privkey flag")
		}

		switch generateKeyType {
		case "ed25519":
		default:
			log.Fatalf("Unhandled key type: %s", generateKeyType)
		}

		pub, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(generatePublicKeyOutput, []byte(base64.StdEncoding.EncodeToString(pub)), 0644)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(generatePrivateKeyOutput, []byte(base64.StdEncoding.EncodeToString(priv)), 0644)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func initGenerate() {
	generateKeyCmd.Flags().StringVar(&generateKeyType, "type", "ed25519", "Key type to generate")
	generateKeyCmd.Flags().StringVar(&generatePrivateKeyOutput, "privkey", "", "Output private key to file")
	generateKeyCmd.Flags().StringVar(&generatePublicKeyOutput, "pubkey", "", "Output public key to file")
}
