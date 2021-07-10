package server

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApiKeys(t *testing.T) {
	ds := &MockedDataService{}
	ds.On("ReserveKey", mock.Anything).Return("test", nil).Times(1)

	httpServer := &HttpServer{DataService: ds}
	server := httptest.NewServer(httpServer.router())
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/keys", "application/json", bytes.NewBuffer([]byte{}))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Cannot decode the body %v", err)
	}

	key := string(bodyBytes)
	assert.NotEmpty(t, key)
}

type MockedDataService struct {
	mock.Mock
}

func (m *MockedDataService) ReserveKey(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
