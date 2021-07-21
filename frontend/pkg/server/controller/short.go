package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type ShortService interface {
	Short(ctx context.Context, urlToShort string, short string) (string, error)
	CreateShortURL(r *http.Request, short string) string
}

type Short struct {
	ShortService ShortService
}

func NewShort(service ShortService) *Short {
	return &Short{
		ShortService: service,
	}
}

type CreateShortRequest struct {
	URL   string `json:"url"`
	Short string `json:"short"`
}

func (r *CreateShortRequest) Validate() error {
	if r.URL == "" {
		return errors.New("url field is missing")
	}

	_, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return errors.New("url field is malformed")
	}

	return nil
}

func (s *Short) CreateShort(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json") // move to middleware

	req := &CreateShortRequest{}
	json.NewDecoder(r.Body).Decode(req)

	err := req.Validate()
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	short, err := s.ShortService.Short(context.Background(), req.URL, req.Short)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Error: , %v", err))
		return
	}

	sURL := s.ShortService.CreateShortURL(r, short)
	respondCreated(w, sURL)
}
