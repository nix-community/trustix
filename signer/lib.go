package signer

import (
	"encoding/base64"
	"io/ioutil"
)

func decodeKey(b64Key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64Key)
}

func readKey(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return decodeKey(string(data))
}
