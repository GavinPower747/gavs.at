package model

import (
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"

	"gavs.at/shortener/pkg/web"
)

type Redirect struct {
	aztables.Entity
	Slug    string `json:"Slug"`
	FullURL string `json:"FullURL"`
}

type UpsertRedirectRequest struct {
	Slug    string `json:"Slug"`
	FullURL string `json:"FullURL"`
}

func (r UpsertRedirectRequest) Validate() (bool, web.ValidationError) {
	errs := &web.ValidationError{Errors: make(map[string]string)}

	if r.Slug == "" {
		errs.AddError("Slug", "Slug is required")
	}

	if r.FullURL == "" {
		errs.AddError("FullURL", "FullURL is required")
	}

	if _, err := url.ParseRequestURI(r.FullURL); err != nil {
		errs.AddError("FullURL", "URL provided is not valid")
	}

	return !errs.HasErrors(), *errs
}
