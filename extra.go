package main

import (
	"log"
)

// Extra is a type that corresponds to Trillian Log ExtraData
type Extra struct {
	name string
}

func newExtra(name string) *Extra {
	return &Extra{
		name: name,
	}
}

// Marshal Extra.name into []byte
func (e *Extra) Marshal() ([]byte, error) {
	return []byte(e.name), nil
}
