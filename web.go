// Package amalgam provides a brutally simple webserver for no-nonsense news aggregation.
// It aims for lightweight hardware requirements, extensibility, and simplicity.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/profile"
)

const numFetchers = 3

type renderLinks struct {
	source string
	link   NewsLink
}

func frontPageHandler(newContent chan string) func(w http.ResponseWriter, r *http.Request) {
	latestContent := "<h1>Loading content...</h1>"
	defaultHeader := "<h1>Latest Links</h1>"

	go func() {
		for content := range newContent {
			latestContent = defaultHeader + content
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
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
	sources := []newsSource{
		newsSource{source: redditSource()},
		newsSource{source: hackerNewsSource()},
	}

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	// memory profiling
	defer profile.Start(profile.MemProfile).Stop()

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
