package main

// Request is a convenience type that we'll extend in the next iteration
type Request struct {
	thing Thing
	extra Extra
}

// Response is a convenience type that we'll extend in the next iteration
type Response struct {
	status string
}
