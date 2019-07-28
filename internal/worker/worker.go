package worker

import (
	"log"
	"time"

	"github.com/maxArturo/amalgam"
	"github.com/maxArturo/amalgam/internal/link"
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
	spawnFetchers(count int, pending chan *source, done chan *source, updated chan *[]link.RenderedLinker)
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
		sleeper:       newSleepProvider(),
	}
}

// Start kicks off workers to fetch new content.
func (f *FetchJob) Start(providers *[]amalgam.Provider) chan *[]link.RenderedLinker {
	// create our pending/done/new content channels
	pending, done := make(chan *source),
		make(chan *source)
	updated := make(chan *[]link.RenderedLinker)

	f.fetcher.spawnFetchers(f.numFetchers, pending, done, updated)

	f.sleeper.sleepSources(done, pending, time.Duration(f.fetchInterval)*time.Second)

	go func() {
		for _, provider := range *providers {
			source := &source{
				provider: provider,
			}
			pending <- source
		}
	}()

	return updated
}
