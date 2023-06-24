package model

import (
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
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

func (r *UpsertRedirectRequest) Validate() (bool, ValidationError) {
	errs := &ValidationError{Errors: make(map[string]string)}

	if r.Slug == "" {
		errs.AddError("Slug", "Slug is required")
	}

	if r.FullURL == "" {
		errs.AddError("FullURL", "FullURL is required")
	}

	if _, err := url.Parse(r.FullURL); err != nil {
		errs.AddError("FullURL", "URL provided is not valid")
	}

	return !errs.HasErrors(), *errs
}
