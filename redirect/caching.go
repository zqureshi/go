package redirect

import (
	"time"

	"github.com/coocood/freecache"
	"github.com/zqureshi/go/log"
)

const cacheExpiration = 7 * 24 * time.Hour

// CachingClient uses an in-memory cache in front of a Redirector to minimize reads.
type CachingClient struct {
	redirector Redirector
	cache      *freecache.Cache
}

// NewCachingClient constructs a CachingClient with a predefined in-memory cache.
func NewCachingClient(r Redirector) (*CachingClient, error) {
	return &CachingClient{redirector: r, cache: freecache.NewCache(1024)}, nil
}

// Get implements Redirector.
func (c *CachingClient) Get(key string) (*Redirect, error) {
	val, err := c.cache.Get([]byte(key))
	if err == nil {
		return &Redirect{key, string(val)}, nil
	}

	redirect, err := c.redirector.Get(key)
	if err != nil {
		return redirect, err
	}

	log.Logger.Infof("Caching %s -> %s", key, redirect.URL)
	err = c.cache.Set([]byte(key), []byte(redirect.URL), int(cacheExpiration.Seconds()))

	return redirect, err
}
