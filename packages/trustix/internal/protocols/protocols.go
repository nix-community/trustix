// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package protocols

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"sync"

	"github.com/hashicorp/go-uuid"
)

type ProtocolDescriptor struct {
	// A globally unique identifier for the protocol
	ID string

	// A human friendly name that can be used in config files and such
	Name string

	NewHash func() hash.Hash
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
func (pd *ProtocolDescriptor) LogID(keyType string, publicKey []byte) string {
	h := pd.NewHash()

	h.Write([]byte(keyType))
	h.Write([]byte(":"))

	h.Write(publicKey)
	h.Write([]byte(":"))

	return hex.EncodeToString(h.Sum(nil))
}

var once sync.Once

var descriptors = []*ProtocolDescriptor{
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

func Get(id string) (*ProtocolDescriptor, error) {
	for _, pd := range descriptors {
		if pd.ID == id || pd.Name == id {
			return pd, pd.Validate()
		}
	}
	return nil, fmt.Errorf("Could not find matching protocol for id: %s", id)
}
