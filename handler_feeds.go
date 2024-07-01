package main

import (
	"encoding/json"
	"github.com/benjamin-vq/goggregator/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type feedsParams struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, dbUser *database.User) {
	log.Println("Received a request to create a feed")

	params := feedsParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding create feed parameters: %q", err)
		respondWithError(w, http.StatusBadRequest, "Could not decode request")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    dbUser.ID,
		Createdat: time.Now(),
		Updatedat: time.Now(),
	})
	if err != nil {
		log.Printf("Could not create feed in database: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create feed")
		return
	}
	ff, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    dbUser.ID,
		FeedID:    feed.ID,
		Createdat: time.Now(),
		Updatedat: time.Now(),
	})
	if err != nil {
		log.Printf("Could not create feed follow after creating feed: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, mapCreateFeedResponse(&feed, &ff))
}

func (cfg *apiConfig) handlerFeedsGetAll(w http.ResponseWriter, r *http.Request) {

	dbFeeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		log.Printf("Error retrieving all feeds from database: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve feeds")
		return
	}

	var feeds []FeedResponse
	for _, feed := range dbFeeds {
		feeds = append(feeds, mapDbFeedResponse(&feed))
	}

	log.Printf("Returning a total of %d feeds", len(feeds))
	respondWithJSON(w, http.StatusOK, feeds)
}

type feedFollowParams struct {
	FeedId string `json:"feed_id"`
}

func (cfg *apiConfig) handlerFeedFollowCreate(w http.ResponseWriter, r *http.Request, dbUser *database.User) {

	params := feedFollowParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Could not decode params to create a feed follow: %q", err)
		respondWithError(w, http.StatusBadRequest, "Could not decode parameters")
		return
	}

	parsedUuid, err := uuid.Parse(params.FeedId)
	if err != nil {
		log.Printf("Could not parse received feed id: %q", err)
		respondWithError(w, http.StatusBadRequest, "Invalid feed id")
		return
	}

	ff, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    dbUser.ID,
		FeedID:    parsedUuid,
		Createdat: time.Now(),
		Updatedat: time.Now(),
	})
	if err != nil {
		log.Printf("Could not create feed follow: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create feed follow")
		return
	}

	log.Printf("Succesfully created feed follow for user id %s and feed id %s",
		dbUser.ID.String(), parsedUuid.String())
	respondWithJSON(w, http.StatusCreated, mapDbFeedFollowResponse(&ff))
}

func (cfg *apiConfig) handlerFeedFollowDelete(w http.ResponseWriter, r *http.Request, dbUser *database.User) {

	ffID := r.PathValue("feedFollowID")
	userId := dbUser.ID

	parsedFfID, err := uuid.Parse(ffID)
	if err != nil {
		log.Printf("Could not parse feed follow ID from URL: %q", err)
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow")
		return
	}

	deletedFf, err := cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: parsedFfID,
		UserID: userId,
	})
	if err != nil {
		log.Printf("Could not delete feed follow: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not delete feed follow")
		return
	}

	log.Printf("Successfully deleted feed follow %s for user id %s",
		parsedFfID, userId)
	respondWithJSON(w, http.StatusOK, mapDbFeedFollowResponse(&deletedFf))
}

func (cfg *apiConfig) handlerFFByUser(w http.ResponseWriter, r *http.Request, dbUser *database.User) {

	ffs, err := cfg.DB.GetFeedFollowsByUser(r.Context(), dbUser.ID)
	if err != nil {
		log.Printf("Could not geet feed follows for user id %s: %q",
			dbUser.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Could not get feed follows")
		return
	}

	var ffResponse []FeedFollowResponse
	for _, ff := range ffs {
		ffResponse = append(ffResponse, mapDbFeedFollowResponse(&ff))
	}

	log.Printf("Returning %d feed follows for user id %s", len(ffResponse), dbUser.ID)
	respondWithJSON(w, http.StatusOK, ffResponse)
}
