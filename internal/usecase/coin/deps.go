package coin

import (
	"context"

	"github.com/pvpender/avito-shop/internal/models"
)

type TransmissionType string

const (
	Received TransmissionType = "received"
	Sent     TransmissionType = "sent"
)

//go:generate mockgen --source=deps.go --destination=mocks/mock.go

type CoinRepository interface {
	CreateTransmission(ctx context.Context, request *models.CoinOperationWithIds) (int32, error)
	GetUserTransmissions(
		ctx context.Context,
		userId uint32,
		transmissionType TransmissionType,
	) ([]*models.CoinOperationWithUsernames, error)
}

type CoinUseCase interface {
	SendCoin(ctx context.Context, userId uint32, request *models.SendCoinRequest) error
}
