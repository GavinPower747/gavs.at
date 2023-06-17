package handlers

import "gavs.at/shortener/internal/storage"

type Handlers struct {
	storage storage.Account
}

func NewHandlers(storage storage.Account) *Handlers {
	return &Handlers{
		storage: storage,
	}
}