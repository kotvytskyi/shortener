package app

import (
	"context"
	"net/http"
)

type UrlRepository interface {
	Create(ctx context.Context, url string, short string) error
}

type ShortClient interface {
	Generate(ctx context.Context) (string, error)
}

type AppShortService struct {
	Repository UrlRepository
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
