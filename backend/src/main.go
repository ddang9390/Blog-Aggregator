package main

import (
	"blog-aggregator/backend/internal/database"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// Using config struct to store shared data that http handlers need access to
type apiConfig struct {
	DB        *database.Queries
	jwtSecret string
	jwtToken  string
}

func main() {
	//Get port from env file
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	fmt.Println(port)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
	}

	dbQueries := database.New(db)
	cfg := &apiConfig{DB: dbQueries, jwtSecret: jwtSecret}

	router := mux.NewRouter()

	//Handlers
	//router.HandleFunc("/", handlePage)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "blog-aggregator/frontend/login.html")
	})
	router.HandleFunc("/v1/healthz", readyHandler()).Methods("GET")
	router.HandleFunc("/v1/err", errorHandler()).Methods("GET")

	//User handlers
	// For registering a user
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.ServeFile(w, r, "../../frontend/register.html")
		} else if r.Method == "POST" {
			createUser(cfg, w, r)
			//router.HandleFunc("/v1/users", createUser(cfg)).Methods("POST")
		}
	}).Methods("GET", "POST")

	// For logging in a user
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.ServeFile(w, r, "../../frontend/login.html")
		} else if r.Method == "POST" {
			loginUser(cfg, w, r)
		}
	}).Methods("GET", "POST")

	//Feed handlers
	router.HandleFunc("/feeds", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Authorization", cfg.jwtToken)
		if r.Method == "GET" {
			getAllFeeds(cfg, w, r)
			//http.ServeFile(w, r, "../../frontend/feeds.html")

		} else if r.Method == "POST" {
			createFeed(cfg, w, r)
		}
	}).Methods("GET", "POST")

	//router.HandleFunc("/v1/feeds", createFeed(cfg)).Methods("POST")
	//router.HandleFunc("/v1/feeds", getAllFeeds(cfg)).Methods("GET")

	//Feed Follow Handlers
	router.HandleFunc("/feed_follows", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Authorization", cfg.jwtToken)
		if r.Method == "GET" {
			getAllFeedFollowsForUser(cfg, w, r)
			//http.ServeFile(w, r, "../../frontend/feeds.html")

		}
	}).Methods("GET")

	router.HandleFunc("/v1/feed_follows", createFeedFollow(cfg)).Methods("POST")
	router.HandleFunc("/v1/feed_follows", deleteFeedFollow(cfg)).Methods("DELETE")
	//router.HandleFunc("/v1/feed_follows", getAllFeedFollowsForUser(cfg)).Methods("GET")

	//Post Handlers
	router.HandleFunc("/posts/{feedID}", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Authorization", cfg.jwtToken)
		if r.Method == "GET" {
			getPosts(cfg, w, r)
		}
	}).Methods("GET")

	router.HandleFunc("/v1/posts", getPostsForUser(cfg)).Methods("GET")

	//Testing worker
	limit := 10
	duration := time.Minute
	go fetchWorker(cfg, limit, duration)

	//Keep server running
	router.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/", http.FileServer(http.Dir("../../frontend"))))
	http.ListenAndServe(":"+port, router)

}

// for testing if docker works
func handlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	const page = `<html>
<head></head>
<body>
	
	<p>Hi Docker, I pushed a new version!</p>
</body>
</html>
`
	w.Write([]byte(page))
}
