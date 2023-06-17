package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gavs.at/shortener/internal/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorageAccount struct {
	mock.Mock
}

func (m *MockStorageAccount) QueryEntity(partitionKey string, rowKey string) ([]byte, error) {
    args := m.Called(partitionKey, rowKey)

    if args.Get(0) == nil {
        return []byte{}, args.Error(1)
    }
	
    return args.Get(0).([]byte), args.Error(1)
}

func TestRedirect(t *testing.T) {
    mockStorage := &MockStorageAccount{}
    handlers := &Handlers{storage: mockStorage}

    slug := "test"
    link := &model.Link{Slug: slug, FullUrl: "https://example.com"}
    linkBytes, _ := json.Marshal(link)
    mockStorage.On("QueryEntity", "1", slug).Return(linkBytes, nil)

    r := mux.NewRouter()
    r.HandleFunc("/{slug}", handlers.Redirect)

    req, _ := http.NewRequest("GET", "/"+slug, nil)

    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusFound, rr.Code)
    assert.Equal(t, link.FullUrl, rr.Header().Get("Location"))
}

func TestRedirectNotFound(t *testing.T) {
    mockStorage := &MockStorageAccount{}
    handlers := &Handlers{storage: mockStorage}

    slug := "test"
    mockStorage.On("QueryEntity", "1", slug).Return(nil, nil)

    r := mux.NewRouter()
    r.HandleFunc("/{slug}", handlers.Redirect)

    req, _ := http.NewRequest("GET", "/"+slug, nil)

    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRedirectError(t *testing.T) {
    mockStorage := &MockStorageAccount{}
    handlers := &Handlers{storage: mockStorage}

    slug := "test"
    mockStorage.On("QueryEntity", "1", slug).Return(nil, errors.New("test error"))

    r := mux.NewRouter()
    r.HandleFunc("/{slug}", handlers.Redirect)

    req, _ := http.NewRequest("GET", "/"+slug, nil)

    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
