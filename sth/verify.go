package sth

import (
	"fmt"
	"github.com/tweag/trustix/signer"
)

func verifySMHSig(signer signer.TrustixSigner, smh *SMH) error {

	rootBytes, err := smh.UnmarshalSMHRoot()
	if err != nil {
		return err
	}

	sigBytes, err := smh.UnmarshalSMHSignature()
	if err != nil {
		return err
	}

	if !signer.Verify(rootBytes, sigBytes) {
		return fmt.Errorf("SMH signature verification failed")
	}

	return nil
}

func verifySTHSig(signer signer.TrustixSigner, smh *SMH) error {

	rootBytes, err := smh.UnmarshalSTHRoot()
	if err != nil {
		return err
	}

	sigBytes, err := smh.UnmarshalSTHSignature()
	if err != nil {
		return err
	}

	if !signer.Verify(rootBytes, sigBytes) {
		return fmt.Errorf("SMH signature verification failed")
	}

	return nil
}
