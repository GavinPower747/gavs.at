package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHandler struct {
	Called bool
}

func (m *MockHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
	m.Called = true
}

func TestBasicAuth_MissingAuthorizationHeader(t *testing.T) {
	mockHandler := &MockHandler{}
	authHandler := BasicAuth(mockHandler)

	req, _ := http.NewRequest("GET", "/", http.NoBody)

	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
	assert.False(t, mockHandler.Called)
}

func TestBasicAuth_InvalidAuthorizationHeader(t *testing.T) {
	mockHandler := &MockHandler{}
	authHandler := BasicAuth(mockHandler)

	req, _ := http.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("Authorization", "Bearer token")

	rr := httptest.NewRecorder()

	authHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
	assert.False(t, mockHandler.Called)
}

func TestBasicAuth_InvalidBase64Encoding(t *testing.T) {
	mockHandler := &MockHandler{}
	authHandler := BasicAuth(mockHandler)

	req, _ := http.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("Authorization", "Basic invalid")

	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
	assert.False(t, mockHandler.Called)
}

func TestBasicAuth_InvalidCredentials(t *testing.T) {
	testCases := []struct {
		caseName string
		username string
		password string
	}{
		{"Both_Wrong", "wrong", "wrong"},
		{"Password_Wrong", "admin", "wrong"},
		{"Username_Wrong", "wrong", "admin"},
	}

	mockHandler := &MockHandler{}
	authHandler := BasicAuth(mockHandler)

	t.Setenv("API_BASIC_AUTH_USERNAME", "admin")
	t.Setenv("API_BASIC_AUTH_PASSWORD_HASH", getPasswordHash("admin"))

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", http.NoBody)
			credentials := fmt.Sprintf("%s:%s", tc.username, tc.password)
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentials))))

			rr := httptest.NewRecorder()
			authHandler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
			assert.NotEmpty(t, rr.Body.String())
			assert.False(t, mockHandler.Called)
		})
	}
}

func TestBasicAuth_ValidCredentials(t *testing.T) {
	mockHandler := &MockHandler{}
	authHandler := BasicAuth(mockHandler)

	t.Setenv("API_BASIC_AUTH_USERNAME", "admin")
	t.Setenv("API_BASIC_AUTH_PASSWORD_HASH", getPasswordHash("admin"))

	req, _ := http.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("admin:admin"))))

	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, mockHandler.Called)
}
