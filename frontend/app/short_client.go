package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Key struct {
	Value string `json:"val"`
}

type AppShortClient struct {
	baseUrl url.URL
}

func NewShortClient(baseUrl url.URL) AppShortClient {
	return AppShortClient{baseUrl}
}

func (c *AppShortClient) Generate(ctx context.Context) (string, error) {
	req, err := http.NewRequest("POST", c.baseUrl.String()+"/api/keys", bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	key := &Key{}
	err = json.NewDecoder(resp.Body).Decode(key)
	if err != nil {
		return "", err
	}

	return key.Value, nil
}
