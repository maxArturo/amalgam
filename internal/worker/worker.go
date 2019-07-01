package worker

import (
	"log"
	"time"

	"github.com/maxArturo/amalgam"
)

const defaultErrTimeoutDelay = 5
const defaultFetchInterval = 30
const defaultNumFetchers = 3

// Source holds together a provider and its latest links, along with some stats.
type source struct {
	provider    amalgam.Provider
	errCount    int
	lastUpdated time.Time
}

type SourceJob struct{}

func sleep(s *source, done chan *source) {
	time.Sleep(defaultFetchInterval*time.Second + time.Duration(s.errCount))
	done <- s
}

// Start kicks off workers to fetch new content.
func (s *SourceJob) Start(providers []amalgam.Provider) chan []amalgam.Linker {
	// create our pending/done/new content channels
	pending, done, updated := make(chan *source),
		make(chan *source), make(chan []amalgam.Linker)

	// launch fetchers
	for i := 0; i < defaultNumFetchers; i++ {
		go Fetch(i, pending, done, updated)
	}

	go func() {
		for s := range done {
			newSource := s
			go sleep(newSource, pending)
		}
	}()

	go func() {
		for _, provider := range providers {
			source := &source{
				provider: provider,
			}
			select {
			case pending <- source:
				// TODO review, remove
				time.Sleep(5 * time.Second)
			}
		}
	}()

	return updated
}

// Fetch waits on an incoming channel for Sources and fetches them, to update with new links.
// It reports out on the outgoing and content channels for completed Sources.
func Fetch(label int, in chan *source, out chan *source, content chan []amalgam.Linker) {
	for src := range in {
		log.Printf("[FETCH] fetcher no %d, fetching for %s", label, src.provider.Name())
		newLinks, err := src.provider.Fetch()
		if err != nil {
			log.Println("[FETCH] Error fetching ", src.provider.Name(), err)
			src.errCount++
		} else {
			src.errCount = 0
			log.Printf("in fechtcher no %d, source %s, link count %d", label, src.provider.Name(), len(newLinks))
		}

		content <- newLinks
		out <- src
	}
}
