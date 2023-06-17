package handlers

import "gavs.at/shortener/internal/storage"

type Handlers struct {
	storage storage.Account
}

func NewHandlers(storageAccount storage.Account) *Handlers {
	return &Handlers{
		storage: storageAccount,
	}
}
