// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"os"
	"strings"
	"sync"

	pb "github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

		tmpDir, err := os.MkdirTemp("", "nix-trustix")
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

		interceptors, err := getAuthInterceptors()
		if err != nil {
			log.Fatalf("Could not get auth interceptor: %v", err)
		}

		c, err := client.CreateClient(dialAddress, interceptors)
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
