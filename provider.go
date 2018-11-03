package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Source represents a given source of links for consumption. It also
// contains relevant information for fetching and parsing links.
type Source struct {
	Name          string
	Abbreviation  string
	APIURL        string
	Links         []newsLink
	LastUpdatedAt time.Time
	ParseResponse func([]byte) ([]newsLink, error)
}

// Start kicks off the fetching and periodic refreshing of links with a default
// 20 minute interval between fetching of new links.
func (p *Source) Start() {
	p.StartWithDuration(time.Minute * 20)
}

// StartWithDuration kicks off the fetching and periodic refreshing of links with
// a given interval. Minimum is one second.
func (p *Source) StartWithDuration(duration time.Duration) {
	if duration < time.Minute {
		duration = time.Minute
	}

	ticker := time.NewTicker(duration)

	go func() {
		for range ticker.C {
			fmt.Println("querying %s...", p.Name)
			p.fetch()
		}
	}()
}

func (p *Source) fetch() {
	resp, err := http.Get(p.APIURL)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error parsing response body", err)
	}

	p.Links = p.ParseResponse(body)
}
