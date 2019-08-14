package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/maxArturo/amalgam/internal/cache"
)

func TextHandler(c cache.Cacher) func(w http.ResponseWriter, r *http.Request) {
	textTemplate := template.Must(template.ParseGlob("config/templates/text.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.URL.Path[len("/text/"):]
		link, found := c.Get(hash)
		if !found {
			fmt.Fprintf(w, "<div>link text not found</div>")
		}
		textTemplate.Execute(w, link)
	}
}
