package signer

import (
	"crypto/ed25519"
	"github.com/tweag/trustix/config"
)

func genED25519Signer(signerConfig *config.SignerConfig) (TrustixSigner, error) {

	// Public key is derived from private key for ed25519
	if signerConfig.ED25519 != nil && signerConfig.ED25519.PrivateKeyPath != "" {
		privBytes, err := readKey(signerConfig.ED25519.PrivateKeyPath)
		if err != nil {
			return nil, err
		}

		key := ed25519.NewKeyFromSeed(privBytes[:32])

		return &privkeySigner{
			signer: key,
		}, nil

	} else {
		pubBytes, err := decodeKey(signerConfig.PublicKey)
		if err != nil {
			return nil, err
		}
		return &pubkeySigner{
			pub: ed25519.PublicKey(pubBytes),
		}, nil

	}
}
