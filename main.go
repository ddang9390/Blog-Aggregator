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
	router.HandleFunc("/v1/healthz", readyHandler()).Methods("GET")
	router.HandleFunc("/v1/err", errorHandler()).Methods("GET")

	//Keep server running
	http.Handle("/", router)
	http.ListenAndServe(":"+port, router)
}
