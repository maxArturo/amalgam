// Package amalgam provides a brutally simple webserver for no-nonsense news aggregation.
// It aims for lightweight hardware requirements, extensibility, and simplicity.

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/profile"
)

const (
	numFetchers     = 5
	errTimeoutDelay = 5
)

type update struct {
	name  string
	links *[]NewsLink
}

func defaultFrontPage() string {
	return "<h1>Loading content...</h1>"
}

func frontPageHandler(newContent chan string) func(w http.ResponseWriter, r *http.Request) {
	latestContent := defaultFrontPage()

	go func() {
		for content := range newContent {
			latestContent = content
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handling request")
		fmt.Fprintf(w, latestContent)
	}
}

func contentHandler(in chan update) chan string {
	base := "<h1> links! </h1>"
	out := make(chan string)

	go func() {
		out <- base
		for newContent := range in {
			out <- base + newContent.name
		}
	}()
	return out
}

// Fetcher waits on an incoming channel for Sources and fetches them, to update with new links.
// It reports out on the outgoing and content channels for completed Sources.
func Fetcher(in chan Sourcer, out chan Sourcer, content chan *update) {
	for src := range in {
		links, err := src.Fetch()
		out <- src
		content <- &update{name: src.Name(), links: links}
	}
}

func main() {
	// memory profiling
	defer profile.Start(profile.MemProfile).Stop()

	// create our pending/done/new content channels
	pending, done, updated := make(chan Sourcer), make(chan Sourcer), make(chan *update)

	// launch fetchers
	for i := 0; i < numFetchers; {
		go Fetcher(pending, done, updated)
	}

	// handle new content coming in
	newContent := contentHandler(updated)

	go func() {
		for s := range done {

		}
	}()

	http.HandleFunc("/", frontPageHandler(newContent))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
