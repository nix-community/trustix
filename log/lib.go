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

package log

import (
	"crypto/sha256"
)

func isRightChild(node int) bool {
	return node%2 == 1
}

func splitPoint(n int) int {
	split := 1
	for split < n {
		split <<= 1
	}
	return split >> 1
}

func parent(node int) int {
	return node / 2
}

func branchHash(left []byte, right []byte) []byte {
	h := sha256.New()
	h.Write([]byte{1}) // Write 0x01 prefix
	h.Write(left)
	h.Write(right)
	return h.Sum(nil)
}

func levelSize(treeSize int, level int) int {
	size := treeSize
	for i := 0; i <= level-1; i++ {
		size = size / 2
	}
	return size
}

// How many "buckets" are in the root level for a given tree size
func rootSize(treeSize int) int {
	size := treeSize
	i := 1
	for {
		size = size / 2
		if size == 0 {
			break
		}
		i++
	}
	return i
}
