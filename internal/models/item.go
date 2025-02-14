package models

type Item struct {
	ItemType string `json:"type"`
	Quantity int32  `json:"quantity"`
}
