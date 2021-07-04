package server

import (
	"errors"
	"net/url"
)

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
