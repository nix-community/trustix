// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/bakins/logrus-middleware"
	"github.com/coreos/go-systemd/activation"
	"github.com/nix-community/trustix/packages/trustix-nix/nar"
	"github.com/nix-community/trustix/packages/trustix-nix/schema"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	pb "github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix/client"
	"github.com/nix-community/trustix/packages/trustix/interfaces"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ulikunitz/xz"
)

var nixProtocolId string = "5138a791-8d00-4182-96bc-f1f2688cdde2"

var listenAddresses []string
var binaryCachePrivKey string

func getCaches(c interfaces.RpcAPI) ([]string, error) {
	ctx, cancel := client.CreateContext(30)
	defer cancel()
	resp, err := c.Logs(ctx, &api.LogsRequest{})
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

		if binaryCachePrivKey == "" {
			log.Fatalf("Missing required binary cache key")
		}

		keyPrefix, signer, err := readKey(binaryCachePrivKey)
		if err != nil {
			panic(err)
		}

		c, err := client.CreateClient(dialAddress)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer c.Close()

		caches, err := getCaches(c.RpcAPI)
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

					resp, err := c.RpcAPI.Decide(r.Context(), &pb.DecideRequest{
						Key:      storeHash,
						Protocol: &nixProtocolId,
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
					err = json.Unmarshal(resp.Decision.Value, narinfo)
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
						tokens := strings.Split(narinfo.StorePath, "/")
						storeBase := tokens[len(tokens)-1]
						storePrefix := strings.Split(storeBase, "-")[0]

						var narHash string
						{
							tokens := strings.Split(narinfo.NarHash, ":")
							narHash = tokens[len(tokens)-1]
						}

						w.Header().Set("Content-Type", "text/x-nix-narinfo")
						fmt.Fprint(w, narinfo.ToString(
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

							var responseReader io.Reader
							switch narinfo.Compression {
							case "none":
								responseReader = resp.Body
							case "xz":
								responseReader, err = xz.NewReader(resp.Body)
								if err != nil {
									w.WriteHeader(500)
									panic(err)
								}
							default:
								w.WriteHeader(500)
								panic(fmt.Errorf("Unhandled NAR compression '%s'", narinfo.Compression))
							}

							w.WriteHeader(200)
							w.Header().Add("Content-Type", resp.Header.Get("Content-Type"))
							_, err = io.Copy(w, responseReader)
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

		loggedHandler := l.Handler(http.HandlerFunc(handler), "/")

		var listeners []net.Listener

		{
			systemdListeners, err := activation.Listeners()
			if err != nil {
				panic(err)
			}

			for _, lis := range systemdListeners {
				log.WithFields(log.Fields{
					"address": lis.Addr(),
				}).Info("Using socket activated listener")

				listeners = append(listeners, lis)
			}
		}

		for _, addr := range listenAddresses {
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			log.WithFields(log.Fields{
				"address": addr,
			}).Info("Listening to address")

			listeners = append(listeners, lis)
		}

		if len(listeners) == 0 {
			panic(fmt.Errorf("No listeners configured"))
		}

		errChan := make(chan error)

		for _, listener := range listeners {
			go func(l net.Listener) {
				err := http.Serve(l, loggedHandler)
				if err != nil {
					errChan <- err
				}
			}(listener)
		}
		for err := range errChan {
			panic(err)
		}

		return nil
	},
}

func initBinaryCache() {
	binaryCacheCommand.Flags().StringVar(&binaryCachePrivKey, "privkey", "", "Binary cache private key (generated by nix-store)")

	binaryCacheCommand.Flags().StringSliceVar(&listenAddresses, "listen", []string{}, "Listen to address")
}
