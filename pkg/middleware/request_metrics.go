package middleware

import (
	"log"
	"net/http"
)

func RequestMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("RequestMetricsMiddleware")
		next.ServeHTTP(w, r)
	})
}
