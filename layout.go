package main

import (
	"html/template"
	"time"
)

type linkLayout struct {
	LastUpdated time.Time
	Links       *[]NewsLink
}

func mainLayout() *template.Template {
	return template.Must(template.ParseFiles("main.html"))
}
