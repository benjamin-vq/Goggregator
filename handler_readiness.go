package main

import (
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	log.Println("Request received on readiness handler")
	ok := struct {
		Status string
	}{Status: "OK"}

	respondWithJSON(w, http.StatusOK, ok)
}
