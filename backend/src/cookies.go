package main

import (
	"net/http"
)

func setCookieHandler(w http.ResponseWriter, r *http.Request, cfg *apiConfig, userID string) {
	session, err := createSession(cfg, userID, r)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:     "blog-aggregator",
		Value:    session.sessionID,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("blog-aggregator")
	if err != nil {
		http.Error(w, "Error getting cookie", http.StatusBadRequest)
		return ""
	}

	return cookie.Value
}
