package server

import (
	"log"
	"net/http"

	"github.com/maxArturo/amalgam/internal/provider/hackernews"
	"github.com/maxArturo/amalgam/internal/provider/reddit"
	"github.com/maxArturo/amalgam/internal/util"
	"github.com/maxArturo/amalgam/internal/worker"

	"github.com/maxArturo/amalgam"
)

type fetcher interface {
	Start(providers []amalgam.Provider) chan []amalgam.Linker
}

type layoutHandler interface {
	newHandler(in chan []amalgam.Linker) func(w http.ResponseWriter, r *http.Request)
}

type portResolver interface {
	ResolveAddress(addr string) string
}

// Server is the main Amalgam news aggregator.
type Server struct {
	queue        fetcher
	layoutRender layoutHandler
	portResolver
}

// New creates a configured server.
func New() *Server {
	return &Server{
		queue:        worker.New(),
		layoutRender: &linkView{},
		portResolver: util.New(),
	}
}

// Run starts the amalgam server.
func (s *Server) Run(port string, sources ...amalgam.Provider) {
	if len(sources) == 0 {
		log.Println("Using default news sources...")
		sources = []amalgam.Provider{
			reddit.New(),
			hackernews.New(),
		}
	}

	updated := s.queue.Start(sources)

	// handle new content coming in
	handler := s.layoutRender.newHandler(updated)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(s.portResolver.ResolveAddress(port), nil))
}
