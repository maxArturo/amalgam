package main

import (
	"fmt"
	"time"
)

// FeedProvider represents a given source of links for consumption. It also
// contains relevant information for fetching and parsing links.
type FeedProvider struct {
	name          string
	apiURL        string
	links         []newsLink
	parseResponse func([]byte) bool
}

func (p *FeedProvider) getLinks(c chan bool) {
	c <- true
}

func (p *FeedProvider) startTimer(c chan []byte) {
	ticker := time.NewTicker(60 * time.Second)

	go func() {
		for range ticker.C {
			fmt.Println("querying %s...", p.name)
			// TODO query HN api instead of counting up
			// 	i++
			// 	select {
			// 	case c <- i:
			// 	default:
			// 		// empty out the single-buffer channel to keep it updated with the latest data
			// 		<-c
			// 		c <- i
			// 	}
		}
	}()
}
