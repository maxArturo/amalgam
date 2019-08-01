package text_extraction

import (
	"github.com/maxArturo/amalgam/internal/cache"
	"github.com/maxArturo/amalgam/internal/link"
)

type Extract struct {
	c cache.Cacher
}

func New(c cache.Cacher) *Extract {
	return &Extract{
		c,
	}
}

func (e *Extract) extract(pending chan link.RenderedLinker, done chan cache.Cacher) {
	for incomingLink := range pending {
		_, found := e.c.Get(incomingLink.Hash())
		if !found {
			incomingLink.FetchLinkText()
			e.c.Set(incomingLink.Hash(), incomingLink)

			done <- e.c
		}
	}
}

func (e *Extract) SpawnExtractor(count int, pending chan link.RenderedLinker, done chan cache.Cacher) {
	for i := 0; i < count; i++ {
		go e.extract(pending, done)
	}
}
