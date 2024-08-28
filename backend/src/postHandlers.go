package main

import (
	"blog-aggregator/backend/internal/database"
	"fmt"
	"net/http"
)

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
