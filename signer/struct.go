package signer

import (
	"crypto"
	"io"
)

// Implements crypto.Signer
// Extend with CanSign to know if it can be used only for verification or for actual signing
type TrustixSigner interface {
	Public() crypto.PublicKey
	Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error)
	CanSign() bool
}
