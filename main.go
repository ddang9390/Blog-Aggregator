package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	//Get port from env file
	godotenv.Load(".env")
	port := os.Getenv("PORT")

	router := mux.NewRouter()

	//Handlers
	router.HandleFunc("/", handlePage)
	router.HandleFunc("/v1/healthz", readyHandler()).Methods("GET")
	router.HandleFunc("/v1/err", errorHandler()).Methods("GET")

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
