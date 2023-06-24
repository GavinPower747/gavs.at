package web

import (
	"encoding/json"
	"io"
	"net/http"
)

type Validatable interface {
	Validate() (bool, ValidationError)
}

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "Validation error"
}

func (e *ValidationError) AddError(field, message string) {
	e.Errors[field] = message
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

func ValidateAndParseBody[T Validatable](r *http.Request) (*T, map[string]string) {
	b, err := io.ReadAll(r.Body)

	if err != nil {
		return nil, map[string]string{"body": "Could not read request body"}
	}

	var req T
	err = json.Unmarshal(b, &req)

	if err != nil {
		return nil, map[string]string{"body": "Could not unmarshal request body"}
	}

	isValid, vErr := req.Validate()

	if !isValid {
		return nil, vErr.Errors
	}

	return &req, nil
}
