package controller

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreatedResponse struct {
	URL string
}

func respondCreated(w http.ResponseWriter, url string) {
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreatedResponse{URL: url})
}

func respondError(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(status)

	errResponse := &ErrorResponse{Error: reason}

	_ = json.NewEncoder(w).Encode(errResponse)
}
