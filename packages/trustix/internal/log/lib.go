// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package log

import (
	"hash"
)

func leafDigest(hashFn func() hash.Hash, data []byte) []byte {
	h := hashFn()
	if data != nil {
		h.Write([]byte{0}) // Write 0x00 prefix
		h.Write(data)
	}
	return h.Sum(nil)
}

func leafDigestKV(hashFn func() hash.Hash, key []byte, value []byte) []byte {
	h := hashFn()
	h.Write([]byte{0}) // Write 0x00 prefix
	h.Write(key)
	h.Write([]byte(":"))
	h.Write(value)
	return h.Sum(nil)
}

func isRightChild(node uint64) bool {
	return node%2 == 1
}

func splitPoint(n uint64) uint64 {
	split := uint64(1)
	for split < n {
		split <<= 1
	}
	return split >> 1
}

func parent(node uint64) uint64 {
	return node / 2
}

func branchHash(hashFn func() hash.Hash, left []byte, right []byte) []byte {
	h := hashFn()
	h.Write([]byte{1}) // Write 0x01 prefix
	h.Write(left)
	h.Write(right)
	return h.Sum(nil)
}

func levelSize(treeSize uint64, level int) uint64 {
	size := treeSize
	for i := 0; i <= level-1; i++ {
		size = size / 2
	}
	return size
}

// How many "buckets" are in the root level for a given tree size
func rootSize(treeSize uint64) int {
	if treeSize == 0 {
		return 0
	}

	size := treeSize
	i := 1
	for size > 0 {
		size = size / 2
		i++
	}
	return i
}
