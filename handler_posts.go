package main

import (
	"github.com/benjamin-vq/goggregator/internal/database"
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, dbUser *database.User) {

	limitString := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		log.Printf("Invalid 'limit' query param: %s. Defaulting to 10 posts", limit)
		limit = 10
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		ID:    dbUser.ID,
		Limit: int32(limit),
	})
	if err != nil {
		log.Printf("Could not retrieve posts by user with id %s: %q",
			dbUser.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve posts")
		return
	}

	log.Printf("Successfully retrieved %d posts from user with id: %s", len(posts), dbUser.ID)
	respondWithJSON(w, http.StatusOK, posts)
}
