package text_extraction

import (
	"sync"

	"github.com/maxArturo/amalgam/internal/link"
)

const MAX_LINKS = 500

type linkCache struct {
	sync.RWMutex
	links      map[string]link.RenderedLinker
	oldestLink string
}

type Extract struct{}

func (e *Extract) extract(cache *linkCache, pending chan link.RenderedLinker, done chan link.RenderedLinker) {
	for incomingLink := range pending {
		cache.RLock()
		cachedLink := cache.links[incomingLink.Hash()]
		cache.RUnlock()

		if cachedLink == nil {

		} else {
			done <- incomingLink
		}
	}
}

func (e *Extract) SpawnExtractor(workerCount int, pending chan link.RenderedLinker, done chan link.RenderedLinker) {
	cache := &linkCache{
		links: make(map[string]link.RenderedLinker),
	}

	for i := 0; i < workerCount; i++ {
		go e.extract(cache, pending, done)
	}
}
