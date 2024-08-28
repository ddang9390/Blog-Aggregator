package main

import (
	"net/http"
)

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "blog-aggregator",
		Value:    "hello world",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	w.Write([]byte("hello world"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("blog-aggregator")
	if err != nil {
		http.Error(w, "Error getting cookie", http.StatusBadRequest)
		return
	}

	w.Write([]byte(cookie.Value))
}
