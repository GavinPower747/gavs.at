package handlers

import "net/http"

func (h *Handlers) CreateRedirect(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}