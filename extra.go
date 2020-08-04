package main

import (
	"log"
)

// Extra is a type that corresponds to Trillian Log ExtraData
type Extra struct {
	name string
}

func newExtra(name string) *Extra {
	log.Printf("[extra:new] Creating: %s", name)
	return &Extra{
		name: name,
	}
}

// Marshal Extra.name into []byte
func (e *Extra) Marshal() ([]byte, error) {
	log.Printf("[extra:marshal] Marshaling: %s", e.name)
	return []byte(e.name), nil
}
