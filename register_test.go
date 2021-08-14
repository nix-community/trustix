package unixtransport_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/peterbourgon/unixtransport"
)

func TestBasics(t *testing.T) {
	t.Parallel()

	// This first server will do HTTP.
	var (
		tempdir = t.TempDir()
		socket1 = filepath.Join(tempdir, "1")
	)
	{
		ln, err := net.Listen("unix", socket1)
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, 1, r.URL.Path) })
		server := httptest.NewUnstartedServer(handler)
		server.Listener = ln
		server.Start()
		defer server.Close()
	}

	// This second server will speak HTTPS. The httptest.Server can do TLS, but
	// it uses a hard-coded cert with "example.com" as a server name. We'll get
	// that cert in the config's pool after we start the server.
	var (
		socket2         = filepath.Join(tempdir, "2")
		tlsClientConfig = &tls.Config{ServerName: "example.com"}
	)
	{
		ln, err := net.Listen("unix", socket2)
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, 2, r.URL.Path) })
		server := httptest.NewUnstartedServer(handler)
		server.Listener = ln
		server.StartTLS()
		defer server.Close()

		certpool := x509.NewCertPool()
		certpool.AddCert(server.Certificate())
		tlsClientConfig.RootCAs = certpool
	}

	// We could just use a plain http.Client, but for the TLS config required by
	// the second server. Create the transport with the TLS config, and a client
	// that utilizes that transport.
	transport := &http.Transport{TLSClientConfig: tlsClientConfig}
	client := &http.Client{Transport: transport}

	// The magic.
	unixtransport.Register(transport)

	// http+unix should work.
	{
		var (
			rawurl = "http+unix://" + socket1 + ":/foo?a=1"
			want   = "1 /foo"
			have   = get(t, client, rawurl)
		)
		if want != have {
			t.Errorf("%s: want %q, have %q", rawurl, want, have)
		}
	}

	// https+unix should also work.
	{
		var (
			rawurl = "https+unix://" + socket2 + ":/bar#fragment"
			want   = "2 /bar"
			have   = get(t, client, rawurl)
		)
		if want != have {
			t.Errorf("%s: want %q, have %q", rawurl, want, have)
		}
	}

	// Do another http+unix request, to kind of verify the connection pool
	// didn't mix things up too badly.
	{
		var (
			rawurl = "http+unix://" + socket1 + ":/baz:baz:baz"
			want   = "1 /baz:baz:baz"
			have   = get(t, client, rawurl)
		)
		if want != have {
			t.Errorf("%s: want %q, have %q", rawurl, want, have)
		}
	}
}

func get(t *testing.T, client *http.Client, rawurl string) string {
	t.Helper()

	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		t.Errorf("GET %s: %v", rawurl, err)
		return ""
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("GET %s: %v", rawurl, err)
		return ""
	}

	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("GET %s: %v", rawurl, err)
		return ""
	}

	return strings.TrimSpace(string(buf))
}
