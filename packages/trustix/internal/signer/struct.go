// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package signer

import (
	"crypto"
)

type Verifier interface {
	Verify(message, sig []byte) bool
	Public() crypto.PublicKey
}
