package main

import (
	"log"
)

// Thing is a type that corresponds to Trillian Log LeafValue
type Thing struct {
	name string
}

func newThing(name string) *Thing {
	log.Printf("[thing:new] Creating: %s", name)
	return &Thing{
		name: name,
	}
}

// Marshal Thing.name into []byte
func (t *Thing) Marshal() ([]byte, error) {
	log.Printf("[thing:marshal] Marshaling: %s", t.name)
	return []byte(t.name), nil
}
