package main

import (
	"encoding/json"
	"fmt"
)

type hackerNewsResponse struct {
	hits []hnHit `json: "hits"`
}

type hnHit struct {
	title        string `json: "title"`
	url          string `json: "url"`
	commentCount string `json: "num_comments"`
	objID        string `json: "objectID"`
}

// HackerNewsSource represents the Hacker News API source.
func HackerNewsSource() *Source {
	return &Source{
		Name:          "Hacker News",
		Abbreviation:  "HN",
		APIURL:        "https://hn.algolia.com/api/v1/search?tags=front_page",
		ParseResponse: parseResponse,
	}
}

func parseResponse(body []byte) ([]newsLink, error) {
	s := &hackerNewsResponse{}
	err := json.Unmarshal(body, s)
	if err != nil {
		return []newsLink{}, err
	}

	links := []newsLink{}
	for _, link := range s.hits {
		commentURL, _ := fmt.Printf("https://news.ycombinator.com/item?id=%s", link.objID)
		links := append(links,
			newsLink{
				Title:        link.title,
				URL:          link.url,
				Source:       "Hacker News",
				CommentCount: link.commentCount,
				CommentsURL:  commentURL})
	}
	return links, nil
}
