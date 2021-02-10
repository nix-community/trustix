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
	"encoding/base32"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/api"
	"github.com/tweag/trustix/contrib/nix/schema"
)

var NixB32Encoding = base32.NewEncoding("0123456789abcdfghijklmnpqrsvwxyz")

func createKVPair(storePath string) (*api.KeyValuePair, error) {

	if storePath == "" {
		return nil, fmt.Errorf("Empty input store path")
	}

	tmpDir, err := ioutil.TempDir("", "nix-trustix")
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(tmpDir)
	if err != nil {
		return nil, err
	}

	var storeHash []byte
	{
		storeHashStr := strings.Split(filepath.Base(storePath), "-")[0]
		storeHash, err = NixB32Encoding.DecodeString(storeHashStr)
		if err != nil {
			log.Fatal(err)
		}
		if len(storeHash) == 0 {
			log.Fatal("Empty decoded store path hash")
		}
	}

	var references []string
	{
		out, err := exec.Command("nix-store", "--query", "--references", storePath).Output()
		if err != nil {
			return nil, err
		}

		for _, p := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			tokens := strings.Split(p, "/")
			references = append(references, tokens[len(tokens)-1])
		}
		sort.Strings(references)
	}

	var narHash string
	{
		out, err := exec.Command("nix-store", "--query", "--hash", storePath).Output()
		if err != nil {
			return nil, fmt.Errorf("Could not query hash: %v", err)
		}

		narHash = strings.TrimSpace(string(out))
	}

	var narSize uint64
	{
		out, err := exec.Command("nix-store", "--query", "--size", storePath).Output()
		if err != nil {
			return nil, fmt.Errorf("Could not query size: %v", err)
		}

		narSize, err = strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.WithFields(log.Fields{
		"storePath": storePath,
	}).Debug("Submitting mapping")

	narinfoBytes, err := proto.Marshal(&schema.NarInfo{
		StorePath:  &storePath,
		NarHash:    &narHash,
		NarSize:    &narSize,
		References: references,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &api.KeyValuePair{
		Key:   storeHash,
		Value: narinfoBytes,
	}, nil
}
