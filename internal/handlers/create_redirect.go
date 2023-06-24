package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gavs.at/shortener/internal/model"
	"gavs.at/shortener/pkg/web"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func (h *Handlers) UpsertRedirect(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		web.BadRequest(w, map[string]string{"body": "Could not read request body"})

		return
	}

	var req model.UpsertRedirectRequest

	err = json.Unmarshal(b, &req)

	if err != nil {
		web.BadRequest(w, map[string]string{"body": "Could not unmarshal request body"})

		return
	}

	isValid, vErr := req.Validate()

	if !isValid {
		web.BadRequest(w, vErr.Errors)

		return
	}

	redirect := model.Redirect{
		Slug:    req.Slug,
		FullURL: req.FullURL,
		Entity:  aztables.Entity{PartitionKey: "pk001", RowKey: req.Slug},
	}

	err = h.storage.UpsertEntity(redirect)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}
