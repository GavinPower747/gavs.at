package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gavs.at/shortener/internal/handlers"
	"gavs.at/shortener/internal/model"
)

const (
	slug = "test"
)

type MockStorageAccount struct {
	mock.Mock
}

func (m *MockStorageAccount) QueryEntity(partitionKey, rowKey string) ([]byte, error) {
	args := m.Called(partitionKey, rowKey)

	if args.Get(0) == nil {
		return []byte{}, args.Error(1)
	}

	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorageAccount) UpsertEntity(entity interface{}) error {
	args := m.Called(entity)

	return args.Error(0)
}

func TestRedirect(t *testing.T) {
	mockStorage := &MockStorageAccount{}
	reqHandlers := handlers.NewHandlers(mockStorage)

	link := &model.Redirect{Slug: slug, FullURL: "https://example.com"}
	linkBytes, _ := json.Marshal(link)
	mockStorage.On("QueryEntity", "pk001", slug).Return(linkBytes, nil)

	r := mux.NewRouter()
	r.HandleFunc("/{slug}", reqHandlers.Redirect)

	req, _ := http.NewRequest("GET", "/"+slug, http.NoBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, link.FullURL, rr.Header().Get("Location"))
}

func TestRedirectNotFound(t *testing.T) {
	mockStorage := &MockStorageAccount{}
	reqHandlers := handlers.NewHandlers(mockStorage)

	mockStorage.On("QueryEntity", "pk001", slug).Return(nil, nil)

	r := mux.NewRouter()
	r.HandleFunc("/{slug}", reqHandlers.Redirect)

	req, _ := http.NewRequest("GET", "/"+slug, http.NoBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRedirectError(t *testing.T) {
	mockStorage := &MockStorageAccount{}
	reqHandlers := handlers.NewHandlers(mockStorage)

	mockStorage.On("QueryEntity", "pk001", slug).Return(nil, errors.New("test error"))

	r := mux.NewRouter()
	r.HandleFunc("/{slug}", reqHandlers.Redirect)

	req, _ := http.NewRequest("GET", "/"+slug, http.NoBody)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
