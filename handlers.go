package main

import "net/http"

func readyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func errorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
