package main

import (
	"blog-aggregator/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type feed struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func createFeed(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var f feed

		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		apiString := r.Header.Get("ApiKey")
		if apiString == "" {
			http.Error(w, "Api key required", http.StatusUnauthorized)
			return
		}

		user, err := getUser(cfg, w, r)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Issue getting user", http.StatusUnauthorized)
			return
		}

		response := map[string]interface{}{
			"id":         user.ApiKey,
			"created_at": time.Now(),
			"updated_at": time.Now(),
			"name":       f.Name,
			"url":        f.Url,
			"user_id":    user.ID,
		}

		ctx := r.Context()
		_, err2 := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
			Name:   sql.NullString{String: f.Name, Valid: true},
			Url:    sql.NullString{String: f.Url, Valid: true},
			UserID: user.ID,
		})
		if err2 != nil {
			fmt.Println(err2)
			http.Error(w, "Issue creating feed", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
