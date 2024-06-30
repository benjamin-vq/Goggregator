package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const (
	readinessPath = "/v1/healthz"
	errorPath     = "/v1/err"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load environment variables: %q", err)
	}
	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc(readinessPath, handlerReadiness)
	mux.HandleFunc(errorPath, handlerError)

	log.Printf("Starting server on port :%s", port)
	log.Fatal(server.ListenAndServe())
}
