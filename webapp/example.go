package main

import (
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/example", exampleHandler)
}

var exampleTemplate = template.Must(template.ParseFiles("static/example.html"))

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	exampleTemplate.Execute(w, nil)
}
