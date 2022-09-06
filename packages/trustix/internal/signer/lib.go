// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package signer

import (
	"encoding/base64"
	"io/ioutil"
)

func readKey(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(data))
}
