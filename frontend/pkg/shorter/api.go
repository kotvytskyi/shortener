package shorter

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Key struct {
	Value string `json:"val"`
}

type AppShortApi struct {
	base *url.URL
}

func NewShortApi(baseUrl string) (*AppShortApi, error) {
	addr, err := url.Parse(baseUrl)
	if err != nil {
		return nil, errors.New("SHORTSRV is not a valid address")
	}
	return &AppShortApi{addr}, nil
}

func (c *AppShortApi) Get(ctx context.Context) (string, error) {
	req, err := http.NewRequest("POST", "http://"+c.base.String()+"/api/keys", bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
