package handlers

import "gavs.at/shortener/internal/storage"

type Handlers struct {
	storage storage.Storage
}

func NewHandlers(storage storage.Storage) *Handlers {
	return &Handlers{
		storage: storage,
	}
}