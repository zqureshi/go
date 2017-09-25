package redirect

import (
	"github.com/coocood/freecache"
	"log"
)

const cacheExpiration = 86400 * 7 // 7 days

type CachingClient struct {
	redirect *Client
	cache    *freecache.Cache
}

func NewCaching(apiKey string, baseId string) (*CachingClient, error) {
	client, err := New(apiKey, baseId)
	if err != nil {
		return nil, err
	}

	return &CachingClient{redirect: client, cache: freecache.NewCache(1024)}, nil
}

func (client *CachingClient) Get(key string) (*Redirect, error) {
	val, err := client.cache.Get([]byte(key))
	if err == nil {
		return &Redirect{key, string(val)}, nil
	}

	redirect, err := client.redirect.Get(key)
	if err != nil {
		return redirect, err
	}

	log.Println("Caching", key, "->", redirect.URL)
	err = client.cache.Set([]byte(key), []byte(redirect.URL), cacheExpiration)

	return redirect, err
}

func (client *CachingClient) GetDefault() (*Redirect, error) {
	return client.Get(defaultRedirectKey)
}
