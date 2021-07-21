package controller

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreatedResponse struct {
	Url string
}

func respondCreated(w http.ResponseWriter, url string) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatedResponse{Url: url})
}

func respondError(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(status)
	errResponse := &ErrorResponse{Error: reason}
	json.NewEncoder(w).Encode(errResponse)
}
