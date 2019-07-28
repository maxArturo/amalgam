package worker

import (
	"log"

	"github.com/maxArturo/amalgam/internal/link"
)

type fetchProvider struct{}

func (f *fetchProvider) spawnFetchers(count int, pending chan *source, done chan *source, updated chan *[]link.RenderedLinker) {
	for i := 0; i < count; i++ {
		go f.fetch(i, pending, done, updated)
	}
}

// Fetch waits on an incoming channel for Sources and fetches them, to update with new links.
// It reports out on the outgoing and content channels for completed Sources.
func (f *fetchProvider) fetch(label int, in chan *source, out chan *source, content chan *[]link.RenderedLinker) {
	for src := range in {
		log.Printf("[FETCH] fetcher no %d, fetching for %s", label, src.provider.Name())
		newLinks, err := src.provider.Fetch()
		if err != nil {
			log.Println("[FETCH] Error fetching ", src.provider.Name(), err)
			src.errCount++
		} else {
			src.errCount = 0
			extractedLinks := make([]link.RenderedLinker, len(*newLinks))
			for _, l := range *newLinks {
				extractedLinks = append(extractedLinks, link.RenderedLinker(link.New(l)))
			}
			content <- &extractedLinks
		}

		out <- src
	}
}
