package storage

import (
	"fmt"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/storage/errors"
	// git "github.com/tweag/trustix/storage/git"
)

// Re-export from subpackage errors for conciseness
var ObjectNotFoundError = errors.ObjectNotFoundError

func FromConfig(conf *config.StorageConfig) (TrustixStorage, error) {
	switch conf.Type {
	// case "git":
	// 	return git.FromConfig(conf.Git)
	case "native":
		return NativeStorageFromConfig(conf.Native)
	}

	return nil, fmt.Errorf("Storage type '%s' is not supported.", conf.Type)

}
