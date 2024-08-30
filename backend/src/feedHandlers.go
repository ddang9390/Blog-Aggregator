package main

import (
	"blog-aggregator/backend/internal/database"
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

func createFeed(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	var f feed

	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// JWT validation for user authentication
	userID, err := jwtValidate(r, cfg.jwtSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	feedID := uuid.New().String()

	ctx := r.Context()
	_, err2 := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:     feedID,
		Name:   f.Name,
		Url:    f.Url,
		UserID: userID,
	})
	if err2 != nil {
		fmt.Println(err2)
		http.Error(w, "Issue creating feed", http.StatusBadRequest)
		return
	}

	//Create feed follow when creating feed
	cfg.DB.CreateFeedFollows(ctx, database.CreateFeedFollowsParams{
		UserID: userID,
		FeedID: feedID,
	})

	response := map[string]interface{}{
		"id":         feedID,
		"created_at": time.Now(),
		"updated_at": time.Now(),
		"name":       f.Name,
		"url":        f.Url,
		"user_id":    userID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func getAllFeeds(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	// JWT validation for user authentication
	userID, err := jwtValidate(r, cfg.jwtSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	user := getUserHelper(cfg, w, r, userID)

	ctx := r.Context()
	databaseFeeds, err := cfg.DB.GetAllFeeds(ctx)
	feeds := make([]feed, len(databaseFeeds))
	for i, dbFeed := range databaseFeeds {
		feeds[i] = feed{
			Name: dbFeed.Name,
			Url:  dbFeed.Url,
			ID:   dbFeed.ID,
		}
	}
	if err != nil {
		http.Error(w, "Issue getting feeds", http.StatusInternalServerError)
		return
	}

	// response := map[string]interface{}{
	// 	"feeds": feeds,
	// }

	fmt.Println(feeds)
	//json.NewEncoder(w).Encode(response)
	data := pageData{
		User:  user,
		Feeds: feeds,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	outputHTML(w, "../../frontend/feeds.html", data)

}

func createFeedFollow(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var f feed_follow

		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// JWT validation for user authentication
		userID, err := jwtValidate(r, cfg.jwtSecret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		feedFollowID := uuid.New()
		response := map[string]interface{}{
			"id":         feedFollowID,
			"created_at": time.Now(),
			"updated_at": time.Now(),
			"feed_id":    f.Feed_id,
			"user_id":    userID,
		}
		ctx := r.Context()
		cfg.DB.CreateFeedFollows(ctx, database.CreateFeedFollowsParams{
			UserID: userID,
			FeedID: f.Feed_id,
		})

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func deleteFeedFollow(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func getAllFeedFollowsForUser(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	// JWT validation for user authentication
	userID, err := jwtValidate(r, cfg.jwtSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	feeds, err := cfg.DB.GetAllFeedFollowsForUser(ctx, userID)
	if err != nil {
		http.Error(w, "Issue getting feed follows", http.StatusInternalServerError)
		return
	}
	user := getUserHelper(cfg, w, r, userID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(feeds)

	data := pageData{
		User: user,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	outputHTML(w, "../../frontend/feeds.html", data)
}
