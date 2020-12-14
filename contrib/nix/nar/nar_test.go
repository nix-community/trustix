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
	"testing"

	"github.com/stretchr/testify/assert"
)

const testInput = `
StorePath: /nix/store/byghkc1k0xmrrl2jk04lp0qipmpmz547-hello-2.10
URL: nar/06vv5hjmdyrklwsxxq5d4fnslkgfzpy3z3ri4s7a9fawi2d20ivb.nar.xz
Compression: xz
FileHash: sha256:06vv5hjmdyrklwsxxq5d4fnslkgfzpy3z3ri4s7a9fawi2d20ivb
FileSize: 41272
NarHash: sha256:1llxabk0xq0gc15yi6kkysfbvn5gzisj9dxk6g29sh5ncqx3if8y
NarSize: 206000
References: a6rnjp15qgp8a699dlffqj94hzy1nldg-glibc-2.32 byghkc1k0xmrrl2jk04lp0qipmpmz547-hello-2.10
Deriver: m0i10ghpcwhi2dml0dj6b437jjrh8ia3-hello-2.10.drv
Sig: cache.nixos.org-1:IHkSz9VMQC/KGYgah2Vr2ISz0uawXUKqm/yP4JtcaBkyLO13B3yD2k578ZrP3RyJIVyvMdn4KOjUxvCungdkDA==
`

func TestParseNarInfo(t *testing.T) {

	assert := assert.New(t)

	n, err := ParseNarInfo(testInput)
	assert.Nil(err)
	assert.NotNil(n)

	assert.Equal("/nix/store/byghkc1k0xmrrl2jk04lp0qipmpmz547-hello-2.10", *n.StorePath)
	assert.Equal("nar/06vv5hjmdyrklwsxxq5d4fnslkgfzpy3z3ri4s7a9fawi2d20ivb.nar.xz", *n.URL)
	assert.Equal("xz", *n.Compression)
	assert.Equal("sha256:06vv5hjmdyrklwsxxq5d4fnslkgfzpy3z3ri4s7a9fawi2d20ivb", *n.FileHash)
	assert.Equal(uint64(41272), *n.FileSize)
	assert.Equal("sha256:1llxabk0xq0gc15yi6kkysfbvn5gzisj9dxk6g29sh5ncqx3if8y", *n.NarHash)
	assert.Equal(uint64(206000), *n.NarSize)
	assert.Equal(2, len(n.References))

}
