package cache

import (
	"time"

	"github.com/maxArturo/amalgam/internal/link"
	"github.com/patrickmn/go-cache"
)

const DEFAULT_CACHE_EXPIRATION_HRS = 48 * time.Hour

type Cacher interface {
	Items() map[string]link.RenderedLinker
	Get(hash string) (link.RenderedLinker, bool)
	Set(hash string, link link.RenderedLinker)
}

type Cache struct {
	c *cache.Cache
}

func New() *Cache {
	return &Cache{
		c: cache.New(DEFAULT_CACHE_EXPIRATION_HRS, 10*time.Minute),
	}
}

func (c *Cache) Items() map[string]link.RenderedLinker {
	items := c.c.Items()
	res := make(map[string]link.RenderedLinker)

	for k, v := range items {
		res[k] = v.Object.(link.RenderedLinker)
	}
	return res
}

func (c *Cache) Get(hash string) (link.RenderedLinker, bool) {
	val, found := c.c.Get(hash)
	if val != nil {
		return val.(link.RenderedLinker), found
	}
	return nil, found
}

func (c *Cache) Set(hash string, link link.RenderedLinker) {
	c.c.SetDefault(hash, link)
}
