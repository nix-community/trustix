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
	"bytes"
	"encoding/base32"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/api"
)

const NIX_STORE_DIR = "/nix/store"

var nixHookCommand = &cobra.Command{
	Use:   "post-build-hook",
	Short: "Submit hashes for inclusion in the log (Nix post-build hook)",
	RunE: func(cmd *cobra.Command, args []string) error {

		storePaths := strings.Split(os.Getenv("OUT_PATHS"), " ")
		if len(storePaths) < 1 {
			log.Fatal("OUT_PATHS is empty, expected at least one path to submit")
		}

		req := &api.SubmitRequest{}

		errChan := make(chan error, len(storePaths))
		wg := new(sync.WaitGroup)
		mux := new(sync.Mutex)

		encoding := base32.NewEncoding("0123456789abcdfghijklmnpqrsvwxyz")

		for _, storePath := range storePaths {
			storePath := storePath
			if storePath == "" {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				storeHashStr := strings.Split(strings.TrimPrefix(storePath, NIX_STORE_DIR+"/"), "-")[0]
				storeHash, err := encoding.DecodeString(storeHashStr)
				if err != nil {
					errChan <- err
					return
				}

				if len(storeHash) == 0 {
					errChan <- fmt.Errorf("Empty decoded store path hash")
					return
				}

				cmd := exec.Command("nix-store", "--query", "--hash", storePath)
				var stdout bytes.Buffer
				cmd.Stdout = &stdout

				err = cmd.Run()
				if err != nil {
					errChan <- err
					return
				}

				stdoutBytes := bytes.TrimSpace(stdout.Bytes())
				if len(stdoutBytes) == 0 {
					errChan <- fmt.Errorf("Empty decoded store hash")
					return
				}

				components := bytes.Split(stdoutBytes, []byte(":"))
				if len(components) != 2 {
					errChan <- fmt.Errorf("Malformed store hash output")
					return
				}

				if !bytes.Equal(components[0], []byte("sha256")) {
					errChan <- fmt.Errorf("Store hash type mismatch")
					return
				}

				// Pad
				b32Hash := string(components[1])
				for i := len(b32Hash); i <= 56; i++ {
					b32Hash = b32Hash + "="
				}

				hash, err := encoding.DecodeString(b32Hash)
				if err != nil {
					errChan <- err
					return
				}

				log.WithFields(log.Fields{
					"inputHash":  storeHashStr,
					"outputHash": string(components[1]),
				}).Debug("Submitting mapping")

				mux.Lock()
				req.Items = append(req.Items, &api.KeyValuePair{
					Key:   storeHash,
					Value: hash,
				})
				mux.Unlock()
			}()
		}

		wg.Wait()
		close(errChan)

		for err := range errChan {
			log.Fatalf("Could not hash store path: %v", err)
		}

		conn, err := createClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := createContext()
		defer cancel()

		c := api.NewTrustixLogAPIClient(conn)
		r, err := c.Submit(ctx, req)
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		fmt.Println(r.Status)

		return nil
	},
}
