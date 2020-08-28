package main

import (
	"encoding/hex"
)

// Input - Corresponds to Trillian Log LeafValue
type Input struct {
	inputHash  []byte
	outputHash []byte
}

// Create a new Input instance from inputHash string
func newInput(inputHash string, outputHash string) *Input {
	i, _ := hex.DecodeString(inputHash)
	o, _ := hex.DecodeString(outputHash)

	return &Input{
		inputHash:  i,
		outputHash: o,
	}
}

func (t *Input) IdentityHash() []byte {
	return t.inputHash[:]
}

func (t *Input) OutputHash() []byte {
	return t.outputHash[:]
}

func (t *Input) IdentityHashString() string {
	return hex.EncodeToString(t.inputHash)
}

func (t *Input) OutputHashString() string {
	return hex.EncodeToString(t.outputHash)
}
