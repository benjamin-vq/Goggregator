package main

import (
	"github.com/benjamin-vq/goggregator/internal/database"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthedHandler func(http.ResponseWriter, *http.Request, *database.User)

func (cfg *apiConfig) authMiddleware(authedHandler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		before := time.Now()
		defer func() {
			duration := time.Since(before).Milliseconds()
			log.Printf("Request finished after %d milliseconds", duration)
		}()

		authHeader := r.Header.Get("Authorization")
		apiKey, found := strings.CutPrefix(authHeader, "Bearer ")
		if apiKey == "" || !found {
			log.Printf("Incorrect authorization header: %s", authHeader)
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		dbUser, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			log.Printf("Could not retrieve user from database: %q", err)
			respondWithError(w, http.StatusNotFound, "Could not retrieve user")
			return
		}

		authedHandler(w, r, &dbUser)
	}
}
