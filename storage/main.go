package storage

import (
	"fmt"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/storage/errors"
)

// Re-export from subpackage errors for conciseness
var ObjectNotFoundError = errors.ObjectNotFoundError

func FromConfig(name string, stateDirectory string, conf *config.StorageConfig) (TrustixStorage, error) {
	switch conf.Type {
	case "git":
		return GitStorageFromConfig(name, stateDirectory, conf.Git)
	case "native":
		return NativeStorageFromConfig(name, stateDirectory, conf.Native)
	}

	return nil, fmt.Errorf("Storage type '%s' is not supported.", conf.Type)

}
