package signer

import (
	"fmt"
	"github.com/tweag/trustix/config"
)

func FromConfig(signerConfig *config.SignerConfig) (TrustixSigner, error) {

	if signerConfig.Type == "" {
		fmt.Errorf("Missing signer config field 'type'.", signerConfig.Type)
	}

	switch signerConfig.Type {
	case "ed25519":
		return genED25519Signer(signerConfig)
	}

	return nil, fmt.Errorf("Signer type '%s' is not supported.", signerConfig.Type)

}
