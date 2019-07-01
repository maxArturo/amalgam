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

// Run starts the amalgam server.
func Run(port string, sources ...amalgam.Provider) {

	if len(sources) == 0 {
		log.Println("Using default news sources...")
		sources = []amalgam.Provider{
			reddit.New(),
			hackernews.New(),
		}
	}

	updated := worker.Start(sources)

	// handle new content coming in
	layoutHandler := contentHandler(updated)

	http.HandleFunc("/", layoutHandler)
	log.Fatal(http.ListenAndServe(util.ResolveAddress(port), nil))
}
