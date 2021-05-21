// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package log

import (
	"hash"

	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"github.com/tweag/trustix/packages/trustix/internal/storage"
)

type VerifiableLog struct {
	treeSize uint64
	storage  *LogStorage
	hashFn   func() hash.Hash
}

func NewVerifiableLog(hashFn func() hash.Hash, txn *storage.BucketTransaction, treeSize uint64) (*VerifiableLog, error) {
	storage := NewLogStorage(txn)
	return &VerifiableLog{
		storage:  storage,
		treeSize: treeSize,
		hashFn:   hashFn,
	}, nil
}

func (l *VerifiableLog) Size() uint64 {
	return l.treeSize
}

func (l *VerifiableLog) Root() ([]byte, error) {
	if l.treeSize == 0 {
		return leafDigest(l.hashFn, nil), nil
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

	digest := leaf.LeafDigest

	storageSize := rootSize(l.treeSize)
	for i := level + 1; i < storageSize; i++ {
		levelSize := levelSize(l.treeSize, i)
		if levelSize%2 == 1 {
			newLeaf, err := l.storage.Get(i, levelSize-1)
			if err != nil {
				return nil, err
			}

			digest = branchHash(l.hashFn, newLeaf.LeafDigest, digest)
		}
	}

	return digest, nil
}

func (l *VerifiableLog) Append(data []byte) (*schema.LogLeaf, error) {
	l.treeSize += 1

	leaf := &schema.LogLeaf{
		LeafDigest: leafDigest(l.hashFn, data),
	}

	err := l.addNodeToLevel(0, leaf)
	if err != nil {
		return nil, err
	}

	return leaf, nil
}

func (l *VerifiableLog) AppendKV(key []byte, value []byte) (*schema.LogLeaf, error) {
	l.treeSize += 1

	h := l.hashFn()
	h.Write(value)
	valueDigest := h.Sum(nil)

	leaf := &schema.LogLeaf{
		Key:         key,
		ValueDigest: valueDigest,
		LeafDigest:  leafDigestKV(l.hashFn, key, valueDigest),
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

		newHash := branchHash(l.hashFn, ll.LeafDigest, rl.LeafDigest)
		err = l.addNodeToLevel(level+1, &schema.LogLeaf{
			// We don't save the raw value for a branch hash
			LeafDigest: newHash,
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

	lastHash := lastLeaf.LeafDigest

	for i := 0; i < level; i++ {
		if isRightChild(lastNode) {
			lastLeaf, err = l.storage.Get(i, lastNode-1)
			if err != nil {
				return nil, err
			}

			lastHash = branchHash(l.hashFn, lastLeaf.LeafDigest, lastHash)
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

			path = append(path, siblingLeaf.LeafDigest)
		} else if sibling == lastNode {
			path = append(path, lastHash)
		}

		if isRightChild(lastNode) {
			lastLeaf, err = l.storage.Get(level, lastNode-1)
			if err != nil {
				return nil, err
			}

			lastHash = branchHash(l.hashFn, lastLeaf.LeafDigest, lastHash)
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
		proof = append(proof, leaf.LeafDigest)
	}

	other, err := l.pathFromNodeToRootAtSnapshot(node, level, sndSize)
	if err != nil {
		return nil, err
	}

	proof = append(proof, other...)
	return proof, nil
}
