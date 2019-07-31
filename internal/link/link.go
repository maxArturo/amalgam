package link

import (
	b64 "encoding/base64"
	"time"

	"github.com/maxArturo/amalgam"

	"github.com/go-shiori/go-readability"
)

type RenderedLinker interface {
	ParsedText() string
	Hash() string
	FetchedAt() time.Time
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

func (l RenderedLink) ParsedText() string {
	return l.TextContent
}

func (l RenderedLink) Hash() string {
	return l.hash
}

func (l RenderedLink) FetchedAt() time.Time {
	return l.fetchedAt
}
