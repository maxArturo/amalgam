package text_extraction

import (
	"github.com/maxArturo/amalgam/internal/cache"
	"github.com/maxArturo/amalgam/internal/link"
)

type Extract struct{}

func New() *Extract {
	return &Extract{}
}

func (e *Extract) extract(c cache.Cacher, pending chan link.RenderedLinker, done chan cache.Cacher) {
	for incomingLink := range pending {
		_, found := c.Get(incomingLink.Hash())
		if !found {
			incomingLink.FetchLinkText()
			c.Set(incomingLink.Hash(), incomingLink)

			done <- c
		}
	}
}

func (e *Extract) SpawnExtractor(count int, pending chan link.RenderedLinker, done chan cache.Cacher) {
	c := cache.New()
	for i := 0; i < count; i++ {
		go e.extract(c, pending, done)
	}
}
