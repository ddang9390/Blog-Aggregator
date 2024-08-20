package main

import (
	"blog-aggregator/internal/database"
	"fmt"
	"net/http"
)

func getPostsForUser(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := getUserHelper(cfg, w, r)
		if user == nil {
			fmt.Println("Error getting user")
			return
		}

		ctx := r.Context()
		posts, err := cfg.DB.GetPostsForUser(ctx, database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  10,
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(posts)
	}
}
