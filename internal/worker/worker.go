package worker

import (
	"log"
	"time"

	"github.com/maxArturo/amalgam"
	"github.com/maxArturo/amalgam/internal/util"
)

const defaultFetchInterval = 30
const defaultNumFetchers = 3

// Source holds together a provider and its latest links, along with some stats.
type source struct {
	provider    amalgam.Provider
	errCount    int
	lastUpdated time.Time
}

type fetcher interface {
	spawnFetchers(count int, pending chan *source, done chan *source, updated chan []amalgam.Linker)
}

type sleeper interface {
	sleepSources(done chan *source, pending chan *source, duration time.Duration)
}

// FetchJob contains the config needed for fetching provider links.
type FetchJob struct {
	fetchInterval int
	numFetchers   int
	fetcher
	sleeper
}

// New creates a configured FetchJob ready to use.
func New() *FetchJob {
	utilService := util.New()
	fetchInterval, err := utilService.GetEnvVarInt("FETCH_INTERVAL")
	if err != nil {
		fetchInterval = defaultFetchInterval
	}

	numFetchers, err := utilService.GetEnvVarInt("FETCH_INTERVAL")
	if err != nil {
		numFetchers = defaultNumFetchers
		log.Println(err)
	}

	return &FetchJob{
		fetchInterval: fetchInterval,
		numFetchers:   numFetchers,
		fetcher:       &fetchProvider{},
		sleeper:       &sleepProvider{},
	}
}

// Start kicks off workers to fetch new content.
func (f *FetchJob) Start(providers []amalgam.Provider) chan []amalgam.Linker {
	// create our pending/done/new content channels
	pending, done, updated := make(chan *source),
		make(chan *source), make(chan []amalgam.Linker)

	f.fetcher.spawnFetchers(f.numFetchers, pending, done, updated)

	f.sleeper.sleepSources(done, pending, time.Duration(f.fetchInterval)*time.Second)

	go func() {
		for _, provider := range providers {
			source := &source{
				provider: provider,
			}
			pending <- source
		}
	}()

	return updated
}

type sleepProvider struct{}

func (provider *sleepProvider) sleepSources(done chan *source, pending chan *source, duration time.Duration) {
	go func() {
		for s := range done {
			newSource := s
			go provider.sleep(newSource, pending, duration+time.Duration(newSource.errCount))
		}
	}()
}

func (provider *sleepProvider) sleep(s *source, done chan *source, sleepDuration time.Duration) {
	time.Sleep(sleepDuration)
	done <- s
}

type fetchProvider struct{}

func (f *fetchProvider) spawnFetchers(count int, pending chan *source, done chan *source, updated chan []amalgam.Linker) {
	for i := 0; i < count; i++ {
		go f.fetch(i, pending, done, updated)
	}
}

// Fetch waits on an incoming channel for Sources and fetches them, to update with new links.
// It reports out on the outgoing and content channels for completed Sources.
func (f *fetchProvider) fetch(label int, in chan *source, out chan *source, content chan []amalgam.Linker) {
	for src := range in {
		log.Printf("[FETCH] fetcher no %d, fetching for %s", label, src.provider.Name())
		newLinks, err := src.provider.Fetch()
		if err != nil {
			log.Println("[FETCH] Error fetching ", src.provider.Name(), err)
			src.errCount++
		} else {
			src.errCount = 0
		}

		content <- newLinks
		out <- src
	}
}
