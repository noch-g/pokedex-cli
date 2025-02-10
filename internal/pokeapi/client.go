package pokeapi

import (
	"net/http"
	"os"
	"time"

	"github.com/noch-g/pokedex-cli/internal/logger"
	"github.com/noch-g/pokedex-cli/internal/pokecache"
	"github.com/noch-g/pokedex-cli/internal/settings"
)

type Client struct {
	httpClient http.Client
	cache      pokecache.Cache
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	cache := pokecache.NewCache(cacheInterval)
	client := Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: cache,
	}
	client.LoadCache()
	return client
}

func (c *Client) SaveCache() error {
	return c.cache.SaveToFile(settings.CacheFilePath)
}

func (c *Client) LoadCache() error {
	if _, err := os.Stat(settings.CacheFilePath); !os.IsNotExist(err) {
		logger.Debug("Loading cache", "file", settings.CacheFilePath)
		return c.cache.LoadFromFile(settings.CacheFilePath)
	} else {
		logger.Debug("No cache found", "file", settings.CacheFilePath)
	}
	return nil
}
