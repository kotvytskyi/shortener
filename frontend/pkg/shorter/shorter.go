package shorter

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type UrlRepository interface {
	Create(ctx context.Context, url string, short string) error
}

type ShortApi interface {
	Get(ctx context.Context) (string, error)
}

type Shorter struct {
	Repository UrlRepository
	API        ShortApi
}

func NewShorter(r UrlRepository, a ShortApi) *Shorter {
	s := &Shorter{}
	s.Repository = r
	s.API = a
	return s
}

func (s *Shorter) Short(ctx context.Context, from string, to string) (shortened string, err error) {
	if to == "" {
		ctx, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()

		s, err := s.API.Get(ctx)
		if err != nil {
			return "", err
		}

		to = s
	}

	err = s.Repository.Create(ctx, from, to)
	return to, err
}

func (s *Shorter) CreateShortURL(r *http.Request, short string) string {
	return fmt.Sprintf("http://%s/short/%s", r.Host, short)
}
