// Package amalgam provides a brutally simple webserver for no-nonsense news aggregation.
// It aims for lightweight hardware requirements, extensibility, and simplicity.

package amalgam

import (
	"log"
	"net/http"
	"os"
	"time"
)

const numFetchers = 3

func frontPageHandler(newContent chan *[]NewsLink) func(w http.ResponseWriter, r *http.Request) {
	var latestContent *linkLayout

	go func() {
		for content := range newContent {
			latestContent = &linkLayout{
				LastUpdated: time.Now().UTC(),
				Links:       content,
			}
		}
	}()

	layout := mainLayout()
	return func(w http.ResponseWriter, r *http.Request) {
		layout.Execute(w, latestContent)
	}
}

func contentHandler(in chan *newsSource) chan *[]NewsLink {
	out := make(chan *[]NewsLink)
	latestLinks := make(map[string]*[]NewsLink)

	go func() {
		for s := range in {
			latestLinks[s.source.Name()] = s.links
			renderLinks := make([]NewsLink, 0)
			for _, links := range latestLinks {
				renderLinks = append(renderLinks, *links...)
			}

			out <- &renderLinks
		}
	}()
	return out
}

// Fetcher waits on an incoming channel for Sources and fetches them, to update with new links.
// It reports out on the outgoing and content channels for completed Sources.
func Fetcher(label int, in chan *newsSource, out chan *newsSource, content chan *newsSource) {
	for src := range in {
		log.Printf("[FETCH] fetcher no %d, fetching for %s", label, src.source.Name())
		newLinks, err := src.source.Fetch()
		if err != nil {
			log.Println("[FETCH] Error fetching ", src.source.Name(), err)
			src.errCount++
		} else {
			src.errCount = 0
			src.links = newLinks
			log.Printf("in fechtcher no %d, source %s, link count %d", label, src.source.Name(), len(*src.links))
		}

		content <- src
		out <- src
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	sources := []newsSource{
		newsSource{source: redditSource()},
		newsSource{source: hackerNewsSource()},
	}

	// create our pending/done/new content channels
	pending, done, updated := make(chan *newsSource),
		make(chan *newsSource), make(chan *newsSource)

	// launch fetchers
	for i := 0; i < numFetchers; i++ {
		go Fetcher(i, pending, done, updated)
	}

	// handle new content coming in
	newContent := contentHandler(updated)

	go func() {
		for s := range done {
			newSource := s
			go newSource.sleep(pending)
		}
	}()

	go func() {
		for _, src := range sources {
			copySrc := src
			select {
			case pending <- &copySrc:
				time.Sleep(10 * time.Second)
			}
		}
	}()

	http.HandleFunc("/", frontPageHandler(newContent))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
