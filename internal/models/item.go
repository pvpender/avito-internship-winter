package models

import "errors"

type Item struct {
	ItemType string `json:"type"`
	Quantity int32  `json:"quantity"`
}

func CreateItem(itemType string, quantity int32) (*Item, error) {
	if itemType == "" || quantity <= 0 {
		return nil, errors.New("itemType or quantity must be positive")
	}

	return &Item{ItemType: itemType, Quantity: quantity}, nil
}
