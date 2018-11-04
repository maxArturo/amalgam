package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type hackerNews struct {
	name   string
	APIURL string
}

func (s *hackerNews) Fetch() (*[]NewsLink, error) {
	resp, err := http.Get(s.APIURL)
	if err != nil {
		log.Println("Error fetching url", s.APIURL, err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseResponse(body)
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

func parseResponse(body []byte) (*[]NewsLink, error) {
	s := &hackerNewsResponse{}
	err := json.Unmarshal(body, s)
	if err != nil {
		return nil, err
	}

	links := []NewsLink{}
	for _, link := range s.hits {
		commentURL := fmt.Sprintf("https://news.ycombinator.com/item?id=%s", link.objID)
		links := append(links,
			NewsLink{
				Title:        link.title,
				URL:          link.url,
				CommentCount: link.commentCount,
				CommentsURL:  commentURL})
	}
	return &links, nil
}
