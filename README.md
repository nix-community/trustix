# unixtransport [![Go Reference](https://pkg.go.dev/badge/github.com/peterbourgon/unixtransport.svg)](https://pkg.go.dev/github.com/peterbourgon/unixtransport) ![tests](https://github.com/peterbourgon/unixtransport/actions/workflows/test.yaml/badge.svg?branch=main)

This package adds support for Unix domain sockets in Go HTTP clients.

```go
t := &http.Transport{...}

unixtransport.Register(t)

client := &http.Client{Transport: t}
```

Now you can make requests with URLs like this:

```go
resp, err := client.Get("https+unix://unix:/path/to/socket:/request/path?a=b")
```

Use scheme `http+unix` or `https+unix`. The host has to be just `unix`.

Inspiration taken from, and thanks given to, both
[tv42/httpunix](https://github.com/tv42/httpunix) and
[agorman/httpunix](https://github.com/agorman/httpunix).
