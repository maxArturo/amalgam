package main

import (
	"github.com/turnage/graw/reddit"
	"time"
)

type redditSource struct {
	name    string
	harvest reddit.Harvest
}

func (s *redditSource) Fetch() (*[]NewsLink, error) {
	links := []NewsLink{}
	posts := s.harvest.Posts

	for _, post := range posts {
		links = append(links,
			NewsLink{
				Title:        post.Title,
				URL:          post.URL,
				CommentCount: len(post.Replies),
				CommentsURL:  post.Permalink},
		)
	}

	return &links, nil
}

func (s *redditSource) Name() string {
	return s.name
}

func newRedditSource() *redditSource {
	rate := 60 * time.Second
	script, _ := reddit.NewScript("amalgam:reddit_script:0.0.1 by /u/aadvark_dev", rate)
	harvest, _ := script.Listing("r/programming/hot", "")

	return &redditSource{
		name:    "reddit",
		harvest: harvest,
	}
}
