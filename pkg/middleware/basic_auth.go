package middleware

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gavs.at/shortener/pkg/web"
	"go.opentelemetry.io/otel"
)

const (
	usernameEnvVar     = "API_BASIC_AUTH_USERNAME"
	passwordHashEnvVar = "API_BASIC_AUTH_PASSWORD_HASH" //nolint:gosec // It's just the password hash that we're storing, not the password itself.
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.GetTracerProvider().Tracer("gavs.at/shortener")
		_, span := tracer.Start(r.Context(), "BasicAuth")

		endSpan := func() {
			if span.IsRecording() {
				span.End()
			}
		}

		defer endSpan()

		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))

		if authHeader == "" {
			web.NotAuthorized(w, "Missing Authorization Header")
			span.RecordError(fmt.Errorf("Missing Authorization Header"))

			return
		}

		headerSections := strings.Split(authHeader, " ")

		if !strings.HasPrefix(authHeader, "Basic ") {
			web.NotAuthorized(w, fmt.Sprintf("invalid authorization header, %s authentication scheme is not supported", headerSections[0]))
			span.RecordError(fmt.Errorf("invalid authorization header, %s authentication scheme is not supported", headerSections[0]))

			return
		}

		decodedCreds, err := base64.StdEncoding.DecodeString(headerSections[1])

		if err != nil {
			web.NotAuthorized(w, "Invalid Authorization Header")
			span.RecordError(fmt.Errorf("invalid authorization header"))

			return
		}

		creds := strings.Split(string(decodedCreds), ":")
		username := creds[0]
		password := creds[1]

		passwordHash := getPasswordHash(password)

		expectedUsername := os.Getenv(usernameEnvVar)
		expectedPasswordHash := os.Getenv(passwordHashEnvVar)

		if username != expectedUsername || passwordHash != expectedPasswordHash {
			web.NotAuthorized(w, "Invalid Credentials")
			span.RecordError(fmt.Errorf("invalid credentials"))

			return
		}

		endSpan()

		next.ServeHTTP(w, r)
	})
}

func getPasswordHash(password string) string {
	hasher := sha256.New()

	hasher.Write([]byte(password))
	bytes := hasher.Sum(nil)

	return base64.StdEncoding.EncodeToString(bytes)
}
