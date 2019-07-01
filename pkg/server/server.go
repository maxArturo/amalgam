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

type server struct {
	queue        fetcher
	layoutRender layoutHandler
}

func new() *server {
	return &server{
		queue:        &worker.SourceJob{},
		layoutRender: &linkView{},
	}
}

// Run starts the amalgam server.
func Run(port string, sources ...amalgam.Provider) {
	if len(sources) == 0 {
		log.Println("Using default news sources...")
		sources = []amalgam.Provider{
			reddit.New(),
			hackernews.New(),
		}
	}

	server := new()
	updated := server.queue.Start(sources)

	// handle new content coming in
	handler := server.layoutRender.newHandler(updated)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(util.ResolveAddress(port), nil))
}
