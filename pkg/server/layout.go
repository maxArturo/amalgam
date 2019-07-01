package server

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/maxArturo/amalgam"
)

type linkLayout struct {
	SourceName  string
	LastUpdated time.Time
	Links       []amalgam.Linker
}

func mainLayout() *template.Template {
	cwd, _ := os.Executable()
	return template.Must(template.ParseFiles(filepath.Join(filepath.Dir(cwd), "config/templates/main.html")))
}

func frontPageHandler(newContent chan *map[string][]amalgam.Linker) func(w http.ResponseWriter, r *http.Request) {
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
func (l *linkView) newHandler(in chan []amalgam.Linker) func(w http.ResponseWriter, r *http.Request) {
	out := make(chan *map[string][]amalgam.Linker)
	latestLinks := make(map[string][]amalgam.Linker)

	go func() {
		for s := range in {
			if linkLen := len(s); linkLen != 0 {

				name := s[0].Source()
				latestLinks[name] = s

				out <- &latestLinks
			}
		}
	}()
	return frontPageHandler(out)
}
