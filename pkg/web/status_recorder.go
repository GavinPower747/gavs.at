package web

import "net/http"

// Wrapper for http.ResponseWriter that keeps track of the status code for logging/metrics
// Satisfies the http.ResponseWriter interface so it can be passed in its place
type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (r *StatusRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
