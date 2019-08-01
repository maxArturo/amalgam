package server

import (
	"log"
	"net/http"

	"github.com/maxArturo/amalgam/internal/cache"
	"github.com/maxArturo/amalgam/internal/handler"

	"github.com/maxArturo/amalgam/internal/provider/hackernews"
	"github.com/maxArturo/amalgam/internal/provider/reddit"
	"github.com/maxArturo/amalgam/internal/util"
	"github.com/maxArturo/amalgam/internal/worker"

	"github.com/maxArturo/amalgam"
)

type fetcher interface {
	Start(providers *[]amalgam.Provider) chan cache.Cacher
}

type layoutHandler interface {
	newHandler(in chan cache.Cacher) func(w http.ResponseWriter, r *http.Request)
}

type portResolver interface {
	ResolveAddress(addr string) string
}

type logger interface {
	fatal(v ...interface{})
}

type osLogger struct{}

func (f *osLogger) fatal(v ...interface{}) {
	log.Fatal(v...)
}

type httpServer interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ListenAndServe(addr string, handler http.Handler) error
}

// Server is the main Amalgam news aggregator.
type Server struct {
	fetcher
	layoutHandler
	portResolver
	httpServer
	logger
	cacher           cache.Cacher
	defaultProviders *[]amalgam.Provider
}

type osHTTPMux struct{}

func (m *osHTTPMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}

func (m *osHTTPMux) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

// New creates a configured server.
func New() *Server {
	c := cache.New()

	return &Server{
		fetcher:       worker.New(c),
		layoutHandler: &linkView{},
		portResolver:  util.New(),
		defaultProviders: &[]amalgam.Provider{
			reddit.New(),
			hackernews.New(),
		},
		cacher:     c,
		httpServer: &osHTTPMux{},
		logger:     &osLogger{},
	}
}

// Run starts the amalgam server.
func (s *Server) Run(port string, sources ...amalgam.Provider) {
	providers := sources
	if len(sources) == 0 {
		log.Println("No providers given. Using default news sources...")
		providers = *s.defaultProviders
	}

	updated := s.fetcher.Start(&providers)

	// handle new content coming in
	layout := s.layoutHandler.newHandler(updated)

	s.httpServer.HandleFunc("/text/", handler.TextHandler(s.cacher))
	s.httpServer.HandleFunc("/", layout)
	s.logger.fatal(s.httpServer.ListenAndServe(s.portResolver.ResolveAddress(port), nil))
}
