package model

type Link struct {
	PartitionKey string `json:"PartitionKey"`
	RowKey       string `json:"RowKey"`
	Slug         string `json:"Slug"`
	FullUrl      string `json:"FullUrl"`
}