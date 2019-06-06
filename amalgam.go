package amalgam

import (
	"log"
	"net/http"
	"time"

	"github.com/maxArturo/amalgam/pkg/util"
	provider "github.com/maxArturo/amalgam/internal/provider"
)

// Engine is the main way to run an amalgam server.
type Engine struct {
	Providers []Provider
}

// Run starts the amalgam server.
func (e *Engine) Run(addr ...string) (err error) {
	defer func() { log.Println(err.Error()) }()
	port := util.ResolveAddress(addr)
	err = http.ListenAndServe(port, nil)
	return
}

// Default creates a normal engine with default settings:
// * Hacker News and Reddit enabled
func Default() *Engine {
	return &Engine{
		Providers {
			provider.HackerNews(), 
			provider.Reddit(), 
		}
	}
}

// New creates a base engine, *without* any defaults.
// Default port 8080 is used. No news sources are provided.
func New() *Engine {
	return &Engine{}
}

// Provider defines the only requirement from a source: a Fetch() function that returns
// a slice of NewsLink structs, and a Name().
type Provider interface {
	Fetch() (*[]WebPager, error)
	Name() string
}

// WebPager represents a link with a set of fields needed to be rendered by the news aggregator.
type WebPager interface {
	Source() string
	Title() string
	URL() string
	CommentsURL() string
	CommentCount() int
}

// newsSource represents the higher-level source of links for rendering.
type newsSource struct {
	source      Provider
	links       *[]WebPager
	errCount    int
	lastUpdated time.Time
}

func (s *newsSource) sleep(done chan *newsSource) {
	time.Sleep(fetchInterval*time.Second + time.Duration(s.errCount))
	done <- s
}
