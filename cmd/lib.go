package cmd

import (
	"fmt"
	"os"
)

// Error - Write error and exit with status
func ErrorStatus(err error, status int) {
	fmt.Fprintf(os.Stderr, "Encountered error: %s, exiting", err.Error())
	os.Exit(status)
}

// Error - Write error and exit with status 1
func Error(err error) {
	ErrorStatus(err, 1)
}
