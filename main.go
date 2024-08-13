package main

import (
	"blog-aggregator/internal/database"
	"database/sql"
	"fmt"
	"net/http"
	"os"

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
	router.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createUser(cfg)
		default:
			getUser(cfg, w, r)
		}
	})

	//Feed handlers
	// router.HandleFunc("/v1/feeds", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		createFeed(cfg)
	// 	default:
	// 		createFeed(cfg)
	// 	}
	// })
	router.HandleFunc("/v1/feeds", createFeed(cfg)).Methods("POST")
	//router.Get("/users", cfg.middlewareAuth(cfg.handlerUsersGet))

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
