package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/frontend/pkg/mongo"
	"github.com/kotvytskyi/frontend/pkg/server/controller"
	"github.com/kotvytskyi/frontend/pkg/shorter"
	"github.com/kotvytskyi/testmongo"
	"github.com/stretchr/testify/require"
)

func TestApiShort(t *testing.T) {
	godotenv.Load()

	cases := []struct {
		name       string
		url, short string

		status   int
		response string

		apiMock string
	}{
		{name: "Empty URL", url: "", short: "", status: http.StatusBadRequest, response: "url field is missing"},
		{name: "Malformed URL", url: "notAURL", short: "", status: http.StatusBadRequest, response: "url field is malformed"},
		{name: "Custom short", url: "http://www.google.com", short: "ggl", status: http.StatusCreated, response: "/short/ggl"},
		{name: "Empty short", url: "http://www.google.com", short: "", status: http.StatusCreated, response: "/short/abc123", apiMock: "abc123"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			server, teardown := CreateTestServer(t, test.apiMock)
			defer teardown()

			srv := httptest.NewServer(server.router())
			defer srv.Close()

			request := controller.CreateShortRequest{
				URL:   test.url,
				Short: test.short,
			}
			rBytes, _ := json.Marshal(request)

			resp, err := http.Post(srv.URL+"/api/shorts", "application/json", bytes.NewBuffer(rBytes))
			require.Nil(t, err)
			require.Equal(t, test.status, resp.StatusCode)

			if resp.StatusCode >= 400 {
				eResp := &controller.ErrorResponse{}
				json.NewDecoder(resp.Body).Decode(eResp)
				require.Equal(t, eResp.Error, test.response)
				return
			}

			r := &controller.CreatedResponse{}
			json.NewDecoder(resp.Body).Decode(r)

			u, err := url.Parse(r.Url)
			require.Nil(t, err)

			require.Equal(t, test.response, u.Path)
		})
	}
}

func CreateTestServer(t *testing.T, mockedShort string) (*RestServer, func()) {
	t.Helper()

	c, teardown := testmongo.CreateTestMongoConnection(t)

	api := &MockedApi{mockedShort: mockedShort}
	r := &mongo.Short{Coll: c}

	s := shorter.NewShorter(r, api)

	controller := controller.NewShort(s)
	restServer := &RestServer{controller: controller}

	return restServer, teardown
}

type MockedApi struct {
	mockedShort string
}

func (a *MockedApi) Get(ctx context.Context) (string, error) {
	return a.mockedShort, nil
}
