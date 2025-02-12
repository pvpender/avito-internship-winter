package models

import "errors"

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

type ReceivedCoin struct {
	FromUser string `json:"fromUser"`
	Amount   int32  `json:"amount"`
}

type CoinOperation struct {
	FromUser string
	ToUser   string
	Amount   int32
}

type CoinHistory struct {
	Received []*ReceivedCoin    `json:"received"`
	Sent     []*SendCoinRequest `json:"sent"`
}

func CreateSendCoinRequest(toUser string, amount int32) (*SendCoinRequest, error) {
	if toUser == "" || amount <= 0 {
		return nil, errors.New("toUser and amount must be positive")
	}

	return &SendCoinRequest{toUser, amount}, nil
}

func CreateReceivedCoin(fromUser string, amount int32) (*ReceivedCoin, error) {
	if fromUser == "" || amount <= 0 {
		return nil, errors.New("fromUser and amount must be positive")
	}

	return &ReceivedCoin{fromUser, amount}, nil
}

func CreateCoinOperation(fromUser string, toUser string, amount int32) (*CoinOperation, error) {
	if fromUser == "" || toUser == "" || amount <= 0 {
		return nil, errors.New("fromUser and toUser and amount must be positive")
	}

	return &CoinOperation{fromUser, toUser, amount}, nil
}
