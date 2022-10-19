// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"fmt"
	"os"

	connect "github.com/bufbuild/connect-go"
	"github.com/nix-community/trustix/packages/trustix/auth"
	log "github.com/sirupsen/logrus"
)

func getAuthInterceptors() (connect.Option, error) {
	defaultTokenPath := os.Getenv("TRUSTIX_TOKEN")
	if defaultTokenPath == "" {
		return nil, nil
	}

	f, err := os.Open(defaultTokenPath)
	if err != nil {
		log.Fatalf("Error opening private token file '%s': %v", defaultTokenPath, err)
	}

	tok, err := auth.NewPrivateToken(f)
	if err != nil {
		return nil, fmt.Errorf("Error creating token: %w", err)
	}

	return connect.WithInterceptors(auth.NewAuthInterceptor(tok, nil)), nil
}
