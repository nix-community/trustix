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
	"archive/tar"
	"compress/gzip"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix/api"
	"github.com/tweag/trustix/packages/trustix/client"
	"google.golang.org/protobuf/proto"
)

var exportFile string

var exportCommand = &cobra.Command{
	Use:   "export",
	Short: "Export a log to archive",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := client.CreateClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := api.NewTrustixLogAPIClient(conn)
		ctx, cancel := client.CreateContext(timeout)
		defer cancel()

		log.Debug("Requesting STH")
		sth, err := c.GetSTH(ctx, &api.STHRequest{})
		if err != nil {
			log.Fatalf("could not get STH: %v", err)
		}

		tarFile, err := os.Create(exportFile)
		if err != nil {
			log.Fatal(err)
		}
		defer tarFile.Close()

		tw := tar.NewWriter(tarFile)
		if strings.HasSuffix(exportFile, ".gz") {
			gz := gzip.NewWriter(tarFile)
			defer gz.Close()
			tw = tar.NewWriter(gz)
		}
		defer tw.Close()

		prev := uint64(0)
		id := 0
		i := uint64(50)
		for {
			if i > *sth.TreeSize {
				i = *sth.TreeSize - 1
			}

			pprev := prev
			ii := i
			req := &api.GetLogEntriesRequest{
				Start:  &pprev,
				Finish: &ii,
			}
			prev = i
			resp, err := c.GetLogEntries(ctx, req)
			if err != nil {
				log.Fatalf("Could not get entries: %v", err)
			}

			for _, leaf := range resp.Leaves {
				content, err := proto.Marshal(leaf)
				if err != nil {
					log.Fatal(err)
				}

				hdr := &tar.Header{
					Name: fmt.Sprintf("log-%d", id),
					Mode: 0600,
					Size: int64(len(content)),
				}
				if err := tw.WriteHeader(hdr); err != nil {
					log.Fatal(err)
				}
				if _, err := tw.Write([]byte(content)); err != nil {
					log.Fatal(err)
				}

				id++
			}

			if i >= *sth.TreeSize-1 {
				break
			}

			i += 50
		}

		prev = uint64(0)
		id = 0
		i = uint64(50)
		for {
			if i > *sth.MHTreeSize {
				i = *sth.MHTreeSize - 1
			}

			pprev := prev
			ii := i
			req := &api.GetLogEntriesRequest{
				Start:  &pprev,
				Finish: &ii,
			}
			prev = i
			resp, err := c.GetMHLogEntries(ctx, req)
			if err != nil {
				log.Fatalf("Could not get entries: %v", err)
			}

			for _, leaf := range resp.Leaves {
				content, err := proto.Marshal(leaf)
				if err != nil {
					log.Fatal(err)
				}

				hdr := &tar.Header{
					Name: fmt.Sprintf("maplog-%d", id),
					Mode: 0600,
					Size: int64(len(content)),
				}
				if err := tw.WriteHeader(hdr); err != nil {
					log.Fatal(err)
				}
				if _, err := tw.Write([]byte(content)); err != nil {
					log.Fatal(err)
				}

				id++
			}

			if i >= *sth.MHTreeSize-1 {
				break
			}

			i += 50
		}

		content, err := proto.Marshal(sth)
		if err != nil {
			log.Fatalf("Could not marsal STH: %v", err)
		}

		hdr := &tar.Header{
			Name: "STH",
			Mode: 0600,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			return err
		}

		return nil
	},
}

func initExport() {
	exportCommand.Flags().StringVar(&exportFile, "output", "trustix-dump.tar.gz", "File to dump log to")
}
