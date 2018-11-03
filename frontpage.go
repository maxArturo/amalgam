package main

// FrontPage is the main struct for the landing page.
type FrontPage struct {
	title string
	links []newsLink
}

type newsLink struct {
	Title        string
	URL          string
	Source       string
	CommentsURL  string
	CommentCount string
}

// DefaultFrontPage provides the default template for the landing page.
// It should be replaced with content fetched from providers.
func DefaultFrontPage() *FrontPage {
	return &FrontPage{title: "Amalgam", links: []newsLink{{Title: "Loading..."}}}
}
