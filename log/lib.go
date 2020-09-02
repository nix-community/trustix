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
