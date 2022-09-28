// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package protocols

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/hashicorp/go-uuid"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
)

type ProtocolDescriptor struct {
	// A globally unique identifier for the protocol
	ID string

	// A human friendly name that can be used in config files and such
	Name string

	NewHash func() hash.Hash
}

var Descriptors = []*ProtocolDescriptor{
	// Internal test
	&ProtocolDescriptor{
		ID:      "cddab738-75cf-4685-94e2-4df58a0f51e7",
		Name:    "test",
		NewHash: sha256.New,
	},
	// Nix protocol
	&ProtocolDescriptor{
		ID:      "5138a791-8d00-4182-96bc-f1f2688cdde2",
		Name:    "nix",
		NewHash: sha256.New,
	},
}

func (pd *ProtocolDescriptor) Validate() error {
	if _, err := uuid.ParseUUID(pd.ID); err != nil {
		return err
	}

	if pd.Name == "" {
		return fmt.Errorf("Protocol descriptor is missing name")
	}

	if pd.NewHash == nil {
		return fmt.Errorf("Hash constructor is nil")
	}

	return nil
}

// LogID - Generate a deterministic ID based on known facts
func (pd *ProtocolDescriptor) LogID(keyType string, publicKey []byte, mode api.Log_LogModes) string {
	h := pd.NewHash()

	writeField := func(fieldName string, data []byte) {
		h.Write([]byte(fieldName + ":"))
		h.Write([]byte(data))
		h.Write([]byte(";"))
	}

	writeField("protocol", []byte(pd.ID))
	writeField("mode", []byte(fmt.Sprintf("%d", mode)))
	writeField("keyType", []byte(keyType))
	writeField("pubKey", publicKey)

	return hex.EncodeToString(h.Sum(nil))
}

func Get(id string) (*ProtocolDescriptor, error) {
	for _, pd := range Descriptors {
		if pd.ID == id || pd.Name == id {
			return pd, pd.Validate()
		}
	}
	return nil, fmt.Errorf("Could not find matching protocol for id: %s", id)
}
