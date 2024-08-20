package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"gavs.at/shortener/internal/handlers"
	"gavs.at/shortener/internal/storage"
	"gavs.at/shortener/pkg/middleware"
	"gavs.at/shortener/pkg/observability"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := observability.SetupOTelSDK(ctx)
	if err != nil {
		log.Fatalf("Failed to setup otel: %v", err)
		return
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	listenAddr := ":80"

	storageAccount, err := storage.NewStorageAccount()
	if err != nil {
		log.Fatal(err)
	}

	reqHandlers := handlers.NewHandlers(storageAccount)

	r := configureRouter(reqHandlers)

	const timeoutDuration = 5 * time.Second

	srv := &http.Server{
		Handler:      r,
		Addr:         listenAddr,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		WriteTimeout: timeoutDuration,
		ReadTimeout:  timeoutDuration,
	}

	srvErr := make(chan error, 1)
	go func() {
		log.Println("Listening on", listenAddr)
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		stop()
	}

	err = srv.Shutdown(context.Background())
	return
}

func configureRouter(reqHandlers *handlers.Handlers) *mux.Router {
	r := mux.NewRouter()

	r.Use(otelmux.Middleware("gavs.at"))
	r.Use(middleware.RequestMetrics)

	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.Use(middleware.BasicAuth)

	apiRouter.HandleFunc("/redirect", reqHandlers.UpsertRedirect).Methods("POST")

	r.HandleFunc("/{slug}", reqHandlers.Redirect).Methods("GET")

	return r
}
