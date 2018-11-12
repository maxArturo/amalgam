package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// redditResponse represents the icoming payload from Reddit.
type redditResponse struct {
	Data []hnHit `json:"data"`
}

type hnHit struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	CommentCount int    `json:"num_comments"`
	ObjID        string `json:"objectID"`
}

type redditPost struct {
}

type reddit struct {
	name   string
	APIURL string
}

func (s *reddit) Fetch() (*[]NewsLink, error) {
	log.Println("querying reddit api...")
	resp, err := http.Get(s.APIURL)
	if err != nil {
		log.Println("Error fetching url", s.APIURL, err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading reddit response", s.APIURL, err)
		return nil, err
	}
	return parseResponse(body)
}

func (s *reddit) Name() string {
	return s.name
}

func redditSource() *reddit {
	return &reddit{
		name:   "reddit",
		APIURL: "https://www.reddit.com/r/programming/hot/.json?top=20",
	}
}

func parseResponse(body []byte) (*[]NewsLink, error) {
	s := &hackerNewsResponse{}
	err := json.Unmarshal(body, s)
	if err != nil {

		log.Println("Error pasing HN response JSON", err)
		return nil, err
	}

	links := []NewsLink{}
	for _, link := range s.Hits {
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
