package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kotvytskyi/shortener/testutils"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Url      string
	Short    string
	Expected ExpectedResult
}

type ExpectedResult struct {
	Error  string
	Status int
	Url    string
	Short  string
}

func TestApiShort(t *testing.T) {
	testCases := []TestCase{
		{Url: "", Expected: ExpectedResult{Status: http.StatusBadRequest, Error: "url field is missing."}},
		{Url: "i'm not a url", Expected: ExpectedResult{Status: http.StatusBadRequest, Error: "url field is malformed."}},
		{
			Url:      "http://www.google.com",
			Short:    "short_google_url",
			Expected: ExpectedResult{Status: http.StatusCreated, Url: "short_google_url"},
		},
	}

	for _, test := range testCases {
		t.Run("", func(t *testing.T) {
			restServer, teardown := CreateTestServer(t)
			defer teardown()

			httpServer := httptest.NewServer(restServer.router())

			request := CreateShortRequest{
				URL:   test.Url,
				Short: test.Short,
			}
			rBytes, _ := json.Marshal(request)

			resp, err := http.Post(httpServer.URL+"/api/shorts", "application/json", bytes.NewBuffer(rBytes))
			assert.Nil(t, err)
			assert.Equal(t, test.Expected.Status, resp.StatusCode)

			if resp.StatusCode >= 400 {
				eResp := &ErrorResponse{}
				json.NewDecoder(resp.Body).Decode(eResp)
				assert.Equal(t, eResp.Error, test.Expected.Error)
				return
			}

			short := &CreatedResponse{}
			json.NewDecoder(resp.Body).Decode(short)

			assert.Equal(t, ShortUrl(httpServer, test.Expected.Url), short.Url)
			assert.NotNil(t, short)
		})
	}
}

func CreateTestServer(t *testing.T) (*RestServer, func()) {
	coll, teardown := testutils.CreateTestMongoConnection(t)

	restServer := &RestServer{}
	service := &AppShortService{}
	repo := &MongoShortRepository{coll}
	service.Repository = repo
	restServer.ShortService = service

	return restServer, teardown
}

func ShortUrl(s *httptest.Server, url string) string {
	return s.URL + "/short/" + url
}
