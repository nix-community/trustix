package sth

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/lazyledger/smt"
	vlog "github.com/tweag/trustix/log"
	"github.com/tweag/trustix/signer"
)

type STH struct {
	Signature string `json:"sig"`
	Root      string `json:"root"`
	Size      int    `json:"tree_size"`
}

type SMH struct {
	Signature string `json:"sig"`
	Root      string `json:"root"`
	LogSth    *STH   `json:"log_sth"`
}

func SignHead(smTree *smt.SparseMerkleTree, vLog *vlog.VerifiableLog, signer crypto.Signer) (*SMH, error) {

	opts := crypto.SignerOpts(crypto.Hash(0))

	smTreeRoot := smTree.Root()
	smTreeSig, err := signer.Sign(rand.Reader, smTreeRoot, opts)
	if err != nil {
		return nil, err
	}

	vLogRoot := vLog.Root()
	sthSig, err := signer.Sign(rand.Reader, vLogRoot, opts)
	if err != nil {
		return nil, err
	}

	return &SMH{
		Signature: base64.StdEncoding.EncodeToString(smTreeSig),
		Root:      base64.StdEncoding.EncodeToString(smTreeRoot),
		LogSth: &STH{
			Signature: base64.StdEncoding.EncodeToString(sthSig),
			Root:      base64.StdEncoding.EncodeToString(vLogRoot),
			Size:      vLog.Size(),
		},
	}, nil
}

func NewSMHFromJSON(j []byte) (*SMH, error) {
	smh := &SMH{}
	err := json.Unmarshal(j, &smh)
	if err != nil {
		return nil, err
	}
	return smh, nil
}

func (smh *SMH) Verify(signer signer.TrustixSigner) error {
	err := verifySMHSig(signer, smh)
	if err != nil {
		return err
	}
	return verifySTHSig(signer, smh)
}

func (smh *SMH) UnmarshalSMHSignature() ([]byte, error) {
	return base64.StdEncoding.DecodeString(smh.Signature)
}

func (smh *SMH) UnmarshalSTHSignature() ([]byte, error) {
	return base64.StdEncoding.DecodeString(smh.LogSth.Signature)
}

func (smh *SMH) UnmarshalSMHRoot() ([]byte, error) {
	return base64.StdEncoding.DecodeString(smh.Root)
}

func (smh *SMH) UnmarshalSTHRoot() ([]byte, error) {
	return base64.StdEncoding.DecodeString(smh.LogSth.Root)
}

func (smh *SMH) Marshal() ([]byte, error) {
	return json.Marshal(smh)
}
