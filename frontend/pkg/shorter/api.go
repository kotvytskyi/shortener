package shorter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/kotvytskyi/frontend/pkg/config"
)

type apiResponse struct {
	Key string `json:"key"`
}

type AppShortAPI struct {
	base *url.URL
}

func NewShortAPI(cfg config.ShortServerConfig) (*AppShortAPI, error) {
	addr, err := url.Parse(cfg.Address)
	if err != nil {
		return nil, errors.New("SHORTSRV is not a valid address")
	}

	return &AppShortAPI{addr}, nil
}

func (c *AppShortAPI) Get(ctx context.Context) (string, error) {
	req, err := http.NewRequest("POST", "http://"+c.base.String()+"/api/keys", bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("cannot obtain a new short")
	}

	key := &apiResponse{}

	err = json.NewDecoder(resp.Body).Decode(key)
	if err != nil {
		return "", err
	}

	return key.Key, nil
}
