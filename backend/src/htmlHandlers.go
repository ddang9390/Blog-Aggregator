package main

import (
	"html/template"
	"net/http"
)

type pageData struct {
	User  User
	Feeds []feed
}

func outputHTML(w http.ResponseWriter, file string, data pageData) {
	template, err := template.ParseFiles(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = template.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
