package handler

import (
	"fmt"
	"net/http"

	"github.com/maxArturo/amalgam/internal/cache"
)

func TextHandler(c cache.Cacher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.URL.Path[len("/text/"):]
		link, found := c.Get(hash)
		if !found {
			fmt.Fprintf(w, "<div>link text not found</div>")
		}
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", link.Title(), link.LinkText())
	}
}
