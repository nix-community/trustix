package signer

import (
	"crypto"
	"io"
)

type privkeySigner struct {
	signer crypto.Signer // Underlying signer
}

func (s *privkeySigner) Public() crypto.PublicKey {
	return s.signer.Public()
}

func (s *privkeySigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return s.signer.Sign(rand, digest, opts)
}

func (s *privkeySigner) CanSign() bool {
	return true
}
