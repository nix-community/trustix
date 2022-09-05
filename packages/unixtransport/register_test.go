package unixtransport_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nix-community/trustix/packages/unixtransport"
)

func TestBasics(t *testing.T) {
	t.Parallel()

	// This first server will do HTTP.
	var (
		tempdir = t.TempDir()
		socket1 = filepath.Join(tempdir, "1.sock")
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

	transport := &http.Transport{}
	client := &http.Client{Transport: transport}

	// The magic.
	unixtransport.Register(transport)

	// unix should work.
	{
		var (
			rawurl = "unix://" + socket1 + "/foo"
			want   = "1 /foo"
			have   = get(t, client, rawurl)
		)
		if want != have {
			t.Errorf("%s: want %q, have %q", rawurl, want, have)
		}
	}

	// Do another unix request, to kind of verify the connection pool
	// didn't mix things up too badly.
	{
		var (
			rawurl = "unix://" + socket1 + "/bar"
			want   = "1 /bar"
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
