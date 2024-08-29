package main

import (
	"html/template"
	"net/http"
)

func outputHTML(w http.ResponseWriter, file string, user User) {
	template, err := template.ParseFiles(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = template.Execute(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
