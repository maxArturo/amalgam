package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/maxArturo/amalgam/internal/cache"
)

func TextHandler(c cache.Cacher) func(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Executable()
	textTemplate := template.Must(template.ParseFiles(filepath.Join(filepath.Dir(cwd), "config/templates/text.html")))
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.URL.Path[len("/text/"):]
		link, found := c.Get(hash)
		if !found {
			fmt.Fprintf(w, "<div>link text not found</div>")
		}
		textTemplate.Execute(w, link)
	}
}
