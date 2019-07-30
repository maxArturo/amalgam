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
type Source struct {
	provider    amalgam.Provider
	errCount    int
	lastUpdated time.Time
}

type fetcher interface {
	spawnFetcher(count int, pending chan *Source, done chan *Source, updated chan *[]link.RenderedLinker, interval time.Duration)
}

// FetchJob contains the config needed for fetching provider links.
type FetchJob struct {
	fetchInterval int
	numFetchers   int
	fetcher
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
		fetcher:       &SourceFetch{},
	}
}

// Start kicks off workers to fetch new content.
func (f *FetchJob) Start(providers *[]amalgam.Provider) chan *[]link.RenderedLinker {
	// create our pending/done/new content channels
	pending, done := make(chan *Source),
		make(chan *Source)
	updated := make(chan *[]link.RenderedLinker)

	f.fetcher.spawnFetcher(f.numFetchers, pending, done, updated, time.Duration(f.fetchInterval)*time.Second)

	go func() {
		for _, provider := range *providers {
			Source := &Source{
				provider: provider,
			}
			pending <- Source
		}
	}()

	return updated
}
