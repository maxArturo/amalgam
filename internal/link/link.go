package link

import (
	"time"

	"github.com/maxArturo/amalgam"
)

type RenderedLinker interface {
	Link() amalgam.Linker
	ParsedText() string
}

type RenderedLink struct {
	link      amalgam.Linker
	Text      string
	FetchedAt time.Time
}

func New(link amalgam.Linker) RenderedLink {
	return RenderedLink{
		link: link,
	}
}

func (l RenderedLink) ParsedText() string {
	return l.Text
}

func (l RenderedLink) Link() amalgam.Linker {
	return l.link
}
