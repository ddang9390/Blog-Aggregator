package main

import (
	"html/template"
	"net/http"
)

type pageData struct {
	User         User
	Feeds        []feed
	Feed_follows []feed_follow
	Posts        []post
}

var navbar = "../../frontend/navbar.html"

func outputHTML(w http.ResponseWriter, file string, data pageData) {
	template, err := template.ParseFiles(file, navbar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = template.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
