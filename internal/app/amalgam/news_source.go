package main

import (
	"time"
)

const errTimeoutDelay = 5
const fetchInterval = 30

// Sourcer defines the only requirement from a source: a Fetch() function that returns
// a slice of NewsLink structs, and a Name().
type Sourcer interface {
	Fetch() (*[]NewsLink, error)
	Name() string
}

// NewsLink represents a link with a set of fields needed to be rendered by the news aggregator.
type NewsLink struct {
	Source       string
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
	time.Sleep(fetchInterval*time.Second + time.Duration(s.errCount))
	done <- s
}
