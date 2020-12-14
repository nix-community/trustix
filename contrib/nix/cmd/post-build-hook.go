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
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/api"
	"github.com/tweag/trustix/client"
	"github.com/tweag/trustix/contrib/nix/nar"
)

const NIX_STORE_DIR = "/nix/store"

func runCommand(command ...string) (string, error) {
	fmt.Println(command)

	cmd := exec.Command(command[0], command[1:]...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}

var nixHookCommand = &cobra.Command{
	Use:   "post-build-hook",
	Short: "Submit hashes for inclusion in the log (Nix post-build hook)",
	RunE: func(cmd *cobra.Command, args []string) error {

		upstreamCache := "https://cache.nixos.org"

		storePaths := strings.Split(os.Getenv("OUT_PATHS"), " ")
		if len(storePaths) < 1 {
			log.Fatal("OUT_PATHS is empty, expected at least one path to submit")
		}

		req := &api.SubmitRequest{}

		errChan := make(chan error, len(storePaths))
		wg := new(sync.WaitGroup)
		mux := new(sync.Mutex)

		encoding := base32.NewEncoding("0123456789abcdfghijklmnpqrsvwxyz")

		tmpDir, err := ioutil.TempDir("", "nix-trustix")
		if err != nil {
			return err
		}
		err = os.RemoveAll(tmpDir)
		if err != nil {
			return err
		}

		for _, storePath := range storePaths {
			storePath := storePath
			if storePath == "" {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				storeHashStr := strings.Split(filepath.Base(storePath), "-")[0]
				storeHash, err := encoding.DecodeString(storeHashStr)
				if err != nil {
					errChan <- err
					return
				}
				if len(storeHash) == 0 {
					errChan <- fmt.Errorf("Empty decoded store path hash")
					return
				}

				URL, err := url.Parse(upstreamCache)
				if err != nil {
					errChan <- err
					return
				}
				URL.Path = fmt.Sprintf("%s.narinfo", storeHashStr)

				resp, err := http.Get(URL.String())
				if err != nil {
					errChan <- err
					return
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)

				narinfo, err := nar.ParseNarInfo(string(body))
				if err != nil {
					errChan <- err
					return
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

		conn, err := client.CreateClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := client.CreateContext(30)
		defer cancel()

		c := api.NewTrustixLogAPIClient(conn)
		_, err = c.Submit(ctx, req)
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		return nil
	},
}
