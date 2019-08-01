package link

import (
	b64 "encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/maxArturo/amalgam"

	"github.com/go-shiori/go-readability"
)

type RenderedLinker interface {
	LinkText() string
	Hash() string
	FetchedAt() time.Time
	FetchLinkText()
	amalgam.Linker
}

type RenderedLink struct {
	source       string
	title        string
	url          string
	commentsURL  string
	commentCount int
	hash         string
	fetchedAt    time.Time
	readability.Article
}

func New(link amalgam.Linker) RenderedLink {
	return RenderedLink{
		source:       link.Source(),
		title:        link.Title(),
		url:          link.URL(),
		commentsURL:  link.CommentsURL(),
		commentCount: link.CommentCount(),
		hash:         b64.StdEncoding.EncodeToString([]byte(link.URL())),
		fetchedAt:    time.Now(),
	}
}

func (l RenderedLink) Source() string {
	return l.source
}

func (l RenderedLink) Title() string {
	return l.title
}

func (l RenderedLink) URL() string {
	return l.url
}

func (l RenderedLink) CommentsURL() string {
	return l.commentsURL
}

func (l RenderedLink) CommentCount() int {
	return l.commentCount
}

func (l RenderedLink) LinkText() string {
	return l.TextContent
}

func (l RenderedLink) Hash() string {
	return l.hash
}

func (l RenderedLink) FetchedAt() time.Time {
	return l.fetchedAt
}

func (l RenderedLink) FetchLinkText() {
	url := l.URL()
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to download %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	article, err := readability.FromReader(resp.Body, url)
	if err != nil {
		log.Printf("failed to parse %s: %v\n", l.URL(), err)
		return
	}

	l.Article = article
}
