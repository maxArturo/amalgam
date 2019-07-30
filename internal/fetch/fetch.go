package fetch

import (
	"log"
	"time"

	"github.com/maxArturo/amalgam/internal/link"
	"github.com/maxArturo/amalgam/internal/worker"
)

type osSleeper interface {
	sleep(duration time.Duration)
}

type SourceFetch struct {
	osSleeper
}

func New() *SourceFetch {
	return &SourceFetch{
		&osSleep{}
	}
}

func (f *SourceFetch) spawnFetcher(count int, pending chan *worker.Source, done chan *worker.Source, updated chan *[]link.RenderedLinker, interval time.Duration) {
	f.sleepSources(done, pending, interval)
	for i := 0; i < count; i++ {
		go f.fetch(i, pending, done, updated)
	}
}

func (f *SourceFetch) sleepSources(done chan *worker.Source, pending chan *worker.Source, duration time.Duration) {
	go func() {
		for s := range done {
			newSource := s
			// go f.sleep(newSource, pending, duration+time.Duration(newSource.errCount))
			go func(){
				f.sleep(duration)
				pending <- newSource
			}()
		}
	}()
}

// Fetch waits on an incoming channel for Sources and fetches them, to update with new links.
// It reports out on the outgoing and content channels for completed Sources.
func (f *SourceFetch) fetch(label int, in chan *worker.Source, out chan *worker.Source, content chan *[]link.RenderedLinker) {
	for src := range in {
		log.Printf("[FETCH] fetcher no %d, fetching for %s", label, src.provider.Name())
		newLinks, err := src.provider.Fetch()
		if err != nil {
			log.Println("[FETCH] Error fetching ", src.provider.Name(), err)
			src.errCount++
		} else {
			src.errCount = 0
			extractedLinks := make([]link.RenderedLinker, len(*newLinks))
			for i := range extractedLinks {
				extractedLinks[i] = link.New((*newLinks)[i])
			}
			content <- &extractedLinks
		}

		out <- src
	}
}

type osSleep struct{}

func (s *osSleep) sleep(d time.Duration) {
	time.Sleep(d)
}
