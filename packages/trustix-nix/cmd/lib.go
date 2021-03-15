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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/packages/trustix-nix/schema"
	"github.com/tweag/trustix/packages/trustix-proto/api"
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

	var narinfo *schema.NarInfo
	{
		out, err := exec.Command("nix", "path-info", "--json", storePath).Output()
		if err != nil {
			return nil, err
		}

		var narinfos []*schema.NarInfo
		err = json.Unmarshal(out, &narinfos)
		if err != nil {
			log.Fatalf("Could not get path info: %v", err)
		}

		if len(narinfos) != 1 {
			log.Fatalf("Unexpected number of narinfos returned: %d", len(narinfos))
		}

		narinfo = narinfos[0]

		sort.Strings(narinfo.References)
	}

	log.WithFields(log.Fields{
		"storePath": storePath,
	}).Debug("Submitting mapping")

	narinfoBytes, err := json.Marshal(narinfo)
	if err != nil {
		log.Fatalf("Could not marshal narinfo: %v", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return &api.KeyValuePair{
		Key:   storeHash,
		Value: narinfoBytes,
	}, nil
}
