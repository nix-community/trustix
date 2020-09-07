package sth

import (
	"encoding/base64"
	"encoding/json"
)

type STH struct {
	Signature string `json:"sig"`
	Root      string `json:"root"`
}

func newSTH(root []byte, sig []byte) *STH {
	return &STH{
		Signature: base64.StdEncoding.EncodeToString(sig),
		Root:      base64.StdEncoding.EncodeToString(root),
	}
}

func (sth *STH) FromJSON(j []byte) error {
	return json.Unmarshal(j, &sth)
}

func (sth *STH) UnmarshalSignature() ([]byte, error) {
	return base64.StdEncoding.DecodeString(sth.Signature)
}

func (sth *STH) UnmarshalRoot() ([]byte, error) {
	return base64.StdEncoding.DecodeString(sth.Root)
}

func (sth *STH) Marshal() ([]byte, error) {
	return json.Marshal(sth)
}
