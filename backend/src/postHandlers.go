package main

import (
	"blog-aggregator/backend/internal/database"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type post struct {
	Title       string `json:"title"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

func getPostsForUser(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := jwtValidate(r, cfg.jwtSecret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		posts, err := cfg.DB.GetPostsForUser(ctx, database.GetPostsForUserParams{
			UserID: userID,
			Limit:  10,
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(posts)
	}
}

func getPosts(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	// JWT validation for user authentication
	userID, err := jwtValidate(r, cfg.jwtSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	//TODO
	//find out how to get feed id
	vars := mux.Vars(r)
	feedID := vars["feedID"]
	fmt.Println(feedID)

	db_posts, err := cfg.DB.GetPostsByFeed(ctx, feedID)
	if err != nil {
		http.Error(w, "Issue getting posts", http.StatusInternalServerError)
		return
	}
	user := getUserHelper(cfg, w, r, userID)

	posts := make([]post, len(db_posts))
	for i, dbPost := range db_posts {
		posts[i] = post{
			Title:       dbPost.Title.String,
			Description: dbPost.Description.String,
			Url:         dbPost.Url,
		}
	}

	data := pageData{
		User:  user,
		Posts: posts,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	outputHTML(w, "../../frontend/posts.html", data)
}
