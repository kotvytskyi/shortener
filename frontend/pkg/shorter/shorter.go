package shorter

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type URLRepository interface {
	Create(ctx context.Context, url string, short string) error
	GetURL(ctx context.Context, short string) (string, error)
}

type ShortAPI interface {
	Get(ctx context.Context) (string, error)
}

type Shorter struct {
	Repository URLRepository
	API        ShortAPI
}

func NewShorter(r URLRepository, a ShortAPI) *Shorter {
	return &Shorter{Repository: r, API: a}
}

func (s *Shorter) Short(ctx context.Context, from string, to string) (shortened string, err error) {
	if to == "" {
		const timeout = time.Second * 2

		ctx, cancel := context.WithTimeout(ctx, timeout)
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

func (s *Shorter) GetURL(ctx context.Context, short string) (url string, err error) {
	return s.Repository.GetURL(ctx, short)
}

func (s *Shorter) CreateShortURL(r *http.Request, short string) string {
	return fmt.Sprintf("http://%s/short/%s", r.Host, short)
}
