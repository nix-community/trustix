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
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/client"
	"github.com/tweag/trustix/contrib/nix/schema"
	pb "github.com/tweag/trustix/proto"
)

var binaryCacheCommand = &cobra.Command{
	Use:   "binary-cache-proxy",
	Short: "Run a Trustix based binary cache proxy",
	RunE: func(cmd *cobra.Command, args []string) error {

		// TODO: Get from remote trustix node & check at startup
		caches := []string{"https://cache.nixos.org"}
		keyPrefix := "trustix-1"

		_, signer, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}

		conn, err := client.CreateClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewTrustixCombinedRPCClient(conn)

		encoding := base32.NewEncoding("0123456789abcdfghijklmnpqrsvwxyz")

		handler := func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			switch r.URL.Path {
			case "/robots.txt":
				w.Header().Set("Content-Type", "text/plain")
				fmt.Fprintf(w, "User-Agent: *\nDisallow: /\n")
				return
			case "/nix-cache-info":
				w.Header().Set("Content-Type", "text/plain")
				fmt.Fprintf(w, "StoreDir: /nix/store\n")
				fmt.Fprintf(w, "WantMassQuery: 1\n")
				fmt.Fprintf(w, "Priority: 40\n")
				return
			}

			if r.Method == "HEAD" || r.Method == "GET" {

				if strings.HasSuffix(r.URL.Path, ".narinfo") {

					storeHash, err := encoding.DecodeString(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), ".narinfo"))
					if err != nil {
						panic(err)
					}

					resp, err := c.Decide(r.Context(), &pb.KeyRequest{
						Key: storeHash,
					})
					if err != nil {
						log.WithFields(log.Fields{
							"path": r.URL.Path,
							"err":  err,
						}).Error("Could not reach decision for narinfo")
						w.WriteHeader(404)
						return
					}

					narinfo := &schema.NarInfo{}
					err = proto.Unmarshal(resp.Decision.Value, narinfo)
					if err != nil {
						log.WithFields(log.Fields{
							"path": r.URL.Path,
							"err":  err,
						}).Error("Could not unmarshal narinfo")
						w.WriteHeader(500)
						return
					}

					sig, err := narinfo.Sign(signer)
					if err != nil {
						panic(err)
					}

					w.Header().Set("Content-Type", "text/x-nix-narinfo")
					fmt.Fprintf(w, narinfo.ToString(fmt.Sprintf("Sig: %s:%s", keyPrefix, base64.StdEncoding.EncodeToString(sig))))

					return

				} else if strings.HasPrefix(r.URL.Path, "/nar") {
					for _, cache := range caches {
						URL, err := url.Parse(cache)
						if err != nil {
							panic(err)
						}
						URL.Path = r.URL.Path
						u := URL.String()

						resp, err := http.Get(u)
						if err != nil {
							panic(err)
						}
						if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
							continue
						}

						w.WriteHeader(200)
						_, err = io.Copy(w, resp.Body)
						if err != nil {
							panic(err)
						}

						return
					}
				} else {
					w.WriteHeader(400)
					return
				}

				w.WriteHeader(404)
			} else {
				w.WriteHeader(405)
			}
		}

		http.Handle("/", http.HandlerFunc(handler))
		log.Fatal(http.ListenAndServe(":8080", nil))

		return nil
	},
}
