package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gavs.at/shortener/internal/model"
)

func (h *Handlers) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	log.Println("Attempting to redirect to slug:", slug)

	urlBytes, err := h.storage.QueryEntity(r.Context(), "pk001", slug)

	if err != nil {
		log.Printf("Error querying entity: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if len(urlBytes) == 0 {
		http.NotFound(w, r)

		return
	}

	var url model.Redirect
	err = json.Unmarshal(urlBytes, &url)

	if err != nil {
		log.Printf("Error unmarshalling entity: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, url.FullURL, http.StatusFound)
}
