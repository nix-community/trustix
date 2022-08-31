// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pb "github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix/client"
)

var nixHookCommand = &cobra.Command{
	Use:   "post-build-hook",
	Short: "Submit hashes for inclusion in the log (Nix post-build hook)",
	RunE: func(cmd *cobra.Command, args []string) error {

		storePaths := strings.Split(os.Getenv("OUT_PATHS"), " ")
		if len(storePaths) < 1 {
			log.Fatal("OUT_PATHS is empty, expected at least one path to submit")
		}

		if logID == "" {
			log.Fatal("Missing log id parameter")
		}

		req := &pb.SubmitRequest{
			LogID: &logID,
		}

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

		c, err := client.CreateClientConn(dialAddress)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer c.Close()

		ctx, cancel := client.CreateContext(30)
		defer cancel()

		_, err = c.LogRPC.Submit(ctx, req)
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		return nil
	},
}
