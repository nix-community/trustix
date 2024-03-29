// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: MIT

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: reprod-api/api.proto

package reprod_apiconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	reprod_api "github.com/nix-community/trustix/packages/trustix-nix-r13y/reprod-api"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// ReproducibilityAPIName is the fully-qualified name of the ReproducibilityAPI service.
	ReproducibilityAPIName = "reprod_api.v1.ReproducibilityAPI"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ReproducibilityAPIDerivationReproducibilityProcedure is the fully-qualified name of the
	// ReproducibilityAPI's DerivationReproducibility RPC.
	ReproducibilityAPIDerivationReproducibilityProcedure = "/reprod_api.v1.ReproducibilityAPI/DerivationReproducibility"
	// ReproducibilityAPIAttrReproducibilityTimeSeriesProcedure is the fully-qualified name of the
	// ReproducibilityAPI's AttrReproducibilityTimeSeries RPC.
	ReproducibilityAPIAttrReproducibilityTimeSeriesProcedure = "/reprod_api.v1.ReproducibilityAPI/AttrReproducibilityTimeSeries"
	// ReproducibilityAPIAttrReproducibilityTimeSeriesGroupedbyChannelProcedure is the fully-qualified
	// name of the ReproducibilityAPI's AttrReproducibilityTimeSeriesGroupedbyChannel RPC.
	ReproducibilityAPIAttrReproducibilityTimeSeriesGroupedbyChannelProcedure = "/reprod_api.v1.ReproducibilityAPI/AttrReproducibilityTimeSeriesGroupedbyChannel"
	// ReproducibilityAPISuggestAttributeProcedure is the fully-qualified name of the
	// ReproducibilityAPI's SuggestAttribute RPC.
	ReproducibilityAPISuggestAttributeProcedure = "/reprod_api.v1.ReproducibilityAPI/SuggestAttribute"
	// ReproducibilityAPIDiffProcedure is the fully-qualified name of the ReproducibilityAPI's Diff RPC.
	ReproducibilityAPIDiffProcedure = "/reprod_api.v1.ReproducibilityAPI/Diff"
)

// ReproducibilityAPIClient is a client for the reprod_api.v1.ReproducibilityAPI service.
type ReproducibilityAPIClient interface {
	DerivationReproducibility(context.Context, *connect.Request[reprod_api.DerivationReproducibilityRequest]) (*connect.Response[reprod_api.DerivationReproducibilityResponse], error)
	AttrReproducibilityTimeSeries(context.Context, *connect.Request[reprod_api.AttrReproducibilityTimeSeriesRequest]) (*connect.Response[reprod_api.AttrReproducibilityTimeSeriesResponse], error)
	AttrReproducibilityTimeSeriesGroupedbyChannel(context.Context, *connect.Request[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelRequest]) (*connect.Response[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelResponse], error)
	SuggestAttribute(context.Context, *connect.Request[reprod_api.SuggestsAttributeRequest]) (*connect.Response[reprod_api.SuggestAttributeResponse], error)
	Diff(context.Context, *connect.Request[reprod_api.DiffRequest]) (*connect.Response[reprod_api.DiffResponse], error)
}

// NewReproducibilityAPIClient constructs a client for the reprod_api.v1.ReproducibilityAPI service.
// By default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped
// responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewReproducibilityAPIClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ReproducibilityAPIClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &reproducibilityAPIClient{
		derivationReproducibility: connect.NewClient[reprod_api.DerivationReproducibilityRequest, reprod_api.DerivationReproducibilityResponse](
			httpClient,
			baseURL+ReproducibilityAPIDerivationReproducibilityProcedure,
			opts...,
		),
		attrReproducibilityTimeSeries: connect.NewClient[reprod_api.AttrReproducibilityTimeSeriesRequest, reprod_api.AttrReproducibilityTimeSeriesResponse](
			httpClient,
			baseURL+ReproducibilityAPIAttrReproducibilityTimeSeriesProcedure,
			opts...,
		),
		attrReproducibilityTimeSeriesGroupedbyChannel: connect.NewClient[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelRequest, reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelResponse](
			httpClient,
			baseURL+ReproducibilityAPIAttrReproducibilityTimeSeriesGroupedbyChannelProcedure,
			opts...,
		),
		suggestAttribute: connect.NewClient[reprod_api.SuggestsAttributeRequest, reprod_api.SuggestAttributeResponse](
			httpClient,
			baseURL+ReproducibilityAPISuggestAttributeProcedure,
			opts...,
		),
		diff: connect.NewClient[reprod_api.DiffRequest, reprod_api.DiffResponse](
			httpClient,
			baseURL+ReproducibilityAPIDiffProcedure,
			opts...,
		),
	}
}

// reproducibilityAPIClient implements ReproducibilityAPIClient.
type reproducibilityAPIClient struct {
	derivationReproducibility                     *connect.Client[reprod_api.DerivationReproducibilityRequest, reprod_api.DerivationReproducibilityResponse]
	attrReproducibilityTimeSeries                 *connect.Client[reprod_api.AttrReproducibilityTimeSeriesRequest, reprod_api.AttrReproducibilityTimeSeriesResponse]
	attrReproducibilityTimeSeriesGroupedbyChannel *connect.Client[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelRequest, reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelResponse]
	suggestAttribute                              *connect.Client[reprod_api.SuggestsAttributeRequest, reprod_api.SuggestAttributeResponse]
	diff                                          *connect.Client[reprod_api.DiffRequest, reprod_api.DiffResponse]
}

// DerivationReproducibility calls reprod_api.v1.ReproducibilityAPI.DerivationReproducibility.
func (c *reproducibilityAPIClient) DerivationReproducibility(ctx context.Context, req *connect.Request[reprod_api.DerivationReproducibilityRequest]) (*connect.Response[reprod_api.DerivationReproducibilityResponse], error) {
	return c.derivationReproducibility.CallUnary(ctx, req)
}

// AttrReproducibilityTimeSeries calls
// reprod_api.v1.ReproducibilityAPI.AttrReproducibilityTimeSeries.
func (c *reproducibilityAPIClient) AttrReproducibilityTimeSeries(ctx context.Context, req *connect.Request[reprod_api.AttrReproducibilityTimeSeriesRequest]) (*connect.Response[reprod_api.AttrReproducibilityTimeSeriesResponse], error) {
	return c.attrReproducibilityTimeSeries.CallUnary(ctx, req)
}

// AttrReproducibilityTimeSeriesGroupedbyChannel calls
// reprod_api.v1.ReproducibilityAPI.AttrReproducibilityTimeSeriesGroupedbyChannel.
func (c *reproducibilityAPIClient) AttrReproducibilityTimeSeriesGroupedbyChannel(ctx context.Context, req *connect.Request[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelRequest]) (*connect.Response[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelResponse], error) {
	return c.attrReproducibilityTimeSeriesGroupedbyChannel.CallUnary(ctx, req)
}

// SuggestAttribute calls reprod_api.v1.ReproducibilityAPI.SuggestAttribute.
func (c *reproducibilityAPIClient) SuggestAttribute(ctx context.Context, req *connect.Request[reprod_api.SuggestsAttributeRequest]) (*connect.Response[reprod_api.SuggestAttributeResponse], error) {
	return c.suggestAttribute.CallUnary(ctx, req)
}

// Diff calls reprod_api.v1.ReproducibilityAPI.Diff.
func (c *reproducibilityAPIClient) Diff(ctx context.Context, req *connect.Request[reprod_api.DiffRequest]) (*connect.Response[reprod_api.DiffResponse], error) {
	return c.diff.CallUnary(ctx, req)
}

// ReproducibilityAPIHandler is an implementation of the reprod_api.v1.ReproducibilityAPI service.
type ReproducibilityAPIHandler interface {
	DerivationReproducibility(context.Context, *connect.Request[reprod_api.DerivationReproducibilityRequest]) (*connect.Response[reprod_api.DerivationReproducibilityResponse], error)
	AttrReproducibilityTimeSeries(context.Context, *connect.Request[reprod_api.AttrReproducibilityTimeSeriesRequest]) (*connect.Response[reprod_api.AttrReproducibilityTimeSeriesResponse], error)
	AttrReproducibilityTimeSeriesGroupedbyChannel(context.Context, *connect.Request[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelRequest]) (*connect.Response[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelResponse], error)
	SuggestAttribute(context.Context, *connect.Request[reprod_api.SuggestsAttributeRequest]) (*connect.Response[reprod_api.SuggestAttributeResponse], error)
	Diff(context.Context, *connect.Request[reprod_api.DiffRequest]) (*connect.Response[reprod_api.DiffResponse], error)
}

// NewReproducibilityAPIHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewReproducibilityAPIHandler(svc ReproducibilityAPIHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	reproducibilityAPIDerivationReproducibilityHandler := connect.NewUnaryHandler(
		ReproducibilityAPIDerivationReproducibilityProcedure,
		svc.DerivationReproducibility,
		opts...,
	)
	reproducibilityAPIAttrReproducibilityTimeSeriesHandler := connect.NewUnaryHandler(
		ReproducibilityAPIAttrReproducibilityTimeSeriesProcedure,
		svc.AttrReproducibilityTimeSeries,
		opts...,
	)
	reproducibilityAPIAttrReproducibilityTimeSeriesGroupedbyChannelHandler := connect.NewUnaryHandler(
		ReproducibilityAPIAttrReproducibilityTimeSeriesGroupedbyChannelProcedure,
		svc.AttrReproducibilityTimeSeriesGroupedbyChannel,
		opts...,
	)
	reproducibilityAPISuggestAttributeHandler := connect.NewUnaryHandler(
		ReproducibilityAPISuggestAttributeProcedure,
		svc.SuggestAttribute,
		opts...,
	)
	reproducibilityAPIDiffHandler := connect.NewUnaryHandler(
		ReproducibilityAPIDiffProcedure,
		svc.Diff,
		opts...,
	)
	return "/reprod_api.v1.ReproducibilityAPI/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ReproducibilityAPIDerivationReproducibilityProcedure:
			reproducibilityAPIDerivationReproducibilityHandler.ServeHTTP(w, r)
		case ReproducibilityAPIAttrReproducibilityTimeSeriesProcedure:
			reproducibilityAPIAttrReproducibilityTimeSeriesHandler.ServeHTTP(w, r)
		case ReproducibilityAPIAttrReproducibilityTimeSeriesGroupedbyChannelProcedure:
			reproducibilityAPIAttrReproducibilityTimeSeriesGroupedbyChannelHandler.ServeHTTP(w, r)
		case ReproducibilityAPISuggestAttributeProcedure:
			reproducibilityAPISuggestAttributeHandler.ServeHTTP(w, r)
		case ReproducibilityAPIDiffProcedure:
			reproducibilityAPIDiffHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedReproducibilityAPIHandler returns CodeUnimplemented from all methods.
type UnimplementedReproducibilityAPIHandler struct{}

func (UnimplementedReproducibilityAPIHandler) DerivationReproducibility(context.Context, *connect.Request[reprod_api.DerivationReproducibilityRequest]) (*connect.Response[reprod_api.DerivationReproducibilityResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("reprod_api.v1.ReproducibilityAPI.DerivationReproducibility is not implemented"))
}

func (UnimplementedReproducibilityAPIHandler) AttrReproducibilityTimeSeries(context.Context, *connect.Request[reprod_api.AttrReproducibilityTimeSeriesRequest]) (*connect.Response[reprod_api.AttrReproducibilityTimeSeriesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("reprod_api.v1.ReproducibilityAPI.AttrReproducibilityTimeSeries is not implemented"))
}

func (UnimplementedReproducibilityAPIHandler) AttrReproducibilityTimeSeriesGroupedbyChannel(context.Context, *connect.Request[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelRequest]) (*connect.Response[reprod_api.AttrReproducibilityTimeSeriesGroupedbyChannelResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("reprod_api.v1.ReproducibilityAPI.AttrReproducibilityTimeSeriesGroupedbyChannel is not implemented"))
}

func (UnimplementedReproducibilityAPIHandler) SuggestAttribute(context.Context, *connect.Request[reprod_api.SuggestsAttributeRequest]) (*connect.Response[reprod_api.SuggestAttributeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("reprod_api.v1.ReproducibilityAPI.SuggestAttribute is not implemented"))
}

func (UnimplementedReproducibilityAPIHandler) Diff(context.Context, *connect.Request[reprod_api.DiffRequest]) (*connect.Response[reprod_api.DiffResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("reprod_api.v1.ReproducibilityAPI.Diff is not implemented"))
}
