// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: MIT

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/api.proto

package apiconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	api "github.com/nix-community/trustix/packages/trustix-proto/api"
	schema "github.com/nix-community/trustix/packages/trustix-proto/schema"
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
	// NodeAPIName is the fully-qualified name of the NodeAPI service.
	NodeAPIName = "trustix_api.v1.NodeAPI"
	// LogAPIName is the fully-qualified name of the LogAPI service.
	LogAPIName = "trustix_api.v1.LogAPI"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// NodeAPILogsProcedure is the fully-qualified name of the NodeAPI's Logs RPC.
	NodeAPILogsProcedure = "/trustix_api.v1.NodeAPI/Logs"
	// NodeAPIGetValueProcedure is the fully-qualified name of the NodeAPI's GetValue RPC.
	NodeAPIGetValueProcedure = "/trustix_api.v1.NodeAPI/GetValue"
	// LogAPIGetHeadProcedure is the fully-qualified name of the LogAPI's GetHead RPC.
	LogAPIGetHeadProcedure = "/trustix_api.v1.LogAPI/GetHead"
	// LogAPIGetLogConsistencyProofProcedure is the fully-qualified name of the LogAPI's
	// GetLogConsistencyProof RPC.
	LogAPIGetLogConsistencyProofProcedure = "/trustix_api.v1.LogAPI/GetLogConsistencyProof"
	// LogAPIGetLogAuditProofProcedure is the fully-qualified name of the LogAPI's GetLogAuditProof RPC.
	LogAPIGetLogAuditProofProcedure = "/trustix_api.v1.LogAPI/GetLogAuditProof"
	// LogAPIGetLogEntriesProcedure is the fully-qualified name of the LogAPI's GetLogEntries RPC.
	LogAPIGetLogEntriesProcedure = "/trustix_api.v1.LogAPI/GetLogEntries"
	// LogAPIGetMapValueProcedure is the fully-qualified name of the LogAPI's GetMapValue RPC.
	LogAPIGetMapValueProcedure = "/trustix_api.v1.LogAPI/GetMapValue"
	// LogAPIGetMHLogConsistencyProofProcedure is the fully-qualified name of the LogAPI's
	// GetMHLogConsistencyProof RPC.
	LogAPIGetMHLogConsistencyProofProcedure = "/trustix_api.v1.LogAPI/GetMHLogConsistencyProof"
	// LogAPIGetMHLogAuditProofProcedure is the fully-qualified name of the LogAPI's GetMHLogAuditProof
	// RPC.
	LogAPIGetMHLogAuditProofProcedure = "/trustix_api.v1.LogAPI/GetMHLogAuditProof"
	// LogAPIGetMHLogEntriesProcedure is the fully-qualified name of the LogAPI's GetMHLogEntries RPC.
	LogAPIGetMHLogEntriesProcedure = "/trustix_api.v1.LogAPI/GetMHLogEntries"
)

// NodeAPIClient is a client for the trustix_api.v1.NodeAPI service.
type NodeAPIClient interface {
	// Get a list of all logs published by this node
	Logs(context.Context, *connect.Request[api.LogsRequest]) (*connect.Response[api.LogsResponse], error)
	// Get values by their content-address
	GetValue(context.Context, *connect.Request[api.ValueRequest]) (*connect.Response[api.ValueResponse], error)
}

// NewNodeAPIClient constructs a client for the trustix_api.v1.NodeAPI service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewNodeAPIClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) NodeAPIClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &nodeAPIClient{
		logs: connect.NewClient[api.LogsRequest, api.LogsResponse](
			httpClient,
			baseURL+NodeAPILogsProcedure,
			opts...,
		),
		getValue: connect.NewClient[api.ValueRequest, api.ValueResponse](
			httpClient,
			baseURL+NodeAPIGetValueProcedure,
			opts...,
		),
	}
}

// nodeAPIClient implements NodeAPIClient.
type nodeAPIClient struct {
	logs     *connect.Client[api.LogsRequest, api.LogsResponse]
	getValue *connect.Client[api.ValueRequest, api.ValueResponse]
}

// Logs calls trustix_api.v1.NodeAPI.Logs.
func (c *nodeAPIClient) Logs(ctx context.Context, req *connect.Request[api.LogsRequest]) (*connect.Response[api.LogsResponse], error) {
	return c.logs.CallUnary(ctx, req)
}

// GetValue calls trustix_api.v1.NodeAPI.GetValue.
func (c *nodeAPIClient) GetValue(ctx context.Context, req *connect.Request[api.ValueRequest]) (*connect.Response[api.ValueResponse], error) {
	return c.getValue.CallUnary(ctx, req)
}

// NodeAPIHandler is an implementation of the trustix_api.v1.NodeAPI service.
type NodeAPIHandler interface {
	// Get a list of all logs published by this node
	Logs(context.Context, *connect.Request[api.LogsRequest]) (*connect.Response[api.LogsResponse], error)
	// Get values by their content-address
	GetValue(context.Context, *connect.Request[api.ValueRequest]) (*connect.Response[api.ValueResponse], error)
}

// NewNodeAPIHandler builds an HTTP handler from the service implementation. It returns the path on
// which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewNodeAPIHandler(svc NodeAPIHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	nodeAPILogsHandler := connect.NewUnaryHandler(
		NodeAPILogsProcedure,
		svc.Logs,
		opts...,
	)
	nodeAPIGetValueHandler := connect.NewUnaryHandler(
		NodeAPIGetValueProcedure,
		svc.GetValue,
		opts...,
	)
	return "/trustix_api.v1.NodeAPI/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case NodeAPILogsProcedure:
			nodeAPILogsHandler.ServeHTTP(w, r)
		case NodeAPIGetValueProcedure:
			nodeAPIGetValueHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedNodeAPIHandler returns CodeUnimplemented from all methods.
type UnimplementedNodeAPIHandler struct{}

func (UnimplementedNodeAPIHandler) Logs(context.Context, *connect.Request[api.LogsRequest]) (*connect.Response[api.LogsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.NodeAPI.Logs is not implemented"))
}

func (UnimplementedNodeAPIHandler) GetValue(context.Context, *connect.Request[api.ValueRequest]) (*connect.Response[api.ValueResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.NodeAPI.GetValue is not implemented"))
}

// LogAPIClient is a client for the trustix_api.v1.LogAPI service.
type LogAPIClient interface {
	// Get signed head
	GetHead(context.Context, *connect.Request[api.LogHeadRequest]) (*connect.Response[schema.LogHead], error)
	GetLogConsistencyProof(context.Context, *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error)
	GetLogAuditProof(context.Context, *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error)
	GetLogEntries(context.Context, *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error)
	GetMapValue(context.Context, *connect.Request[api.GetMapValueRequest]) (*connect.Response[api.MapValueResponse], error)
	GetMHLogConsistencyProof(context.Context, *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error)
	GetMHLogAuditProof(context.Context, *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error)
	GetMHLogEntries(context.Context, *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error)
}

// NewLogAPIClient constructs a client for the trustix_api.v1.LogAPI service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewLogAPIClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) LogAPIClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &logAPIClient{
		getHead: connect.NewClient[api.LogHeadRequest, schema.LogHead](
			httpClient,
			baseURL+LogAPIGetHeadProcedure,
			opts...,
		),
		getLogConsistencyProof: connect.NewClient[api.GetLogConsistencyProofRequest, api.ProofResponse](
			httpClient,
			baseURL+LogAPIGetLogConsistencyProofProcedure,
			opts...,
		),
		getLogAuditProof: connect.NewClient[api.GetLogAuditProofRequest, api.ProofResponse](
			httpClient,
			baseURL+LogAPIGetLogAuditProofProcedure,
			opts...,
		),
		getLogEntries: connect.NewClient[api.GetLogEntriesRequest, api.LogEntriesResponse](
			httpClient,
			baseURL+LogAPIGetLogEntriesProcedure,
			opts...,
		),
		getMapValue: connect.NewClient[api.GetMapValueRequest, api.MapValueResponse](
			httpClient,
			baseURL+LogAPIGetMapValueProcedure,
			opts...,
		),
		getMHLogConsistencyProof: connect.NewClient[api.GetLogConsistencyProofRequest, api.ProofResponse](
			httpClient,
			baseURL+LogAPIGetMHLogConsistencyProofProcedure,
			opts...,
		),
		getMHLogAuditProof: connect.NewClient[api.GetLogAuditProofRequest, api.ProofResponse](
			httpClient,
			baseURL+LogAPIGetMHLogAuditProofProcedure,
			opts...,
		),
		getMHLogEntries: connect.NewClient[api.GetLogEntriesRequest, api.LogEntriesResponse](
			httpClient,
			baseURL+LogAPIGetMHLogEntriesProcedure,
			opts...,
		),
	}
}

// logAPIClient implements LogAPIClient.
type logAPIClient struct {
	getHead                  *connect.Client[api.LogHeadRequest, schema.LogHead]
	getLogConsistencyProof   *connect.Client[api.GetLogConsistencyProofRequest, api.ProofResponse]
	getLogAuditProof         *connect.Client[api.GetLogAuditProofRequest, api.ProofResponse]
	getLogEntries            *connect.Client[api.GetLogEntriesRequest, api.LogEntriesResponse]
	getMapValue              *connect.Client[api.GetMapValueRequest, api.MapValueResponse]
	getMHLogConsistencyProof *connect.Client[api.GetLogConsistencyProofRequest, api.ProofResponse]
	getMHLogAuditProof       *connect.Client[api.GetLogAuditProofRequest, api.ProofResponse]
	getMHLogEntries          *connect.Client[api.GetLogEntriesRequest, api.LogEntriesResponse]
}

// GetHead calls trustix_api.v1.LogAPI.GetHead.
func (c *logAPIClient) GetHead(ctx context.Context, req *connect.Request[api.LogHeadRequest]) (*connect.Response[schema.LogHead], error) {
	return c.getHead.CallUnary(ctx, req)
}

// GetLogConsistencyProof calls trustix_api.v1.LogAPI.GetLogConsistencyProof.
func (c *logAPIClient) GetLogConsistencyProof(ctx context.Context, req *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error) {
	return c.getLogConsistencyProof.CallUnary(ctx, req)
}

// GetLogAuditProof calls trustix_api.v1.LogAPI.GetLogAuditProof.
func (c *logAPIClient) GetLogAuditProof(ctx context.Context, req *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error) {
	return c.getLogAuditProof.CallUnary(ctx, req)
}

// GetLogEntries calls trustix_api.v1.LogAPI.GetLogEntries.
func (c *logAPIClient) GetLogEntries(ctx context.Context, req *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error) {
	return c.getLogEntries.CallUnary(ctx, req)
}

// GetMapValue calls trustix_api.v1.LogAPI.GetMapValue.
func (c *logAPIClient) GetMapValue(ctx context.Context, req *connect.Request[api.GetMapValueRequest]) (*connect.Response[api.MapValueResponse], error) {
	return c.getMapValue.CallUnary(ctx, req)
}

// GetMHLogConsistencyProof calls trustix_api.v1.LogAPI.GetMHLogConsistencyProof.
func (c *logAPIClient) GetMHLogConsistencyProof(ctx context.Context, req *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error) {
	return c.getMHLogConsistencyProof.CallUnary(ctx, req)
}

// GetMHLogAuditProof calls trustix_api.v1.LogAPI.GetMHLogAuditProof.
func (c *logAPIClient) GetMHLogAuditProof(ctx context.Context, req *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error) {
	return c.getMHLogAuditProof.CallUnary(ctx, req)
}

// GetMHLogEntries calls trustix_api.v1.LogAPI.GetMHLogEntries.
func (c *logAPIClient) GetMHLogEntries(ctx context.Context, req *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error) {
	return c.getMHLogEntries.CallUnary(ctx, req)
}

// LogAPIHandler is an implementation of the trustix_api.v1.LogAPI service.
type LogAPIHandler interface {
	// Get signed head
	GetHead(context.Context, *connect.Request[api.LogHeadRequest]) (*connect.Response[schema.LogHead], error)
	GetLogConsistencyProof(context.Context, *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error)
	GetLogAuditProof(context.Context, *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error)
	GetLogEntries(context.Context, *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error)
	GetMapValue(context.Context, *connect.Request[api.GetMapValueRequest]) (*connect.Response[api.MapValueResponse], error)
	GetMHLogConsistencyProof(context.Context, *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error)
	GetMHLogAuditProof(context.Context, *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error)
	GetMHLogEntries(context.Context, *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error)
}

// NewLogAPIHandler builds an HTTP handler from the service implementation. It returns the path on
// which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewLogAPIHandler(svc LogAPIHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	logAPIGetHeadHandler := connect.NewUnaryHandler(
		LogAPIGetHeadProcedure,
		svc.GetHead,
		opts...,
	)
	logAPIGetLogConsistencyProofHandler := connect.NewUnaryHandler(
		LogAPIGetLogConsistencyProofProcedure,
		svc.GetLogConsistencyProof,
		opts...,
	)
	logAPIGetLogAuditProofHandler := connect.NewUnaryHandler(
		LogAPIGetLogAuditProofProcedure,
		svc.GetLogAuditProof,
		opts...,
	)
	logAPIGetLogEntriesHandler := connect.NewUnaryHandler(
		LogAPIGetLogEntriesProcedure,
		svc.GetLogEntries,
		opts...,
	)
	logAPIGetMapValueHandler := connect.NewUnaryHandler(
		LogAPIGetMapValueProcedure,
		svc.GetMapValue,
		opts...,
	)
	logAPIGetMHLogConsistencyProofHandler := connect.NewUnaryHandler(
		LogAPIGetMHLogConsistencyProofProcedure,
		svc.GetMHLogConsistencyProof,
		opts...,
	)
	logAPIGetMHLogAuditProofHandler := connect.NewUnaryHandler(
		LogAPIGetMHLogAuditProofProcedure,
		svc.GetMHLogAuditProof,
		opts...,
	)
	logAPIGetMHLogEntriesHandler := connect.NewUnaryHandler(
		LogAPIGetMHLogEntriesProcedure,
		svc.GetMHLogEntries,
		opts...,
	)
	return "/trustix_api.v1.LogAPI/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case LogAPIGetHeadProcedure:
			logAPIGetHeadHandler.ServeHTTP(w, r)
		case LogAPIGetLogConsistencyProofProcedure:
			logAPIGetLogConsistencyProofHandler.ServeHTTP(w, r)
		case LogAPIGetLogAuditProofProcedure:
			logAPIGetLogAuditProofHandler.ServeHTTP(w, r)
		case LogAPIGetLogEntriesProcedure:
			logAPIGetLogEntriesHandler.ServeHTTP(w, r)
		case LogAPIGetMapValueProcedure:
			logAPIGetMapValueHandler.ServeHTTP(w, r)
		case LogAPIGetMHLogConsistencyProofProcedure:
			logAPIGetMHLogConsistencyProofHandler.ServeHTTP(w, r)
		case LogAPIGetMHLogAuditProofProcedure:
			logAPIGetMHLogAuditProofHandler.ServeHTTP(w, r)
		case LogAPIGetMHLogEntriesProcedure:
			logAPIGetMHLogEntriesHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedLogAPIHandler returns CodeUnimplemented from all methods.
type UnimplementedLogAPIHandler struct{}

func (UnimplementedLogAPIHandler) GetHead(context.Context, *connect.Request[api.LogHeadRequest]) (*connect.Response[schema.LogHead], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.LogAPI.GetHead is not implemented"))
}

func (UnimplementedLogAPIHandler) GetLogConsistencyProof(context.Context, *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.LogAPI.GetLogConsistencyProof is not implemented"))
}

func (UnimplementedLogAPIHandler) GetLogAuditProof(context.Context, *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.LogAPI.GetLogAuditProof is not implemented"))
}

func (UnimplementedLogAPIHandler) GetLogEntries(context.Context, *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.LogAPI.GetLogEntries is not implemented"))
}

func (UnimplementedLogAPIHandler) GetMapValue(context.Context, *connect.Request[api.GetMapValueRequest]) (*connect.Response[api.MapValueResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.LogAPI.GetMapValue is not implemented"))
}

func (UnimplementedLogAPIHandler) GetMHLogConsistencyProof(context.Context, *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.LogAPI.GetMHLogConsistencyProof is not implemented"))
}

func (UnimplementedLogAPIHandler) GetMHLogAuditProof(context.Context, *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.LogAPI.GetMHLogAuditProof is not implemented"))
}

func (UnimplementedLogAPIHandler) GetMHLogEntries(context.Context, *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("trustix_api.v1.LogAPI.GetMHLogEntries is not implemented"))
}