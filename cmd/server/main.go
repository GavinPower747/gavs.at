package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"

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

const (
	ServiceName = "gavs.at"
)

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := observability.SetupOTelSDK(ctx, ServiceName)
	if err != nil {
		log.Printf("Failed to setup otel: %v", err)
		return err
	}

	tracer := otel.GetTracerProvider().Tracer(ServiceName)

	ctx, span := tracer.Start(ctx, "ApplicationStartup")

	defer func() {
		if span.IsRecording() && err != nil {
			span.RecordError(err)
			span.End()
		}

		err = errors.Join(err, otelShutdown(ctx))
	}()

	listenAddr := ":80"

	storageAccount, err := storage.NewStorageAccount(ctx)
	if err != nil {
		log.Println(err)

		return err
	}

	reqHandlers := handlers.NewHandlers(storageAccount)

	r := configureRouter(reqHandlers)

	const timeoutDuration = 5 * time.Second

	srv := &http.Server{
		Handler:      r,
		Addr:         listenAddr,
		WriteTimeout: timeoutDuration,
		ReadTimeout:  timeoutDuration,
	}

	srvErr := make(chan error, 1)

	span.End()

	go func() {
		log.Println("Listening on", listenAddr)
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		if span.IsRecording() {
			span.RecordError(err)
			span.End()
		}

		return err
	case <-ctx.Done():
		stop()
	}

	err = srv.Shutdown(ctx)

	return nil
}

func configureRouter(reqHandlers *handlers.Handlers) *mux.Router {
	r := mux.NewRouter()

	r.Use(otelmux.Middleware(ServiceName))
	r.Use(middleware.RequestMetrics)

	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.Use(middleware.BasicAuth)

	apiRouter.HandleFunc("/redirect", reqHandlers.UpsertRedirect).Methods("POST")

	r.HandleFunc("/{slug}", reqHandlers.Redirect).Methods("GET")

	return r
}
