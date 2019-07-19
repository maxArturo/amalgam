package reddit

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/maxArturo/amalgam"
	"github.com/valyala/fastjson"
)

// Link represents a link in Reddit
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

// Reddit is the base struct for the provider to Reddit
type Reddit struct {
	name   string
	APIURL string
}

// Fetch queries Reddit for new links
func (s *Reddit) Fetch() (*[]amalgam.Linker, error) {
	log.Println("[REDDIT] querying reddit api...")
	req, err := http.NewRequest("GET", s.APIURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "amalgam:reddit_script:0.0.1 by /u/aadvark_dev")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error fetching url", s.APIURL, err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading reddit response", s.APIURL, err)
		return nil, err
	}

	response, err := s.parseResponse(body)

	if err != nil {
		log.Println("Error reading HN response", s.APIURL, err)
		return nil, err
	}

	links := make([]amalgam.Linker, len(*response))
	for i := range links {
		links[i] = (*response)[i]
	}
	return &links, nil
}

// Name is the name to display for the Reddit source
func (s *Reddit) Name() string {
	return s.name
}

// New returns a configured Reddit provider
func New() *Reddit {
	return &Reddit{
		name:   "reddit",
		APIURL: "https://www.reddit.com/r/programming/hot/.json?top=20",
	}
}

func (s *Reddit) parseResponse(body []byte) (*[]Link, error) {
	var p fastjson.Parser
	v, err := p.Parse(string(body))
	if err != nil {
		log.Println("Error reading reddit response", s.APIURL, err)
		return nil, err
	}

	listingArr := v.GetArray("data", "children")
	links := []Link{}

	for _, listing := range listingArr {
		commentURL := fmt.Sprintf("https://www.reddit.com%s",
			string(listing.GetStringBytes("data", "permalink")))

		links = append(links, Link{
			source:       s.name,
			title:        string(listing.GetStringBytes("data", "title")),
			url:          string(listing.GetStringBytes("data", "url")),
			commentCount: listing.GetInt("data", "num_comments"),
			commentsURL:  commentURL})
	}

	return &links, nil
}
