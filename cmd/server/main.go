package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"gavs.at/shortener/internal/handlers"
	"gavs.at/shortener/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	listenAddr := ":8080"

	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	storageAccount, err := storage.NewStorageAccount()

	if err != nil {
		log.Fatal(err)
	}

	reqHandlers := handlers.NewHandlers(storageAccount)

	r := mux.NewRouter()
	r.HandleFunc("/{slug}", reqHandlers.Redirect)

	const timeoutDuration = 5 * time.Second

	srv := &http.Server{
		Handler: r,
		Addr:    listenAddr,

		WriteTimeout: timeoutDuration,
		ReadTimeout:  timeoutDuration,
	}

	log.Printf("About to listen on %s. Go to https://127.0.0.1%s", listenAddr, listenAddr)
	log.Fatal(srv.ListenAndServe())
}
