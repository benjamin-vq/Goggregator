package main

import (
	"encoding/json"
	"net/http"
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
