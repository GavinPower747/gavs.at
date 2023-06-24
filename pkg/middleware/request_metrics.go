package middleware

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"

	"gavs.at/shortener/pkg/web"
)

func RequestMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		recorder := &web.StatusRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(startTime)

		connectionString := os.Getenv("APPLICATIONINSIGHTS_CONNECTION_STRING")

		if connectionString == "" {
			return
		}

		client := appinsights.NewTelemetryClient(connectionString)
		trace := appinsights.NewRequestTelemetry(r.Method, r.URL.String(), duration, strconv.Itoa(recorder.StatusCode))

		client.Track(trace)
	})
}
