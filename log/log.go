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
	"github.com/tweag/trustix/schema"
	"github.com/tweag/trustix/storage"
)

type VerifiableLog struct {
	treeSize int
	storage  *LogStorage
}

func NewVerifiableLog(transaction storage.Transaction, treeSize int) (*VerifiableLog, error) {
	return &VerifiableLog{
		storage: &LogStorage{
			txn: transaction,
		},
		treeSize: treeSize,
	}, nil
}

func (l *VerifiableLog) Size() int {
	return l.treeSize
}

func (l *VerifiableLog) Root() []byte {
	if l.treeSize == 0 {
		h := sha256.New()
		return h.Sum(nil)
	}

	level := 0
	for levelSize(l.treeSize, level)%2 == 0 {
		level = level + 1
	}

	lastIndex := levelSize(l.treeSize, level) - 1
	hash := l.storage.Get(level, lastIndex).Digest

	storageSize := rootSize(l.treeSize)
	for i := level + 1; i < storageSize; i++ {
		levelSize := levelSize(l.treeSize, i)
		if levelSize%2 == 1 {
			hash = branchHash(l.storage.Get(i, levelSize-1).Digest, hash)
		}
	}

	return hash
}

func (l *VerifiableLog) Append(data []byte) {
	l.treeSize += 1
	h := sha256.New()
	h.Write([]byte{0}) // Write 0x00 prefix
	h.Write(data)

	leaf := &schema.LogLeaf{
		Value:  data,
		Digest: h.Sum(nil),
	}

	l.addNodeToLevel(0, leaf)
}

func (l *VerifiableLog) addNodeToLevel(level int, leaf *schema.LogLeaf) {
	l.storage.Append(l.treeSize, level, leaf)

	levelSize := levelSize(l.treeSize, level)
	if levelSize%2 == 0 {
		li := levelSize - 2
		ri := levelSize - 1
		newHash := branchHash(l.storage.Get(level, li).Digest, l.storage.Get(level, ri).Digest)
		l.addNodeToLevel(level+1, &schema.LogLeaf{
			// We don't save the raw value for a branch hash
			Digest: newHash,
		})
	}
}

func (l *VerifiableLog) AuditProof(node int, size int) [][]byte {
	return l.pathFromNodeToRootAtSnapshot(node, 0, size)
}

func (l *VerifiableLog) pathFromNodeToRootAtSnapshot(node int, level int, snapshot int) [][]byte {
	var path [][]byte

	if snapshot == 0 {
		return path
	}

	lastNode := snapshot - 1
	lastHash := l.storage.Get(0, lastNode).Digest

	for i := 0; i < level; i++ {
		if isRightChild(lastNode) {
			lastHash = branchHash(l.storage.Get(i, lastNode-1).Digest, lastHash)
		}
		lastNode = parent(lastNode)
	}

	for lastNode > 0 {
		var sibling int
		if isRightChild(node) {
			sibling = node - 1
		} else {
			sibling = node + 1
		}

		if sibling < lastNode {
			path = append(path, l.storage.Get(level, sibling).Digest)
		} else if sibling == lastNode {
			path = append(path, lastHash)
		}

		if isRightChild(lastNode) {
			lastHash = branchHash(l.storage.Get(level, lastNode-1).Digest, lastHash)
		}
		level += 1
		node = parent(node)
		lastNode = parent(lastNode)
	}

	return path
}

func (l *VerifiableLog) ConsistencyProof(fstSize int, sndSize int) [][]byte {
	var proof [][]byte
	if fstSize == 0 || fstSize >= sndSize || sndSize > l.treeSize {
		return proof
	}

	level := 0
	node := fstSize - 1
	for isRightChild(node) {
		node = parent(node)
		level += 1
	}

	if node > 0 {
		proof = append(proof, l.storage.Get(level, node).Digest)
	}

	proof = append(proof, l.pathFromNodeToRootAtSnapshot(node, level, sndSize)...)
	return proof
}
