package handlers

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"

	"gavs.at/shortener/internal/model"
	"gavs.at/shortener/pkg/web"
)

func (h *Handlers) UpsertRedirect(w http.ResponseWriter, r *http.Request) {
	req, vErr := web.ValidateAndParseBody[model.UpsertRedirectRequest](r)

	if vErr != nil {
		web.BadRequest(w, vErr)

		return
	}

	redirect := model.Redirect{
		Slug:    req.Slug,
		FullURL: req.FullURL,
		Entity:  aztables.Entity{PartitionKey: "pk001", RowKey: req.Slug},
	}

	err := h.storage.UpsertEntity(redirect)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
