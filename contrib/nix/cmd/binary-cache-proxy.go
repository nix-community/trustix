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
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/bakins/logrus-middleware"
	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/client"
	"github.com/tweag/trustix/contrib/nix/nar"
	"github.com/tweag/trustix/contrib/nix/schema"
	pb "github.com/tweag/trustix/proto"
)

func getCaches(c pb.TrustixCombinedRPCClient) ([]string, error) {
	ctx, cancel := client.CreateContext(30)
	defer cancel()
	resp, err := c.Logs(ctx, &pb.LogsRequest{})
	if err != nil {
		return nil, err
	}

	seen := make(map[string]string)
	caches := []string{}

	for _, log := range resp.Logs {
		upstream, ok := log.Meta["upstream"]
		if ok {
			_, hasSeen := seen[upstream]
			if !hasSeen {
				caches = append(caches, upstream)
			}
		}
	}

	return caches, nil
}

func readKey(path string) (string, crypto.Signer, error) {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	components := bytes.Split(keyBytes, []byte(":"))
	if len(components) != 2 {
		return "", nil, fmt.Errorf("Unexpected number of key components: %d", len(components))
	}

	buf := make([]byte, 64)
	n, err := base64.StdEncoding.Decode(buf, components[1])
	if err != nil {
		return "", nil, err
	}
	if n != 64 {
		return "", nil, fmt.Errorf("Expected 64 bytes, wrote %d", n)
	}

	priv := ed25519.NewKeyFromSeed(buf[:32])

	return string(components[0]), priv, nil
}

var binaryCacheCommand = &cobra.Command{
	Use:   "binary-cache-proxy",
	Short: "Run a Trustix based binary cache proxy",
	RunE: func(cmd *cobra.Command, args []string) error {

		keyPrefix, signer, err := readKey("cache-priv-key.pem")
		if err != nil {
			panic(err)
		}

		conn, err := client.CreateClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewTrustixCombinedRPCClient(conn)

		caches, err := getCaches(c)
		if err != nil {
			panic(err)
		}

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

					storeHash, err := NixB32Encoding.DecodeString(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), ".narinfo"))
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

					{
						tokens := strings.Split(*narinfo.StorePath, "/")
						storeBase := tokens[len(tokens)-1]
						storePrefix := strings.Split(storeBase, "-")[0]

						var narHash string
						{
							tokens := strings.Split(*narinfo.NarHash, ":")
							narHash = tokens[len(tokens)-1]
						}

						w.Header().Set("Content-Type", "text/x-nix-narinfo")
						fmt.Fprintf(w, narinfo.ToString(
							fmt.Sprintf("URL: nar/%s/%s", storePrefix, narHash),
							fmt.Sprintf("Sig: %s:%s", keyPrefix, base64.StdEncoding.EncodeToString(sig)),
						))
					}

					return

				} else if strings.HasPrefix(r.URL.Path, "/nar") {

					var storePrefix string
					var narHash string
					{
						tokens := strings.Split(r.URL.Path, "/")
						if len(tokens) != 4 {
							panic(fmt.Errorf("Malformed URL, expected 4 tokens, got %d", len(tokens)))
						}

						storePrefix = tokens[2]
						narHash = tokens[3]
					}

					for _, cache := range caches {

						var narinfo *nar.NarInfo
						{
							URL, err := url.Parse(cache)
							if err != nil {
								panic(err)
							}
							URL.Path = fmt.Sprintf("%s.narinfo", storePrefix)
							u := URL.String()

							resp, err := http.Get(u)
							if err != nil {
								panic(err)
							}
							if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
								continue
							}

							narinfoBytes, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								panic(err)
							}

							narinfo, err = nar.ParseNarInfo(narinfoBytes)
							if err != nil {
								panic(err)
							}
						}

						if strings.Split(narinfo.NarHash, ":")[1] == narHash {

							URL, err := url.Parse(cache)
							if err != nil {
								panic(err)
							}
							URL.Path = narinfo.URL
							u := URL.String()

							resp, err := http.Get(u)
							if err != nil {
								panic(err)
							}

							if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
								continue
							}

							w.WriteHeader(200)
							w.Header().Add("Content-Type", resp.Header.Get("Content-Type"))

							_, err = io.Copy(w, resp.Body)
							if err != nil {
								panic(err)
							}

							return
						}
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

		l := logrusmiddleware.Middleware{
			Name:   "trustix-binary-cache-proxy",
			Logger: log.New(),
		}

		http.Handle("/", l.Handler(http.HandlerFunc(handler), "/"))
		log.Fatal(http.ListenAndServe(":8080", nil))

		return nil
	},
}
