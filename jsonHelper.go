package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(payload)

	if err != nil {
		fmt.Printf("Error marshalling JSON: %s"+"\n", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(body)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type err struct {
		Error string `json:"error"`
	}

	e := &err{
		Error: msg,
	}

	respondWithJSON(w, code, e)
}
