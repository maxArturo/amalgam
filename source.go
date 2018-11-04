package main

// Sourcer defines the only requirement from a source: a Fetch() function that returns
// a slice of NewsLink structs.
type Sourcer interface {
	Fetch() (*[]NewsLink, error)
	Name() string
}

// NewsLink represents a link with a set of fields needed to be rendered by the news aggregator.
type NewsLink struct {
	Title        string
	URL          string
	CommentsURL  string
	CommentCount string
}
