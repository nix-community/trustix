package signer

import (
	"crypto"
	"fmt"
	"io"
)

type pubkeySigner struct {
	pub      crypto.PublicKey
	verifier func(message, sig []byte) bool
}

func (s *pubkeySigner) Public() crypto.PublicKey {
	return s.pub
}

func (s *pubkeySigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return nil, fmt.Errorf("This Signer implementation can only be used for verification, as we lack a private key")
}

func (s *pubkeySigner) CanSign() bool {
	return false
}

func (s *pubkeySigner) Verify(message, sig []byte) bool {
	return s.verifier(message, sig)
}
