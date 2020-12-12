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
	"github.com/tweag/trustix/schema"
	"github.com/tweag/trustix/storage"
)

type VerifiableLog struct {
	treeSize uint64
	storage  *LogStorage
}

func NewVerifiableLog(prefix string, txn storage.Transaction, treeSize uint64) (*VerifiableLog, error) {
	storage := NewLogStorage(prefix, txn)
	return &VerifiableLog{
		storage:  storage,
		treeSize: treeSize,
	}, nil
}

func (l *VerifiableLog) Size() uint64 {
	return l.treeSize
}

func (l *VerifiableLog) Root() ([]byte, error) {
	if l.treeSize == 0 {
		return LeafDigest(nil), nil
	}

	level := 0
	for levelSize(l.treeSize, level)%2 == 0 {
		level = level + 1
	}

	lastIndex := levelSize(l.treeSize, level) - 1
	leaf, err := l.storage.Get(level, lastIndex)
	if err != nil {
		return nil, err
	}

	digest := leaf.Digest

	storageSize := rootSize(l.treeSize)
	for i := level + 1; i < storageSize; i++ {
		levelSize := levelSize(l.treeSize, i)
		if levelSize%2 == 1 {
			newLeaf, err := l.storage.Get(i, levelSize-1)
			if err != nil {
				return nil, err
			}

			digest = branchHash(newLeaf.Digest, digest)
		}
	}

	return digest, nil
}

func (l *VerifiableLog) Append(data []byte) (*schema.LogLeaf, error) {
	l.treeSize += 1

	leaf := &schema.LogLeaf{
		Digest: LeafDigest(data),
	}

	err := l.addNodeToLevel(0, leaf)
	if err != nil {
		return nil, err
	}

	return leaf, nil
}

func (l *VerifiableLog) AppendKV(key []byte, value []byte) (*schema.LogLeaf, error) {
	l.treeSize += 1

	leaf := &schema.LogLeaf{
		Key:    key,
		Digest: LeafDigestKV(key, value),
	}

	err := l.addNodeToLevel(0, leaf)
	if err != nil {
		return nil, err
	}

	return leaf, nil
}

func (l *VerifiableLog) addNodeToLevel(level int, leaf *schema.LogLeaf) error {
	err := l.storage.Append(l.treeSize, level, leaf)
	if err != nil {
		return err
	}

	levelSize := levelSize(l.treeSize, level)
	if levelSize%2 == 0 {
		li := levelSize - 2
		ri := levelSize - 1

		ll, err := l.storage.Get(level, li)
		if err != nil {
			return err
		}

		rl, err := l.storage.Get(level, ri)
		if err != nil {
			return err
		}

		newHash := branchHash(ll.Digest, rl.Digest)
		err = l.addNodeToLevel(level+1, &schema.LogLeaf{
			// We don't save the raw value for a branch hash
			Digest: newHash,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *VerifiableLog) AuditProof(node uint64, size uint64) ([][]byte, error) {
	return l.pathFromNodeToRootAtSnapshot(node, 0, size)
}

func (l *VerifiableLog) pathFromNodeToRootAtSnapshot(node uint64, level int, snapshot uint64) ([][]byte, error) {
	var path [][]byte

	if snapshot == 0 {
		return path, nil
	}

	lastNode := snapshot - 1

	lastLeaf, err := l.storage.Get(0, lastNode)
	if err != nil {
		return nil, err
	}

	lastHash := lastLeaf.Digest

	for i := 0; i < level; i++ {
		if isRightChild(lastNode) {
			lastLeaf, err = l.storage.Get(i, lastNode-1)
			if err != nil {
				return nil, err
			}

			lastHash = branchHash(lastLeaf.Digest, lastHash)
		}
		lastNode = parent(lastNode)
	}

	for lastNode > 0 {
		var sibling uint64
		if isRightChild(node) {
			sibling = node - 1
		} else {
			sibling = node + 1
		}

		if sibling < lastNode {
			siblingLeaf, err := l.storage.Get(level, sibling)
			if err != nil {
				return nil, err
			}

			path = append(path, siblingLeaf.Digest)
		} else if sibling == lastNode {
			path = append(path, lastHash)
		}

		if isRightChild(lastNode) {
			lastLeaf, err = l.storage.Get(level, lastNode-1)
			if err != nil {
				return nil, err
			}

			lastHash = branchHash(lastLeaf.Digest, lastHash)
		}
		level += 1
		node = parent(node)
		lastNode = parent(lastNode)
	}

	return path, nil
}

func (l *VerifiableLog) ConsistencyProof(fstSize uint64, sndSize uint64) ([][]byte, error) {
	var proof [][]byte
	if fstSize == 0 || fstSize >= sndSize || sndSize > l.treeSize {
		return proof, nil
	}

	level := 0
	node := fstSize - 1
	for isRightChild(node) {
		node = parent(node)
		level += 1
	}

	if node > 0 {
		leaf, err := l.storage.Get(level, node)
		if err != nil {
			return nil, err
		}
		proof = append(proof, leaf.Digest)
	}

	other, err := l.pathFromNodeToRootAtSnapshot(node, level, sndSize)
	if err != nil {
		return nil, err
	}

	proof = append(proof, other...)
	return proof, nil
}
