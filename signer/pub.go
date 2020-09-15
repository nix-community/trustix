package signer

import (
	"crypto"
	"fmt"
	"io"
)

type pubkeySigner struct {
	pub crypto.PublicKey
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
