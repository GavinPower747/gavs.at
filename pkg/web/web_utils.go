package web

import (
	"encoding/json"
	"net/http"
)

const (
	ErrTypeBadRequest = "https://tools.ietf.org/html/rfc7231#section-6.5.1"
)

type ErrorResponse struct {
	ErrorType string       `json:"type"`
	Message   string       `json:"message"`
	Errors    []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string
	Message string
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

	json.NewEncoder(w).Encode(resp)
}
