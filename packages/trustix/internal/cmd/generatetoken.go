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

var generateTokenName string
var generateTokenOutputPriv string
var generateTokenOutputPub string

const pubTokenHelp = `  Put this public token into your trustix configuration as such:
	Nix:
	{
	  writeTokens = [ "%s" ]
	}

	TOML:
	write_tokens = [ "%s" ]

`

const privTokenHelp = `  Write the private token to a file and make Trustix aware of it using $TRUSTIX_TOKEN:
    $ echo '%s' > /var/run/trustix.token
    $ export TRUSTIX_TOKEN=/var/run/trustix.token
`

var generateTokenCmd = &cobra.Command{
	Use:   "generate-token",
	Short: "Generate a write token (needed for publishing log entries)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if generateTokenName == "" {
			log.Fatal("Missing required param: name")
		}

		fmtToken := func(t string) string {
			return generateTokenName + ":" + t
		}

		pub, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			log.Fatal(err)
		}

		pubEncoded := fmtToken(base64.StdEncoding.EncodeToString(pub))
		privEncoded := fmtToken(base64.StdEncoding.EncodeToString(priv))

		if generateTokenOutputPub != "" {
			err = os.WriteFile(generateTokenOutputPub, []byte(pubEncoded), 0644)
			if err != nil {
				log.Fatal(err)
			}

		} else {
			fmt.Println("Public token:", pubEncoded)
			fmt.Printf(pubTokenHelp, pubEncoded, pubEncoded)
		}

		if generateTokenOutputPriv != "" {
			err = os.WriteFile(generateTokenOutputPriv, []byte(privEncoded), 0600)
			if err != nil {
				log.Fatal(err)
			}

		} else {
			fmt.Println("Private token:", privEncoded)
			fmt.Printf(privTokenHelp, privEncoded)
		}

		return nil
	},
}

func initGenerateToken() {
	generateTokenCmd.Flags().StringVar(&generateTokenName, "name", "", "Name of the token, used as an identifier")
	generateTokenCmd.Flags().StringVar(&generateTokenOutputPriv, "privkey", "", "Output token private key to file")
	generateTokenCmd.Flags().StringVar(&generateTokenOutputPub, "pubkey", "", "Output token public key to file")
}
