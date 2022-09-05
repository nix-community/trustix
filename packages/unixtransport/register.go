package unixtransport

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// When UNIX sockets URLs come in they have no meaningful way to
// separate the file system path from the server URL
//
// We can work around this issue by enforcing a file extension
// that makes it possible to separate the two.
var socketSuffixes = []string{".sock", ".socket"}

// Register adds a protocol handler to the provided transport that can serve
// requests to Unix domain sockets via the "unix" or "https+unix" schemes.
// Request URLs should have the following form:
//
//    https+unix:///path/to/socket:/request/path?query=val&...
//
// The registered transport is based on a clone of the provided transport, and
// so uses the same configuration: timeouts, TLS settings, and so on. Connection
// pooling should also work as normal. One caveat: only the DialContext and
// DialTLSContext dialers are respected; the Dial and DialTLS dialers are
// explicitly removed and ignored.
func Register(t *http.Transport) {
	copy := t.Clone()

	copy.Dial = nil // nolint <= deprecated function
	copy.DialTLS = nil

	switch {
	case copy.DialContext == nil && copy.DialTLSContext == nil:
		copy.DialContext = dialContextAdapter(defaultDialContextFunc)

	case copy.DialContext == nil && copy.DialTLSContext != nil:
		copy.DialContext = dialContextAdapter(defaultDialContextFunc)
		copy.DialTLSContext = dialContextAdapter(copy.DialTLSContext)

	case copy.DialContext != nil && copy.DialTLSContext == nil:
		copy.DialContext = dialContextAdapter(copy.DialContext)

	case copy.DialContext != nil && copy.DialTLSContext != nil:
		copy.DialContext = dialContextAdapter(copy.DialContext)
		copy.DialTLSContext = dialContextAdapter(copy.DialTLSContext)
	}

	tt := roundTripAdapter(copy)

	t.RegisterProtocol("unix", tt)
}

// dialContextAdapter decorates the provided DialContext function by trying to base64 decode
// the provided address. If successful, the network is changed to "unix" and the address
// is changed to the decoded value.
func dialContextAdapter(next dialContextFunc) dialContextFunc {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		host, _, err := net.SplitHostPort(address)
		if err != nil {
			host = address
		}

		filepath, err := base64.RawURLEncoding.DecodeString(host)
		if err == nil {
			network, address = "unix", string(filepath)
		}

		return next(ctx, network, address)
	}
}

// roundTripAdapter returns an http.RoundTripper which, when used in combination
// with the dialContextAdapter, supports Unix sockets via any scheme with a
// "+unix" suffix.
func roundTripAdapter(next http.RoundTripper) http.RoundTripper {
	return roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL == nil {
			return nil, fmt.Errorf("unix transport: no request URL")
		}

		if req.URL.Scheme != "unix" {
			return nil, fmt.Errorf("unix transport: invalid scheme '%s'", req.URL.Scheme)
		}

		var socketPath string
		var requestPath string

		for _, suffix := range socketSuffixes {
			idx := strings.Index(req.URL.Path, suffix)
			if idx != -1 {
				sepIdx := idx + len(suffix)
				socketPath = req.URL.Path[0:sepIdx]
				requestPath = req.URL.Path[sepIdx:len(req.URL.Path)]
				break
			}
		}

		if socketPath == "" {
			return nil, fmt.Errorf("unix transport: could not extract socket path: missing .sock/.socket file suffix")
		}

		encodedHost := base64.RawURLEncoding.EncodeToString([]byte(socketPath))

		req = req.Clone(req.Context())
		req.URL.Scheme = "http" // Remove?
		req.URL.Host = encodedHost
		req.URL.Path = requestPath

		return next.RoundTrip(req)
	})
}

type dialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)

var defaultDialContextFunc = (&net.Dialer{}).DialContext

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
