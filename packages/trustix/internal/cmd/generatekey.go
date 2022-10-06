// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"
	"os"

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

		err = os.WriteFile(generatePublicKeyOutput, []byte(base64.StdEncoding.EncodeToString(pub)), 0644)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(generatePrivateKeyOutput, []byte(base64.StdEncoding.EncodeToString(priv)), 0644)
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
