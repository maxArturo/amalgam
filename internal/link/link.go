package link

import (
	"github.com/maxArturo/amalgam"

	"github.com/go-shiori/go-readability"
)

type RenderedLinker interface {
	amalgam.Linker
	ParsedText() string
}

type RenderedLink struct {
	source       string
	title        string
	url          string
	commentsURL  string
	commentCount int
	readability.Article
}

func New(link amalgam.Linker) RenderedLink {
	return RenderedLink{
		source:       link.Source(),
		title:        link.Title(),
		url:          link.URL(),
		commentsURL:  link.CommentsURL(),
		commentCount: link.CommentCount(),
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
