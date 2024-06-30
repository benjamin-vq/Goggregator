package main

import (
	"log"
	"net/http"
)

func handlerError(w http.ResponseWriter, _ *http.Request) {
	log.Println("Request received on error handler")
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
