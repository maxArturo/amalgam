package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/valyala/fastjson"
)

type reddit struct {
	name   string
	APIURL string
}

func (s *reddit) Fetch() (*[]NewsLink, error) {
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
	return s.parseResponse(body)
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

func (s *reddit) parseResponse(body []byte) (*[]NewsLink, error) {
	var p fastjson.Parser
	v, err := p.Parse(string(body))
	if err != nil {
		log.Println("Error reading reddit response", s.APIURL, err)
		return nil, err
	}

	listingArr := v.GetArray("data", "children")
	links := []NewsLink{}

	for _, listing := range listingArr {
		commentURL := fmt.Sprintf("https://www.reddit.com%s",
			string(listing.GetStringBytes("data", "permalink")))

		links = append(links,
			NewsLink{
				Title:        string(listing.GetStringBytes("data", "title")),
				URL:          string(listing.GetStringBytes("data", "url")),
				CommentCount: listing.GetInt("data", "num_comments"),
				CommentsURL:  commentURL})
	}

	return &links, nil
}
