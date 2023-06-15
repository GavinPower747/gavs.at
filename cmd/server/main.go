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

	storage_account, err := storage.NewStorageAccount()

	if err != nil {
		log.Fatal(err)
	}

	handlers := handlers.NewHandlers(storage_account)

	r := mux.NewRouter()
	r.HandleFunc("/{slug}", handlers.Redirect)

	srv := &http.Server{
		Handler: r,
		Addr:    listenAddr,

		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("About to listen on %s. Go to https://127.0.0.1%s", listenAddr, listenAddr)
	log.Fatal(srv.ListenAndServe())
}