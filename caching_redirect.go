package main

import (
	"github.com/coocood/freecache"
	"log"
)

const cacheExpiration = 86400 * 7 // 7 days

type CachingRedirectClient struct {
	redirectClient *RedirectClient
	cache          *freecache.Cache
}

func NewCachingRedirectClient(apiKey string, baseId string) (*CachingRedirectClient, error) {
	client, err := NewRedirectClient(apiKey, baseId)
	if err != nil {
		return nil, err
	}

	return &CachingRedirectClient{redirectClient: client, cache: freecache.NewCache(1024)}, nil
}

func (client *CachingRedirectClient) Get(key string) (*Redirect, error) {
	val, err := client.cache.Get([]byte(key))
	if err == nil {
		return &Redirect{key, string(val)}, nil
	}

	redirect, err := client.redirectClient.Get(key)
	if err != nil {
		return redirect, err
	}

	log.Println("Caching", key, "->", redirect.URL)
	err = client.cache.Set([]byte(key), []byte(redirect.URL), cacheExpiration)

	return redirect, err
}

func (client *CachingRedirectClient) GetDefault() (*Redirect, error) {
	return client.Get(defaultRedirectKey)
}
