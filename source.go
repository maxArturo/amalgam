package main

import (
	"log"
	"time"
)

const errTimeoutDelay = 5
const fetchInterval = 5

// Sourcer defines the only requirement from a source: a Fetch() function that returns
// a slice of NewsLink structs.
type Sourcer interface {
	Fetch() (*[]NewsLink, error)
	Name() string
}

type Sources *[]newsSource

// NewsLink represents a link with a set of fields needed to be rendered by the news aggregator.
type NewsLink struct {
	Title        string
	URL          string
	CommentsURL  string
	CommentCount int
}

// newsSource represents the higher-level source of links for rendering.
type newsSource struct {
	source      Sourcer
	links       *[]NewsLink
	errCount    int
	lastUpdated time.Time
}

func (s *newsSource) sleep(done chan *newsSource) {
	log.Println("sleeping for some seconds...")
	time.Sleep(fetchInterval*time.Second + time.Duration(s.errCount))
	done <- s
}
