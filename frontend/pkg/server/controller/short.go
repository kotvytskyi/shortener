package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kotvytskyi/frontend/pkg/config"
	"github.com/kotvytskyi/frontend/pkg/mongo"
	"github.com/kotvytskyi/frontend/pkg/shorter"
)

type Short struct {
	ShortService shortService
}

func NewShort(mongoCfg config.MongoConfig, shortCfg config.ShortServerConfig) (*Short, error) {
	repository, err := mongo.NewShort(mongoCfg)
	if err != nil {
		return nil, err
	}

	api, err := shorter.NewShortApi(shortCfg)
	if err != nil {
		return nil, err
	}

	return &Short{
		ShortService: shorter.NewShorter(repository, api),
	}, nil
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

type shortService interface {
	Short(ctx context.Context, urlToShort string, short string) (string, error)
	CreateShortURL(r *http.Request, short string) string
}
