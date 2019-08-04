package cache

import (
	"log"
	"time"

	"github.com/maxArturo/amalgam/internal/link"
	"github.com/maxArturo/amalgam/internal/util"
	"github.com/patrickmn/go-cache"
)

const defaultCacheExpiration = 24 * time.Hour

type Cacher interface {
	Items() map[string]link.RenderedLinker
	Get(hash string) (link.RenderedLinker, bool)
	Set(hash string, link link.RenderedLinker)
}

type Cache struct {
	c *cache.Cache
}

func New() *Cache {
	var cacheExpiration time.Duration
	utilService := util.New()
	cacheDurationVar, err := utilService.GetEnvVarInt("CACHE_EXPIRATION_HRS")
	if err != nil {
		cacheExpiration = defaultCacheExpiration
		log.Println(err)
	} else {
		cacheExpiration = time.Duration(cacheDurationVar) * time.Hour
	}

	return &Cache{
		c: cache.New(cacheExpiration, 10*time.Minute),
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
