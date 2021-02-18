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

package schema

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
)

type NarInfo struct {
	StorePath  string   `json:"path"`
	NarHash    string   `json:"narHash"`
	NarSize    uint64   `json:"narSize"`
	References []string `json:"references"`
}

func (n *NarInfo) Fingerprint() []byte {
	var b bytes.Buffer

	b.WriteString("1;")
	b.WriteString(n.StorePath)
	b.WriteString(";")

	// TODO: Verify whether to strip prefix or not
	b.WriteString(n.NarHash)
	b.WriteString(";")

	b.WriteString(strconv.FormatUint(n.NarSize, 10))
	b.WriteString(";")

	b.WriteString(strings.Join(n.References, ","))

	return b.Bytes()
}

func (n *NarInfo) Sign(signer crypto.Signer) ([]byte, error) {
	opts := crypto.SignerOpts(crypto.Hash(0))
	return signer.Sign(rand.Reader, n.Fingerprint(), opts)
}

func (n *NarInfo) ToString(extraLines ...string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("StorePath: %s\n", n.StorePath))
	b.WriteString(fmt.Sprintf("Compression: %s\n", "none"))
	b.WriteString(fmt.Sprintf("FileHash: %s\n", n.NarHash))
	b.WriteString(fmt.Sprintf("FileSize: %d\n", n.NarSize))
	b.WriteString(fmt.Sprintf("NarHash: %s\n", n.NarHash))
	b.WriteString(fmt.Sprintf("NarSize: %d\n", n.NarSize))
	b.WriteString(fmt.Sprintf("References: %s\n", strings.Join(n.References, " ")))

	for _, l := range extraLines {
		b.WriteString(l)
		b.WriteString("\n")
	}

	return b.String()
}
