// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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
