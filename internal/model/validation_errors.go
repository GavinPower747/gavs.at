package model

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "Validation error"
}

func (e *ValidationError) AddError(field string, message string) {
	e.Errors[field] = message
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}
