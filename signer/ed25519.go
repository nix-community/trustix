package signer

import (
	"crypto/ed25519"
	"github.com/tweag/trustix/config"
)

func newED25519Verifier(pubkey ed25519.PublicKey) func(message, sig []byte) bool {
	return func(message, sig []byte) bool {
		return ed25519.Verify(pubkey, message, sig)
	}
}

func genED25519Signer(signerConfig *config.SignerConfig) (TrustixSigner, error) {
	pubBytes, err := decodeKey(signerConfig.PublicKey)
	if err != nil {
		return nil, err
	}
	pub := ed25519.PublicKey(pubBytes)

	// Public key is derived from private key for ed25519
	if signerConfig.ED25519 != nil && signerConfig.ED25519.PrivateKeyPath != "" {
		privBytes, err := readKey(signerConfig.ED25519.PrivateKeyPath)
		if err != nil {
			return nil, err
		}

		key := ed25519.NewKeyFromSeed(privBytes[:32])

		return &privkeySigner{
			signer:   key,
			verifier: newED25519Verifier(pub),
		}, nil

	} else {
		return &pubkeySigner{
			pub:      pub,
			verifier: newED25519Verifier(pub),
		}, nil
	}
}
