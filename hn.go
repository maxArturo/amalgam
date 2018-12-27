package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// HackerNewsResponse represents the Hacker News website source for news links.
type hackerNewsResponse struct {
	Hits []hnHit `json:"hits"`
}

type hnHit struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	CommentCount int    `json:"num_comments"`
	ObjID        string `json:"objectID"`
}

type hackerNews struct {
	name   string
	APIURL string
}

func (s *hackerNews) Fetch() (*[]NewsLink, error) {
	log.Println("[HN] querying HN api...")
	resp, err := http.Get(s.APIURL)
	if err != nil {
		log.Println("Error fetching url", s.APIURL, err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading HN response", s.APIURL, err)
		return nil, err
	}
	return s.parseResponse(body)
}

func (s *hackerNews) Name() string {
	return s.name
}

func hackerNewsSource() *hackerNews {
	return &hackerNews{
		name:   "Hacker News",
		APIURL: "https://hn.algolia.com/api/v1/search?tags=front_page",
	}
}

func (s *hackerNews) parseResponse(body []byte) (*[]NewsLink, error) {
	resp := &hackerNewsResponse{}
	err := json.Unmarshal(body, resp)
	if err != nil {

		log.Println("Error pasing HN response JSON", err)
		return nil, err
	}

	links := []NewsLink{}
	for _, link := range resp.Hits {
		commentURL := fmt.Sprintf("https://news.ycombinator.com/item?id=%s", link.ObjID)
		links = append(links,
			NewsLink{
				Title:        link.Title,
				URL:          link.URL,
				CommentCount: link.CommentCount,
				CommentsURL:  commentURL})
	}
	return &links, nil
}
