package log

import (
	"crypto/sha256"
)

type VerifiableLog struct {
	treeSize int

	// TODO: Implement persistent storage
	hashes [][]*Leaf
}

func NewVerifiableLog() (*VerifiableLog, error) {
	return &VerifiableLog{
		treeSize: 0,
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
	for len(l.hashes[level])%2 == 0 {
		level = level + 1
	}

	lastIndex := len(l.hashes[level]) - 1
	hash := l.hashes[level][lastIndex].Digest

	for _, hashes := range l.hashes[level+1:] {
		if len(hashes)%2 == 1 {
			idx := len(hashes) - 1
			hash = branchHash(hashes[idx].Digest, hash)
		}
	}

	return hash
}

func (l *VerifiableLog) Append(data []byte) {
	l.treeSize += 1
	h := sha256.New()
	h.Write([]byte{0}) // Write 0x00 prefix
	h.Write(data)

	leaf := &Leaf{
		Value:  data,
		Digest: h.Sum(nil),
	}

	l.addNodeToLevel(0, leaf)
}

func (l *VerifiableLog) addNodeToLevel(level int, leaf *Leaf) {
	if len(l.hashes) == level {
		h := []*Leaf{}
		l.hashes = append(l.hashes, h)
	}

	hashes := l.hashes[level]
	hashes = append(hashes, leaf)
	l.hashes[level] = hashes

	if len(l.hashes[level])%2 == 0 {
		li := len(hashes) - 2
		ri := len(hashes) - 1
		newHash := branchHash(hashes[li].Digest, hashes[ri].Digest)
		l.addNodeToLevel(level+1, &Leaf{
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

	last_node := snapshot - 1
	last_hash := l.hashes[0][last_node].Digest

	for _, row := range l.hashes[:level] {
		if isRightChild(last_node) {
			last_hash = branchHash(row[last_node-1].Digest, last_hash)
		}
		last_node = parent(last_node)
	}

	for last_node > 0 {
		var sibling int
		if isRightChild(node) {
			sibling = node - 1
		} else {
			sibling = node + 1
		}

		if sibling < last_node {
			path = append(path, l.hashes[level][sibling].Digest)
		} else if sibling == last_node {
			path = append(path, last_hash)
		}

		if isRightChild(last_node) {
			last_hash = branchHash(l.hashes[level][last_node-1].Digest, last_hash)
		}
		level += 1
		node = parent(node)
		last_node = parent(last_node)
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
		proof = append(proof, l.hashes[level][node].Digest)
	}

	proof = append(proof, l.pathFromNodeToRootAtSnapshot(node, level, sndSize)...)
	return proof
}
