// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package log

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

func rootHashFromAuditProof(leafHash []byte, proof [][]byte, idx uint64, treeSize uint64) ([]byte, error) {
	if len(proof) == 0 {
		return leafHash, nil
	}

	if idx%2 == 0 && idx+1 == treeSize {
		if treeSize == 1 {
			return nil, fmt.Errorf("No such level")
		}
		return rootHashFromAuditProof(leafHash, proof, idx/2, (treeSize+1)/2)
	}

	sibling := proof[0]
	if idx%2 == 0 {
		return rootHashFromAuditProof(branchHash(leafHash, sibling), proof, idx/2, (treeSize+1)/2)
	} else {
		return rootHashFromAuditProof(branchHash(sibling, leafHash), proof, idx/2, (treeSize+1)/2)
	}
}

func rootHashFromConsistencyProof(oldSize uint64, newSize uint64, proofNodes [][]byte, oldRoot []byte, computeNewRoot bool, startFromOldRoot bool) []byte {
	if oldSize == newSize {
		if startFromOldRoot {
			return oldRoot
		}
		idx := len(proofNodes) - 1
		return proofNodes[idx]
	}

	k := splitPoint(newSize)
	idx := len(proofNodes) - 1
	nextHash := proofNodes[idx]

	if oldSize <= k {
		leftChild := rootHashFromConsistencyProof(oldSize, k, proofNodes[:idx], oldRoot, computeNewRoot, startFromOldRoot)
		if computeNewRoot {
			return branchHash(leftChild, nextHash)
		} else {
			return leftChild
		}
	} else {
		rightChild := rootHashFromConsistencyProof(oldSize-k, newSize-k, proofNodes[:idx], oldRoot, computeNewRoot, false)
		return branchHash(nextHash, rightChild)
	}

}

func ValidAuditProof(rootHash []byte, treeSize uint64, idx uint64, proof [][]byte, leafData []byte) (bool, error) {
	leafHash := sha256.New()
	leafHash.Write([]byte{0})
	leafHash.Write(leafData)

	fromAuditProof, err := rootHashFromAuditProof(
		leafHash.Sum(nil),
		proof,
		idx,
		treeSize)
	if err != nil {
		return false, err
	}

	return bytes.Compare(rootHash, fromAuditProof) == 0, nil
}

func ValidConsistencyProof(oldRoot []byte, newRoot []byte, oldSize uint64, newSize uint64, proofNodes [][]byte) bool {
	if oldSize == 0 { // Empty tree consistent with any future state
		return true
	}

	if oldSize == newSize {
		return bytes.Compare(oldRoot, newRoot) == 0
	}

	computedOldRoot := rootHashFromConsistencyProof(oldSize, newSize, proofNodes, oldRoot, false, true)
	computedNewRoot := rootHashFromConsistencyProof(oldSize, newSize, proofNodes, oldRoot, true, true)

	return bytes.Compare(oldRoot, computedOldRoot) == 0 && bytes.Compare(newRoot, computedNewRoot) == 0
}
