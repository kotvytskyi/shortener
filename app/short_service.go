package app

import (
	"context"
	"net/http"
)

type ShortRepository interface {
	Create(ctx context.Context, url string, short string) error
}

type AppShortService struct {
	Repository ShortRepository
}

func (s *AppShortService) Create(ctx context.Context, urlToShort string, short string) error {
	// validate url

	err := s.Repository.Create(ctx, urlToShort, short)
	return err
}

func (s *AppShortService) CreateShortURL(r *http.Request, short string) string {
	// solve tests https issues
	return "http://" + r.Host + "/short/" + short
}
