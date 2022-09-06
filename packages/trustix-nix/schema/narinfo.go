// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

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
