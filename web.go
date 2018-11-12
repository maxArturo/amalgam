// Package amalgam provides a brutally simple webserver for no-nonsense news aggregation.
// It aims for lightweight hardware requirements, extensibility, and simplicity.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pkg/profile"
)

const numFetchers = 3

func defaultFrontPage() string {
	return "<h1>Loading content...</h1>"
}

func frontPageHandler(newContent chan string) func(w http.ResponseWriter, r *http.Request) {
	latestContent := defaultFrontPage()
	defaultHeader := "<h1>Latest Links</h1>"

	go func() {
		for content := range newContent {
			latestContent = defaultHeader + content
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling request")
		fmt.Fprintf(w, latestContent)
	}
}

func formatLinks(sourceLinks *map[string]*[]NewsLink) string {
	// TODO make this use templates instead
	content := "<ul>"

	for source, links := range *sourceLinks {
		for _, link := range *links {
			content = content + fmt.Sprintf("<li>[%s] <a href=%s>%s</a>. <a href=%s>[%d]</a> </li>",
				source, link.URL, link.Title, link.CommentsURL, link.CommentCount)
		}
	}

	content = content + "</ul>"
	return content
}

func contentHandler(in chan *newsSource) chan string {
	out := make(chan string)
	latestLinks := make(map[string]*[]NewsLink)

	go func() {
		for s := range in {
			latestLinks[s.source.Name()] = s.links
			out <- formatLinks(&latestLinks)
		}
	}()
	return out
}

// Fetcher waits on an incoming channel for Sources and fetches them, to update with new links.
// It reports out on the outgoing and content channels for completed Sources.
func Fetcher(in chan *newsSource, out chan *newsSource, content chan *newsSource) {
	for src := range in {
		newLinks, err := src.source.Fetch()
		if err != nil {
			src.errCount++
		} else {
			src.errCount = 0
			src.links = newLinks
			log.Println("counts of links for source, ", src.source.Name(), len(*src.links))
		}

		content <- src
		out <- src
	}
}

func main() {
	port := os.Getenv("PORT")
	sources := []newsSource{
		newsSource{source: hackerNewsSource()},
		// newsSource{source: newRedditSource()},
	}
	sourceCount := len(sources)

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	// memory profiling
	defer profile.Start(profile.MemProfile).Stop()

	// create our pending/done/new content channels
	pending, done, updated := make(chan *newsSource, sourceCount),
		make(chan *newsSource, sourceCount), make(chan *newsSource, sourceCount)

	// launch fetchers
	for i := 0; i < numFetchers; i++ {
		go Fetcher(pending, done, updated)
	}

	// handle new content coming in
	newContent := contentHandler(updated)

	go func() {
		for s := range done {
			go s.sleep(pending)
		}
	}()

	// Finally! Send the sources into pending
	go func() {
		for _, src := range sources {
			log.Println("adding source to pending queue: ", src.source.Name())
			pending <- &src
		}
	}()

	http.HandleFunc("/", frontPageHandler(newContent))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
