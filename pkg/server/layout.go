package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/maxArturo/amalgam/internal/cache"

	"github.com/maxArturo/amalgam/internal/link"
)

type linkLayout struct {
	SourceName  string
	LastUpdated time.Time
	Links       []link.RenderedLinker
}

func mainLayout() *template.Template {
	return template.Must(template.ParseGlob("config/templates/main.html"))
}

func frontPageHandler(newContent chan *map[string][]link.RenderedLinker) func(w http.ResponseWriter, r *http.Request) {
	var latestContent []linkLayout

	go func() {
		for content := range newContent {
			latestContent = []linkLayout{}
			for source, links := range *content {
				latestContent = append(latestContent, linkLayout{
					SourceName:  source,
					Links:       links,
					LastUpdated: time.Now().UTC(),
				})
			}
		}
	}()

	layout := mainLayout()
	return func(w http.ResponseWriter, r *http.Request) {
		layout.Execute(w, latestContent)
	}
}

type linkView struct{}

// contentHandler parses out all links from a given Provider when it updates.
func (l *linkView) newHandler(in chan cache.Cacher) func(w http.ResponseWriter, r *http.Request) {
	out := make(chan *map[string][]link.RenderedLinker)

	go func() {
		for c := range in {
			items := c.Items()
			latestLinks := make(map[string][]link.RenderedLinker)

			for _, item := range items {
				if found := latestLinks[item.Source()]; found == nil {
					latestLinks[item.Source()] = []link.RenderedLinker{}
				}
				latestLinks[item.Source()] = append(latestLinks[item.Source()], item)
			}

			out <- &latestLinks
		}
	}()
	return frontPageHandler(out)
}
