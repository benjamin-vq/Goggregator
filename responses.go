package main

import (
	"encoding/json"
	"github.com/benjamin-vq/goggregator/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	ContentTypeHeader = "Content-Type"
	ApplicationJson   = "application/json"
)

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	bytes, err := json.Marshal(&payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(ContentTypeHeader, ApplicationJson)
	w.WriteHeader(code)
	w.Write(bytes)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	assert(code > 399, "HTTP Response code is not an error one")
	errResponse := struct {
		ErrMsg string `json:"error"`
	}{ErrMsg: msg}

	respondWithJSON(w, code, errResponse)
}

type UserResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func mapDbUserResponse(user *database.User) UserResponse {
	return UserResponse{
		Id:        user.ID,
		CreatedAt: user.Createdat,
		UpdatedAt: user.Updatedat,
		Name:      user.Name,
		ApiKey:    user.Apikey,
	}
}

type CreateFeedResponse struct {
	MappedFeed FeedResponse       `json:"feed"`
	MappedFF   FeedFollowResponse `json:"feed_follow"`
}

type FeedResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func mapCreateFeedResponse(feed *database.Feed, ff *database.FeedFollow) CreateFeedResponse {
	return CreateFeedResponse{
		MappedFeed: FeedResponse{
			ID:        feed.ID,
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    feed.UserID,
			CreatedAt: feed.Createdat,
			UpdatedAt: feed.Updatedat,
		},
		MappedFF: FeedFollowResponse{
			ID:        ff.ID,
			UserID:    ff.UserID,
			FeedID:    ff.FeedID,
			CreatedAt: ff.Createdat,
			UpdatedAt: ff.Updatedat,
		},
	}
}

func mapDbFeedResponse(feed *database.Feed) FeedResponse {
	return FeedResponse{
		ID:        feed.ID,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
		CreatedAt: feed.Createdat,
		UpdatedAt: feed.Updatedat,
	}
}

type FeedFollowResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func mapDbFeedFollowResponse(ff *database.FeedFollow) FeedFollowResponse {
	return FeedFollowResponse{
		ID:        ff.ID,
		UserID:    ff.UserID,
		FeedID:    ff.FeedID,
		CreatedAt: ff.Createdat,
		UpdatedAt: ff.Updatedat,
	}
}
