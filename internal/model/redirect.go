package model

type Redirect struct {
	PartitionKey string `json:"PartitionKey"`
	RowKey       string `json:"RowKey"`
	Slug         string `json:"Slug"`
	FullURL      string `json:"FullURL"`
}

type UpsertLinkRequest struct {
	Slug    string `json:"Slug"`
	FullURL string `json:"FullURL"`
}
