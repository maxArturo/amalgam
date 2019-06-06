package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// HNLink a link to a HN post
type HNLink struct {
	source       string
	title        string
	url          string
	commentCount int
	commentsURL  string
}

// Source the HN name
func (l HNLink) Source() string {
	return l.source
}

// Title the article title
func (l HNLink) Title() string {
	return l.title
}

// URL the link URL
func (l HNLink) URL() string {
	return l.url
}

// CommentsURL the link's Comments URL (if any)
func (l HNLink) CommentsURL() string {
	return l.commentsURL
}

// CommentCount the link's Comments URL (if any)
func (l HNLink) CommentCount() int {
	return l.commentCount
}

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

func (s *hackerNews) Fetch() (*[]HNLink, error) {
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

// HackerNews provides the HN source
func HackerNews() *hackerNews {
	return &hackerNews{
		name:   "Hacker News",
		APIURL: "https://hn.algolia.com/api/v1/search?tags=front_page",
	}
}

func (s *hackerNews) parseResponse(body []byte) (*[]HNLink, error) {
	resp := &hackerNewsResponse{}
	err := json.Unmarshal(body, resp)
	if err != nil {

		log.Println("Error pasing HN response JSON", err)
		return nil, err
	}

	links := []HNLink{}
	for _, link := range resp.Hits {
		commentURL := fmt.Sprintf("https://news.ycombinator.com/item?id=%s", link.ObjID)
		links = append(links,
			HNLink{
				source:       s.name,
				title:        link.Title,
				url:          link.URL,
				commentCount: link.CommentCount,
				commentsURL:  commentURL})
	}
	return &links, nil
}
