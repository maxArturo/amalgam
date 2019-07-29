package server

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/maxArturo/amalgam/internal/link"
)

type linkLayout struct {
	SourceName  string
	LastUpdated time.Time
	Links       []link.RenderedLinker
}

func mainLayout() *template.Template {
	cwd, _ := os.Executable()
	return template.Must(template.ParseFiles(filepath.Join(filepath.Dir(cwd), "config/templates/main.html")))
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
func (l *linkView) newHandler(in chan *[]link.RenderedLinker) func(w http.ResponseWriter, r *http.Request) {
	out := make(chan *map[string][]link.RenderedLinker)
	latestLinks := make(map[string][]link.RenderedLinker)

	go func() {
		for s := range in {
			src := *s
			if linkLen := len(src); linkLen != 0 {

				name := src[0].Source()
				latestLinks[name] = src

				out <- &latestLinks
			}
		}
	}()
	return frontPageHandler(out)
}
