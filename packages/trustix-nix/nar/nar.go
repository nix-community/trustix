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

package nar

import (
	"fmt"
	"strconv"
	"strings"
)

type NarInfo struct {
	StorePath   string
	URL         string
	Compression string
	FileHash    string
	FileSize    uint64
	NarHash     string
	NarSize     uint64
	References  []string
	Deriver     string
	Sig         string
	System      string
	CA          string
}

func ParseNarInfo(input []byte) (*NarInfo, error) {

	n := &NarInfo{}

	parseUint := func(in string) (uint64, error) {
		return strconv.ParseUint(in, 10, 64)
	}

	var err error

	for _, line := range strings.Split(string(input), "\n") {
		tok := strings.SplitN(line, ": ", 2)
		if len(tok) <= 1 {
			continue
		}

		if len(tok) != 2 {
			return nil, fmt.Errorf("Unexpected number of tokens '%d' for value '%s'", len(tok), tok)
		}
		value := tok[1]

		switch tok[0] {
		case "StorePath":
			n.StorePath = value
		case "URL":
			n.URL = value
		case "Compression":
			n.Compression = value
		case "FileHash":
			n.FileHash = value
		case "FileSize":
			n.FileSize, err = parseUint(value)
			if err != nil {
				return nil, err
			}
		case "NarHash":
			n.NarHash = value
		case "NarSize":
			n.NarSize, err = parseUint(value)
			if err != nil {
				return nil, err
			}
		case "References":
			n.References = strings.Split(value, " ")
		case "Deriver":
			n.Deriver = value
		case "Sig":
			n.Sig = value
		case "System":
			n.System = value
		case "CA":
			n.CA = value
		default:
			return nil, fmt.Errorf("Unknown field: %s", tok[0])
		}
	}

	return n, nil
}
