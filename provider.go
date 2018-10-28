package main

// FeedProvider represents a given source of links for consumption.
type FeedProvider struct {
	name          string
	apiURL        string
	links         []newsLink
	parseResponse func([]byte) bool
}

type newsLink struct {
	title string
	url   string
}

func (p *FeedProvider) getLinks(c chan bool) {
	c <- true
}
