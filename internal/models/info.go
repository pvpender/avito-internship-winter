package models

import "errors"

type InfoResponse struct {
	Coins       int32        `json:"coins"`
	Inventory   []*Item      `json:"inventory"`
	CoinHistory *CoinHistory `json:"coinHistory"`
}

func CreateInfoResponse(coins int32, inventory []*Item, coinHistory *CoinHistory) (*InfoResponse, error) {
	if coins < 0 {
		return nil, errors.New("invalid coins")
	}

	return &InfoResponse{Coins: coins, Inventory: inventory, CoinHistory: coinHistory}, nil
}
