// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

	connect "github.com/bufbuild/connect-go"
)

const tokenHeader = "Trustix-Token"

type authHeader struct {
	Name string `json:"name"` // name of the key, used for looking up key in map
	Sig  string `json:"sig"`  // Signature (base64 encoded)
}

func fmtAuthHeaderMessage(name string, procedure string) []byte {
	return []byte(name + ":" + procedure)
}

func NewAuthInterceptor(privateToken *PrivateToken, writeTokens map[string]*PublicToken) connect.UnaryInterceptorFunc {

	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			spec := req.Spec()
			procedure := spec.Procedure

			req.Header().Set("X-TRUSTIX-AUTH", "0")

			if req.Spec().IsClient && privateToken != nil {
				sig, err := privateToken.Sign(fmtAuthHeaderMessage(privateToken.Name, procedure))
				if err != nil {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("signing request failed"),
					)
				}

				hdr, err := json.Marshal(&authHeader{
					Name: privateToken.Name,
					Sig:  base64.StdEncoding.EncodeToString(sig),
				})
				if err != nil {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("encoding signed json auth header failed"),
					)
				}

				req.Header().Set(tokenHeader, string(hdr))

			} else {
				tokenStr := req.Header().Get(tokenHeader)
				if tokenStr == "" {
					return next(ctx, req)
				}

				hdr := &authHeader{}
				err := json.Unmarshal([]byte(tokenStr), hdr)
				if err != nil {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("error decoding auth header"),
					)
				}

				sig, err := base64.StdEncoding.DecodeString(hdr.Sig)
				if err != nil {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("error decoding signature"),
					)
				}

				tok, ok := writeTokens[hdr.Name]
				if !ok {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("request signed with unknown key"),
					)
				}

				if !tok.Verify(fmtAuthHeaderMessage(hdr.Name, procedure), sig) {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("signature verification failed"),
					)
				}

				// Propagate that we can write
				req.Header().Set("X-TRUSTIX-AUTH", "1")
			}

			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
