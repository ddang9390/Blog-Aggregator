package main

import (
	"blog-aggregator/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type feed struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	ID   string `json:"id"`
}

type feed_follow struct {
	Feed_id string `json:"feed_id"`
}

func createFeed(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var f feed

		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user := getUserHelper(cfg, w, r)

		feedID := uuid.New().String()
		response := map[string]interface{}{
			"id":         feedID,
			"created_at": time.Now(),
			"updated_at": time.Now(),
			"name":       f.Name,
			"url":        f.Url,
			"user_id":    user.ID,
		}

		ctx := r.Context()
		_, err2 := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
			ID:     feedID,
			Name:   f.Name,
			Url:    f.Url,
			UserID: user.ID,
		})
		if err2 != nil {
			fmt.Println(err2)
			http.Error(w, "Issue creating feed", http.StatusUnauthorized)
			return
		}

		//Create feed follow when creating feed
		cfg.DB.CreateFeedFollows(ctx, database.CreateFeedFollowsParams{
			UserID: user.ID,
			FeedID: feedID,
		})

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func getAllFeeds(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		feeds, err := cfg.DB.GetAllFeeds(ctx)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Issue getting feeds", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"feeds": feeds,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func createFeedFollow(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var f feed_follow

		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user := getUserHelper(cfg, w, r)
		if user == nil {
			http.Error(w, "Issue getting user", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"id":         user.ApiKey,
			"created_at": time.Now(),
			"updated_at": time.Now(),
			"feed_id":    f.Feed_id,
			"user_id":    user.ID,
		}
		ctx := r.Context()
		cfg.DB.CreateFeedFollows(ctx, database.CreateFeedFollowsParams{
			UserID: user.ID,
			FeedID: f.Feed_id,
		})

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func deleteFeedFollow(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)

		parsedURL, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, "Issue parsing url", http.StatusInternalServerError)
			return
		}
		feed_id := parsedURL.Query().Get("feed_id")

		ctx := r.Context()
		cfg.DB.DeleteFeedFollow(ctx, feed_id)
	}
}

func getAllFeedFollowsForUser(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authroization header required", http.StatusUnauthorized)
			return
		}
		user := getUserHelper(cfg, w, r)
		if user == nil {
			http.Error(w, "Issue getting user", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		feeds, err := cfg.DB.GetAllFeedFollowsForUser(ctx, user.ID)
		if err != nil {
			http.Error(w, "Issue getting feed follows", http.StatusInternalServerError)
			return
		}

		fmt.Println(feeds)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(feeds)
	}
}
