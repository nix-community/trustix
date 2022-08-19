package cmd

import (
	"os"

	lru "github.com/hashicorp/golang-lru"
	"github.com/nix-community/go-nix/pkg/derivation"
)

// Arbitrary large number of derivations to cache
const cacheSize = 30_000

type CachedDrvParser struct {
	cache *lru.Cache
}

func NewCachedDrvParser() (*CachedDrvParser, error) {
	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}

	return &CachedDrvParser{
		cache: cache,
	}, nil
}

func (c *CachedDrvParser) ReadPath(drvPath string) (*derivation.Derivation, error) {
	cached, ok := c.cache.Get(drvPath)
	if ok {
		return cached.(*derivation.Derivation), nil
	}

	f, err := os.Open(drvPath)
	if err != nil {
		return nil, err
	}

	drv, err := derivation.ReadDerivation(f)
	if err != nil {
		return nil, err
	}

	c.cache.Add(drvPath, drv)

	return drv, nil
}
