package web

import (
	"encoding/json"
	"net/http"
)

const (
	ErrTypeBadRequest    = "https://tools.ietf.org/html/rfc7231#section-6.5.1"
	ErrTypeNotAuthorized = "https://tools.ietf.org/html/rfc7235#section-3.1"
)

type ErrorResponse struct {
	ErrorType string       `json:"type"`
	Message   string       `json:"message"`
	Errors    []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func BadRequest(w http.ResponseWriter, errors map[string]string) {
	w.WriteHeader(http.StatusBadRequest)

	resp := ErrorResponse{
		ErrorType: ErrTypeBadRequest,
		Message:   "Bad Request",
		Errors:    []FieldError{},
	}

	for field, message := range errors {
		resp.Errors = append(resp.Errors, FieldError{Field: field, Message: message})
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func NotAuthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)

	resp := ErrorResponse{
		ErrorType: ErrTypeNotAuthorized,
		Message:   message,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
