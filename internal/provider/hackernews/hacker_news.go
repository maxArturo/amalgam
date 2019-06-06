package hackernews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/maxArturo/amalgam"
)

// Link represents a link to a HN post
type Link struct {
	source       string
	title        string
	url          string
	commentCount int
	commentsURL  string
}

// Source is the HN name
func (l Link) Source() string {
	return l.source
}

// Title is the article title
func (l Link) Title() string {
	return l.title
}

// URL is the link URL
func (l Link) URL() string {
	return l.url
}

// CommentsURL is the link's Comments URL (if any)
func (l Link) CommentsURL() string {
	return l.commentsURL
}

// CommentCount is the link's Comments URL (if any)
func (l Link) CommentCount() int {
	return l.commentCount
}

type hackerNewsResponse struct {
	Hits []hnHit `json:"hits"`
}

type hnHit struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	CommentCount int    `json:"num_comments"`
	ObjID        string `json:"objectID"`
}

// HackerNews is the base struct for the provider
type HackerNews struct {
	name   string
	APIURL string
}

// Fetch collets new links for processing
func (s *HackerNews) Fetch() ([]amalgam.Linker, error) {
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

	response, err := s.parseResponse(body)
	if err != nil {
		log.Println("Error reading HN response", s.APIURL, err)
		return nil, err
	}

	links := make([]amalgam.Linker, len(response))
	for i, v := range links {
		links[i] = v
	}
	return links, nil
}

// Name is the provider's official name
func (s *HackerNews) Name() string {
	return s.name
}

// New provides a configured HN provider
func New() *HackerNews {
	return &HackerNews{
		name:   "Hacker News",
		APIURL: "https://hn.algolia.com/api/v1/search?tags=front_page",
	}
}

func (s *HackerNews) parseResponse(body []byte) ([]*Link, error) {
	resp := &hackerNewsResponse{}
	err := json.Unmarshal(body, resp)
	if err != nil {

		log.Println("Error pasing HN response JSON", err)
		return nil, err
	}

	links := []*Link{}
	for _, link := range resp.Hits {
		commentURL := fmt.Sprintf("https://news.ycombinator.com/item?id=%s", link.ObjID)
		links = append(links,
			&Link{
				source:       s.name,
				title:        link.Title,
				url:          link.URL,
				commentCount: link.CommentCount,
				commentsURL:  commentURL})
	}
	return links, nil
}
