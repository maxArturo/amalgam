package main

import (
	"html/template"
	"os"
	"path/filepath"
	"time"
)

type linkLayout struct {
	LastUpdated time.Time
	Links       *[]NewsLink
}

func mainLayout() *template.Template {
	cwd, _ := os.Executable()
	return template.Must(template.ParseFiles(filepath.Join(filepath.Dir(cwd), "templates/main.html")))
}
