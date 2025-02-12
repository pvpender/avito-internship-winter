package coin

import (
	"context"
	"github.com/pvpender/avito-shop/internal/models"
)

type CoinRepository interface {
	CreateTransmission(ctx context.Context, request models.SendCoinRequest) (*models.CoinOperation, error)
	GetUserSendTransmission(ctx context.Context, username string) ([]*models.CoinOperation, error)
	GetUserReceiveTransmission(ctx context.Context, username string) ([]*models.CoinOperation, error)
}

type TransmissionRepository interface {
	SendCoin(ctx context.Context, username string, request models.SendCoinRequest) (*models.CoinOperation, error)
}
