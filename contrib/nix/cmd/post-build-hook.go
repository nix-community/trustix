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
	"strconv"
	"strings"
	"sync"

	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/api"
	"github.com/tweag/trustix/contrib/nix/schema"
)

const NIX_STORE_DIR = "/nix/store"

func runCommand(command ...string) (string, error) {
	cmd := exec.Command(command[0], command[1:]...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	stdoutBytes := bytes.TrimSpace(stdout.Bytes())
	if len(stdoutBytes) == 0 {
		return "", fmt.Errorf("Empty decoded store hash")
	}
	return strings.TrimSpace(string(stdoutBytes)), nil
}

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

				narHash, err := runCommand("nix-store", "--query", "--hash", storePath)
				if err != nil {
					errChan <- err
					return
				}

				var narSize uint64
				narSizeStr, err := runCommand("nix-store", "--query", "--size", storePath)
				if err != nil {
					errChan <- err
					return
				}
				if s, err := strconv.ParseUint(narSizeStr, 10, 64); err == nil {
					narSize = s
				}

				refs := []string{}
				refsStr, err := runCommand("nix-store", "--query", "--references", storePath)
				if err != nil {
					errChan <- err
					return
				}
				for _, path := range strings.Split(refsStr, "\n") {
					if len(path) == 0 {
						continue
					}
					refs = append(refs, strings.TrimPrefix(path, NIX_STORE_DIR+"/"))
				}

				narinfo := &schema.NarInfo{
					StorePath:  &storePath,
					NarHash:    &narHash,
					NarSize:    &narSize,
					References: refs,
				}

				log.WithFields(log.Fields{
					"storePath": storePath,
				}).Debug("Submitting mapping")

				narinfoBytes, err := proto.Marshal(narinfo)
				if err != nil {
					errChan <- err
					return
				}

				mux.Lock()
				req.Items = append(req.Items, &api.KeyValuePair{
					Key:   storeHash,
					Value: narinfoBytes,
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
		_, err = c.Submit(ctx, req)
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		return nil
	},
}
