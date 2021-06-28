package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApiKeys(t *testing.T) {
	ds := &MockedDataService{}
	ds.On("ReserveKey", mock.Anything).Return(Key{Value: "test"}, nil).Times(1)

	httpServer := &HttpServer{DataService: ds}
	server := httptest.NewServer(httpServer.router())
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/keys", "application/json", bytes.NewBuffer([]byte{}))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	key := &Key{}
	json.NewDecoder(resp.Body).Decode(key)
	assert.NotNil(t, key.Value)
}

type MockedDataService struct {
	mock.Mock
}

func (m *MockedDataService) ReserveKey(ctx context.Context) (*Key, error) {
	args := m.Called(ctx)
	key := args.Get(0).(Key)
	err := args.Error(1)
	return &key, err
}
