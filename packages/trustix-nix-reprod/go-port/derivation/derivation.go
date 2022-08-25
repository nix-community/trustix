package derivation

import (
	"os"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/nix-community/go-nix/pkg/derivation"
)

type CachedDrvParser struct {
	cache *lru.Cache
	mux   *sync.RWMutex
}

func NewCachedDrvParser(cacheSize int) (*CachedDrvParser, error) {
	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}

	return &CachedDrvParser{
		cache: cache,
		mux:   &sync.RWMutex{},
	}, nil
}

func (c *CachedDrvParser) ReadPath(drvPath string) (*derivation.Derivation, error) {

	{
		c.mux.RLock()

		cached, ok := c.cache.Get(drvPath)
		if ok {
			c.mux.RUnlock()
			return cached.(*derivation.Derivation), nil
		}

		c.mux.RUnlock()
	}

	c.mux.Lock()
	defer c.mux.Unlock()

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
