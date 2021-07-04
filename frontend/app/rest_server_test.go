package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/frontend/app/mongo"
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
	godotenv.Load()

	cases := []TestCase{
		{Url: "", Expected: ExpectedResult{Status: http.StatusBadRequest, Error: "url field is missing."}},
		{Url: "i'm not a url", Expected: ExpectedResult{Status: http.StatusBadRequest, Error: "url field is malformed."}},
		{
			Url:      "http://www.google.com",
			Short:    "short_google_url",
			Expected: ExpectedResult{Status: http.StatusCreated, Url: "/short/short_google_url"},
		},
	}

	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			server, teardown := CreateTestServer(t)
			defer teardown()

			httpServer := httptest.NewServer(server.router())

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

			url, err := url.Parse(short.Url)
			assert.Nil(t, err)

			assert.Equal(t, test.Expected.Url, url.Path)
			assert.NotNil(t, short)
		})
	}

	t.Run("generates random path when short is empty", func(t *testing.T) {
		restServer, teardown := CreateTestServer(t)
		defer teardown()

		httpServer := httptest.NewServer(restServer.router())

		request := CreateShortRequest{
			URL: "http://www.google.com",
		}
		rBytes, _ := json.Marshal(request)

		resp, err := http.Post(httpServer.URL+"/api/shorts", "application/json", bytes.NewBuffer(rBytes))
		assert.Nil(t, err)

		short := &CreatedResponse{}
		json.NewDecoder(resp.Body).Decode(short)

		matched, err := regexp.MatchString("/short/.+", short.Url)
		assert.Nil(t, err)
		assert.True(t, matched)
	})
}

func CreateTestServer(t *testing.T) (*RestServer, func()) {
	c, teardown := testutils.CreateTestMongoConnection(t)

	api, err := NewShortApi(os.Getenv("SHORTSRV"))
	if err != nil {
		panic(err)
	}
	r := &mongo.Short{Coll: c}

	s := NewShorter(r, api)
	restServer := &RestServer{ShortService: s}

	return restServer, teardown
}
