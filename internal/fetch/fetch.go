package fetch

import (
	"log"
	"time"

	"github.com/maxArturo/amalgam/internal/link"
	"github.com/maxArturo/amalgam/internal/source"
)

type osSleeper interface {
	sleep(duration time.Duration)
}

type SourceFetch struct {
	osSleeper
}

func New() *SourceFetch {
	return &SourceFetch{
		&osSleep{},
	}
}

func (f *SourceFetch) SpawnFetcher(count int, pending chan *source.Source, done chan *source.Source, updated chan link.RenderedLinker, interval time.Duration) {
	f.sleepSources(done, pending, interval)
	for i := 0; i < count; i++ {
		go f.fetch(i, pending, done, updated)
	}
}

func (f *SourceFetch) sleepSources(done chan *source.Source, pending chan *source.Source, duration time.Duration) {
	go func() {
		for s := range done {
			newSource := s
			go func() {
				f.sleep(duration)
				pending <- newSource
			}()
		}
	}()
}

// Fetch waits on an incoming channel for Sources and fetches them, to update with new links.
// It reports out on the outgoing and link channels for fetched Sources.
func (f *SourceFetch) fetch(label int, in chan *source.Source, out chan *source.Source, outLinks chan link.RenderedLinker) {
	for src := range in {
		log.Printf("[FETCH] fetcher no %d, fetching for %s", label, src.Provider.Name())
		newLinks, err := src.Provider.Fetch()
		if err != nil {
			log.Println("[FETCH] Error fetching ", src.Provider.Name(), err)
			src.ErrCount++
		} else {
			src.ErrCount = 0
			extractedLinks := *newLinks

			for _, fetchedLink := range extractedLinks {
				outLink := link.New(fetchedLink)
				outLinks <- outLink
			}
		}

		out <- src
	}
}

type osSleep struct{}

func (s *osSleep) sleep(d time.Duration) {
	time.Sleep(d)
}
