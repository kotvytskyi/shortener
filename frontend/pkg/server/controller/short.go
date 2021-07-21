package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type ShortService interface {
	Short(ctx context.Context, urlToShort string, short string) (string, error)
	GetUrl(ctx context.Context, short string) (string, error)
	CreateShortURL(r *http.Request, short string) string
}

type Short struct {
	service ShortService
}

func NewShort(service ShortService) *Short {
	return &Short{
		service: service,
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
	w.Header().Add("Content-Type", "application/json")

	req := &CreateShortRequest{}
	json.NewDecoder(r.Body).Decode(req)

	err := req.Validate()
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	short, err := s.service.Short(context.Background(), req.URL, req.Short)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Error: , %v", err))
		return
	}

	sURL := s.service.CreateShortURL(r, short)
	respondCreated(w, sURL)
}

func (s *Short) ProxyShort(w http.ResponseWriter, r *http.Request) {
	short := mux.Vars(r)["short"]
	if short == "" {
		respondError(w, http.StatusNotFound, "short is empty")
		return
	}

	url, err := s.service.GetUrl(context.Background(), short)
	if err != nil {
		log.Print(fmt.Sprintf("ERROR: %v", err))
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
