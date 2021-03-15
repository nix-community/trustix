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
	"io/ioutil"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/api"
	"github.com/tweag/trustix/client"
)

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

				var err error

				item, err := createKVPair(storePath)
				if err != nil {
					errChan <- err
					return
				}

				mux.Lock()
				req.Items = append(req.Items, item)
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
