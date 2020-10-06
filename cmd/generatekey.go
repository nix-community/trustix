// MIT License
//
// Copyright (c) 2020 Tweag IO
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package cmd

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
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
	generateKeyCmd.Flags().StringVar(&generatePrivateKeyOutput, "privkey", "", "Output private key to file")
	generateKeyCmd.Flags().StringVar(&generatePublicKeyOutput, "pubkey", "", "Output public key to file")
}
