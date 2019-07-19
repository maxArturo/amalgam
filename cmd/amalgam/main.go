package main

import (
	"github.com/maxArturo/amalgam/pkg/server"
)

func main() {
	s := server.New()
	s.Run("")
}
