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

const numFetchers = 2

func defaultFrontPage() string {
	return "<h1>Loading content...</h1>"
}

func frontPageHandler(newContent chan string) func(w http.ResponseWriter, r *http.Request) {
	latestContent := defaultFrontPage()
	defaultHeader := "<h1>LInks@@!!</h1>"

	go func() {
		for content := range newContent {
			latestContent = defaultHeader + content
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handling request")
		fmt.Fprintf(w, latestContent)
	}
}

func contentHandler(in chan *newsSource) chan string {
	out := make(chan string)

	go func() {
		for src := range in {
			// TODO make this use templates instead
			content := "<ul>"
			for _, link := range *src.links {
				content = content + fmt.Sprintf("<li><a href=%s>%s</a>. <a href=%s>[%d]</a> </li>",
					link.URL, link.Title, link.CommentsURL, link.CommentCount)
			}
			content = content + "</ul>"

			out <- content
		}
	}()
	return out
}

// Fetcher waits on an incoming channel for Sources and fetches them, to update with new links.
// It reports out on the outgoing and content channels for completed Sources.
func Fetcher(in chan *newsSource, out chan *newsSource, content chan *newsSource) {
	for linkSrc := range in {
		links, err := linkSrc.source.Fetch()
		if err != nil {
			linkSrc.errCount++
		} else {
			linkSrc.links = links
			content <- linkSrc
		}
		out <- linkSrc
	}
}

func main() {
	port := os.Getenv("PORT")

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
		pending <- &newsSource{source: hackerNewsSource()}
	}()

	http.HandleFunc("/", frontPageHandler(newContent))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
