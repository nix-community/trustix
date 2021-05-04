// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package log

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tweag/trustix/packages/trustix/storage"
)

type testInput struct {
	Input        []byte
	ExpectedRoot []byte
}

type testStorageTxn struct {
	// bucket -> key -> value
	kv map[string]map[string][]byte
}

func newTestStorageTxn() *testStorageTxn {
	return &testStorageTxn{
		kv: make(map[string]map[string][]byte),
	}
}

func (s *testStorageTxn) Get(bucket *storage.Bucket, key []byte) ([]byte, error) {
	b, ok := s.kv[bucket.Join()]
	if !ok {
		return nil, storage.ObjectNotFoundError
	}

	m, ok := b[string(key)]
	if !ok {
		return nil, storage.ObjectNotFoundError
	}

	return m, nil
}
func (s *testStorageTxn) Set(bucket *storage.Bucket, key []byte, value []byte) error {
	_, ok := s.kv[bucket.Join()]
	if !ok {
		s.kv[bucket.Join()] = make(map[string][]byte)
	}
	s.kv[bucket.Join()][string(key)] = value
	return nil
}
func (s *testStorageTxn) Delete(bucket *storage.Bucket, key []byte) error {
	return fmt.Errorf("Not implemented")
}

func newTestStorageBucketTxn() *storage.BucketTransaction {
	b := &storage.Bucket{}
	t, _ := b.Cd("test").Txn(newTestStorageTxn())
	return t
}

func mkInputs() []*testInput {

	decodeHex := func(h string) []byte {
		b, err := hex.DecodeString(h)
		if err != nil {
			panic(err)
		}
		return b
	}

	return []*testInput{
		&testInput{
			Input:        []byte(""),
			ExpectedRoot: decodeHex("6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d"),
		},
		&testInput{
			Input:        []byte("a"),
			ExpectedRoot: decodeHex("688dc6244b041199e7ab4990df6340ce3dc14caa5cd5a0e1131addaa1209e1a6"),
		},
		&testInput{
			Input:        []byte("b"),
			ExpectedRoot: decodeHex("652297b9504045a600942bcdf9ae5c2400be42d51139c7fb63ab3ee439ff110d"),
		},
		&testInput{
			Input:        []byte("c"),
			ExpectedRoot: decodeHex("4a9bab0b70e36b453e967468fc209705d9171fd05e9cf9e0ed6c2dff673fc790"),
		},
		&testInput{
			Input:        []byte("d"),
			ExpectedRoot: decodeHex("ff9cdaec73345d3896e37ff5681084b7be4097839f760e621412f9343d139f22"),
		},
		&testInput{
			Input:        []byte("efghijk"),
			ExpectedRoot: decodeHex("a44a4f5f5190f8bf6acbfecc50e56374072196c17aa5fd46af01a5b9674307cf"),
		},
		&testInput{
			Input:        []byte("lmnopqrstuvwxyz"),
			ExpectedRoot: decodeHex("968244ebd454ce024d380be757b570886f8449f41395c761ec363a08a8f18210"),
		},
	}
}

func mkAssertProof(t *testing.T, proofFunc func(first uint64, second uint64) (proof [][]byte, err error)) func(first int, second int, expected []string) {
	assert := assert.New(t)
	return func(first int, second int, expected []string) {
		proof, err := proofFunc(uint64(first), uint64(second))
		assert.Nil(err)

		assert.Equal(len(expected), len(proof), "Proof length matches expected length")
		for i, proof := range proof {
			hexProof := hex.EncodeToString(proof)
			assert.Equal(expected[i], hexProof, "Proof matches expected")
		}
	}
}

func TestLogRoots(t *testing.T) {

	assert := assert.New(t)

	storageTxn := newTestStorageBucketTxn()

	tree, err := NewVerifiableLog(storageTxn, 0)
	assert.Nil(err)

	root, err := tree.Root()
	assert.Nil(err)

	assert.Equal("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", hex.EncodeToString(root), fmt.Sprintf("Correct root for zero inputs"))

	for _, input := range mkInputs() {
		_, err = tree.Append(input.Input)
		if err != nil {
			t.Fatalf("tree.Append failed: %v", err)
		}

		root, err := tree.Root()
		assert.Nil(err)
		// Encode to hex to get prettier error message
		assert.Equal(hex.EncodeToString(input.ExpectedRoot), hex.EncodeToString(root), fmt.Sprintf("Correct root for input %s", input.Input))

	}

}

func TestAuditProofs(t *testing.T) {

	assert := assert.New(t)

	storageTxn := newTestStorageBucketTxn()

	tree, err := NewVerifiableLog(storageTxn, 0)
	assert.Nil(err)

	inputs := mkInputs()
	for _, input := range inputs {
		_, err = tree.Append(input.Input)
		if err != nil {
			t.Fatalf("tree.Append failed: %v", err)
		}
	}

	assert.Equal(7, len(inputs), "Assert expected inputs to test")

	assertProof := mkAssertProof(t, tree.AuditProof)

	assertProof(0, 0, []string{})

	assertProof(1, 2, []string{
		"6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d",
	})

	assertProof(0, 4, []string{
		"022a6979e6dab7aa5ae4c3e5e45f7e977112a7e63593820dbec1ec738a24f93c",
		"a5eb6e7bcfaaff4957c342e0cbfe88209dbe2058fc3e1a3455cc071922c85741",
	})

	assertProof(4, 7, []string{
		"bc78703cecc49c1119746b7baa573cc39274b72f40fe1e6c242fc524b1acd3f6",
		"e6d714a0c30dbe89616ee317930e7821a18f18c5a80307e08fc92e7809e52d86",
		"4a9bab0b70e36b453e967468fc209705d9171fd05e9cf9e0ed6c2dff673fc790",
	})

	assertProof(3, 5, []string{
		"57eb35615d47f34ec714cacdf5fd74608a5e8e102724e80b24b287c0c27b6a31",
		"688dc6244b041199e7ab4990df6340ce3dc14caa5cd5a0e1131addaa1209e1a6",
		"d070dc5b8da9aea7dc0f5ad4c29d89965200059c9a0ceca3abd5da2492dcb71d",
	})

	assertProof(0, 7, []string{
		"022a6979e6dab7aa5ae4c3e5e45f7e977112a7e63593820dbec1ec738a24f93c",
		"a5eb6e7bcfaaff4957c342e0cbfe88209dbe2058fc3e1a3455cc071922c85741",
		"49ad1f129f0f126dd6b90955fb177ab8941be0d7b5d0085c4813fcabb62b6ec9",
	})

}

func TestConsistencyProofs(t *testing.T) {

	assert := assert.New(t)

	storageTxn := newTestStorageBucketTxn()

	tree, err := NewVerifiableLog(storageTxn, 0)
	assert.Nil(err)

	assertProof := mkAssertProof(t, tree.ConsistencyProof)

	inputs := mkInputs()
	for _, input := range inputs {
		_, err = tree.Append(input.Input)
		if err != nil {
			t.Fatalf("tree.Append failed: %v", err)
		}
	}

	assert.Equal(7, len(inputs), "Assert expected inputs to test")

	assertProof(1, 1, []string{})

	assertProof(2, 5, []string{
		"a5eb6e7bcfaaff4957c342e0cbfe88209dbe2058fc3e1a3455cc071922c85741",
		"d070dc5b8da9aea7dc0f5ad4c29d89965200059c9a0ceca3abd5da2492dcb71d",
	})

	assertProof(1, 7, []string{
		"022a6979e6dab7aa5ae4c3e5e45f7e977112a7e63593820dbec1ec738a24f93c",
		"a5eb6e7bcfaaff4957c342e0cbfe88209dbe2058fc3e1a3455cc071922c85741",
		"49ad1f129f0f126dd6b90955fb177ab8941be0d7b5d0085c4813fcabb62b6ec9",
	})

}
