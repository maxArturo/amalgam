package worker

import (
	"log"
	"time"

	"github.com/maxArturo/amalgam"
	"github.com/maxArturo/amalgam/internal/cache"
	"github.com/maxArturo/amalgam/internal/fetch"
	"github.com/maxArturo/amalgam/internal/link"
	"github.com/maxArturo/amalgam/internal/source"
	extraction "github.com/maxArturo/amalgam/internal/text_extraction"
	"github.com/maxArturo/amalgam/internal/util"
)

const defaultFetchInterval = 30
const defaultNumFetchers = 3
const defaultNumExtractors = 5

type fetcher interface {
	SpawnFetcher(count int, pending chan *source.Source, done chan *source.Source, updated chan link.RenderedLinker, interval time.Duration)
}

type extractor interface {
	SpawnExtractor(count int, pending chan link.RenderedLinker, done chan cache.Cacher)
}

// FetchJob contains the config needed for fetching provider links.
type FetchJob struct {
	fetchInterval int
	numFetchers   int
	numExtractors int
	fetcher
	extractor
}

// New creates a configured FetchJob ready to use.
func New() *FetchJob {
	utilService := util.New()
	fetchInterval, err := utilService.GetEnvVarInt("FETCH_INTERVAL")
	if err != nil {
		fetchInterval = defaultFetchInterval
		log.Println(err)
	}

	numFetchers, err := utilService.GetEnvVarInt("NUM_FETCHERS")
	if err != nil {
		numFetchers = defaultNumFetchers
		log.Println(err)
	}

	numExtractors, err := utilService.GetEnvVarInt("NUM_EXTRACTORS")
	if err != nil {
		numExtractors = defaultNumExtractors
		log.Println(err)
	}

	return &FetchJob{
		fetchInterval: fetchInterval,
		numFetchers:   numFetchers,
		numExtractors: numExtractors,
		fetcher:       fetch.New(),
		extractor:     extraction.New(),
	}
}

// Start kicks off workers to fetch new content.
func (f *FetchJob) Start(providers *[]amalgam.Provider) chan cache.Cacher {
	// create our pending/done/new content channels
	pendingSources, doneSources := make(chan *source.Source), make(chan *source.Source)
	fetchedLinks := make(chan link.RenderedLinker, f.numExtractors)
	updatedCache := make(chan cache.Cacher)

	f.fetcher.SpawnFetcher(f.numFetchers, pendingSources, doneSources, fetchedLinks, time.Duration(f.fetchInterval)*time.Second)
	f.extractor.SpawnExtractor(f.numExtractors, fetchedLinks, updatedCache)

	go func() {
		for _, provider := range *providers {
			src := &source.Source{
				Provider: provider,
			}
			pendingSources <- src
		}
	}()

	return updatedCache
}
