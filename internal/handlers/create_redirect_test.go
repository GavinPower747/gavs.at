package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"gavs.at/shortener/internal/handlers"
	"gavs.at/shortener/internal/model"
)

type MockStorage struct {
	UpsertEntityFunc func(entity interface{}) error
}

func (m *MockStorage) QueryEntity(partitionKey, rowKey string) ([]byte, error) {
	return nil, nil
}

func (m *MockStorage) UpsertEntity(entity interface{}) error {
	if m.UpsertEntityFunc == nil {
		return nil
	}

	return m.UpsertEntityFunc(entity)
}

func TestHandlers_UpsertRedirect_ValidRequest(t *testing.T) {
	// Arrange
	mockStorage := &MockStorage{
		UpsertEntityFunc: func(entity interface{}) error {
			return nil
		},
	}
	h := handlers.NewHandlers(mockStorage)
	reqBody := model.UpsertRedirectRequest{
		Slug:    "test",
		FullURL: "https://example.com",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/redirects", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	// Act
	h.UpsertRedirect(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestHandlers_UpsertRedirect_InvalidRequest(t *testing.T) {
	// Arrange
	mockStorage := &MockStorage{}
	h := handlers.NewHandlers(mockStorage)
	reqBody := model.UpsertRedirectRequest{
		Slug:    "",
		FullURL: "",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/redirects", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	// Act
	h.UpsertRedirect(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandlers_UpsertRedirect_InvalidRequestBody(t *testing.T) {
	// Arrange
	mockStorage := &MockStorage{}
	h := handlers.NewHandlers(mockStorage)
	reqBody := "invalid json"
	req, _ := http.NewRequest("POST", "/redirects", bytes.NewReader([]byte(reqBody)))
	rr := httptest.NewRecorder()

	// Act
	h.UpsertRedirect(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandlers_UpsertRedirect_StorageError(t *testing.T) {
	// Arrange
	mockStorage := &MockStorage{
		UpsertEntityFunc: func(entity interface{}) error {
			return errors.New("test error")
		},
	}
	h := handlers.NewHandlers(mockStorage)
	reqBody := model.UpsertRedirectRequest{
		Slug:    "test",
		FullURL: "https://example.com",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/redirects", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	// Act
	h.UpsertRedirect(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandlers_UpsertRedirect_InvalidRequestBodyFields(t *testing.T) {
	// Arrange
	mockStorage := &MockStorage{}
	h := handlers.NewHandlers(mockStorage)
	reqBody := model.UpsertRedirectRequest{
		Slug:    "test",
		FullURL: "invalid url",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/redirects", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	// Act
	h.UpsertRedirect(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandlers_UpsertRedirect_ReadRequestBodyError(t *testing.T) {
	// Arrange
	mockStorage := &MockStorage{}
	h := handlers.NewHandlers(mockStorage)
	req, _ := http.NewRequest("POST", "/redirects", &errorReader{})
	rr := httptest.NewRecorder()

	// Act
	h.UpsertRedirect(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
