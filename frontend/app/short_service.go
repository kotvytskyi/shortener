package app

import (
	"context"
	"net/http"
	"time"
)

type UrlRepository interface {
	Create(ctx context.Context, url string, short string) error
}

type ShortApi interface {
	Get(ctx context.Context) (string, error)
}

type AppShortService struct {
	Repository UrlRepository
	ShortApi   ShortApi
}

func NewShortService(r UrlRepository, a ShortApi) *AppShortService {
	s := &AppShortService{}
	s.Repository = r
	s.ShortApi = a
	return s
}

func (s *AppShortService) Create(ctx context.Context, urlToShort string, short string) (string, error) {
	if short == "" {
		ctx, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()

		s, err := s.ShortApi.Get(ctx)
		if err != nil {
			return "", err
		}

		short = s
	}

	err := s.Repository.Create(ctx, urlToShort, short)
	return short, err
}

func (s *AppShortService) CreateShortURL(r *http.Request, short string) string {
	// solve tests https issues
	return "http://" + r.Host + "/short/" + short
}
