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
	"bytes"
	"crypto/sha256"
	"fmt"
)

func rootHashFromAuditProof(leafHash []byte, proof [][]byte, idx int, treeSize int) ([]byte, error) {
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

func rootHashFromConsistencyProof(oldSize int, newSize int, proofNodes [][]byte, oldRoot []byte, computeNewRoot bool, startFromOldRoot bool) []byte {
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

func ValidAuditProof(rootHash []byte, treeSize int, idx int, proof [][]byte, leafData []byte) (bool, error) {
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

func ValidConsistencyProof(oldRoot []byte, newRoot []byte, oldSize int, newSize int, proofNodes [][]byte) bool {
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
