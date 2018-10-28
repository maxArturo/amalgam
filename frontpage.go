package main

// FrontPage is the main struct for the landing page.
type FrontPage struct {
	title string
	links []newsLink
}

type newsLink struct {
	title  string
	url    string
	source string
}

// DefaultFrontPage provides the default template for the landing page.
// It should be replaced with content fetched from providers.
func DefaultFrontPage() *FrontPage {
	return &FrontPage{title: "Amalgam", links: []newsLink{{title: "Loading..."}}}
}
