package main

import (
	"encoding/json"
	"github.com/benjamin-vq/goggregator/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type userParams struct {
	Name string `json:"name"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {

	params := userParams{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Could not decode request for user creation: %q", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	if params.Name == "" {
		log.Printf("Received name was empty, responding with error")
		respondWithError(w, http.StatusBadRequest, "Name can not be empty")
		return
	}

	db := cfg.DB

	dbUser, err := db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Createdat: time.Now(),
		Updatedat: time.Now(),
		Name:      params.Name,
	})
	if err != nil {
		log.Printf("Could not create user in database: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, mapDbUserResponse(&dbUser))
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, _ *http.Request, dbUser *database.User) {
	log.Println("Received a request to retrieve a user by API Key")

	respondWithJSON(w, http.StatusOK, mapDbUserResponse(dbUser))
}
