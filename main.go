package main

import (
	"blog-aggregator/internal/database"
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
	DB *database.Queries
}

func main() {
	//Get port from env file
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
	}

	dbQueries := database.New(db)
	cfg := &apiConfig{DB: dbQueries}

	router := mux.NewRouter()

	//Handlers
	router.HandleFunc("/", handlePage)
	router.HandleFunc("/v1/healthz", readyHandler()).Methods("GET")
	router.HandleFunc("/v1/err", errorHandler()).Methods("GET")

	//User handlers
	router.HandleFunc("/v1/users", createUser(cfg)).Methods("POST")
	router.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) { getUser(cfg, w, r) }).Methods("GET")

	//Feed handlers
	router.HandleFunc("/v1/feeds", createFeed(cfg)).Methods("POST")
	router.HandleFunc("/v1/feeds", getAllFeeds(cfg)).Methods("GET")

	router.HandleFunc("/v1/feed_follows", createFeedFollow(cfg)).Methods("POST")
	router.HandleFunc("/v1/feed_follows", deleteFeedFollow(cfg)).Methods("DELETE")
	router.HandleFunc("/v1/feed_follows", getAllFeedFollowsForUser(cfg)).Methods("GET")

	//Testing worker
	limit := 10
	duration := time.Minute
	go fetchWorker(cfg, limit, duration)

	//Keep server running
	//http.Handle("/", router)
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
