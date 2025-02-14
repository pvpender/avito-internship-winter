package models

import "github.com/pvpender/avito-shop/internal/errors"

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

type ReceivedCoin struct {
	FromUser string `json:"fromUser"`
	Amount   int32  `json:"amount"`
}

type CoinOperationWithIds struct {
	FromUser uint32
	ToUser   uint32
	Amount   int32
}

type CoinOperationWithUsernames struct {
	FromUser string
	ToUser   string
	Amount   int32
}

type CoinHistory struct {
	Received []*ReceivedCoin    `json:"received"`
	Sent     []*SendCoinRequest `json:"sent"`
}

func CreateCoinOperationWithIds(fromUser uint32, toUser uint32, amount int32) (*CoinOperationWithIds, error) {
	if (fromUser == toUser) || amount <= 0 {
		return nil, &errors.InvalidCoinOperation{}
	}

	return &CoinOperationWithIds{fromUser, toUser, amount}, nil
}

func (coinRequest *SendCoinRequest) Validate() bool {
	if coinRequest.ToUser == "" || coinRequest.Amount <= 0 {
		return false
	}

	return true
}
