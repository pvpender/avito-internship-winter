package usecase

import (
	"context"

	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase/coin"
	"github.com/pvpender/avito-shop/internal/usecase/purchase"
	"github.com/pvpender/avito-shop/internal/usecase/user"
)

type UserUseCase struct {
	user.UserRepository
	purchase.PurchaseRepository
	coin.CoinRepository
}

func NewUserUseCase(
	userRepository user.UserRepository,
	purchaseRepository purchase.PurchaseRepository,
	coinRepository coin.CoinRepository,
) *UserUseCase {
	return &UserUseCase{
		UserRepository:     userRepository,
		PurchaseRepository: purchaseRepository,
		CoinRepository:     coinRepository,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, user *models.AuthRequest) error {
	_, err := u.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) GetInfo(ctx context.Context, userId uint32) (*models.InfoResponse, error) {
	currentUser, err := u.UserRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	items, err := u.PurchaseRepository.GetUserPurchases(ctx, userId)
	if err != nil {
		return nil, err
	}

	coinOperations, err := u.CoinRepository.GetUserTransmissions(ctx, userId, coin.Received)
	if err != nil {
		return nil, err
	}

	received := make([]*models.ReceivedCoin, len(coinOperations))
	for i, coinOperation := range coinOperations {
		received[i] = &models.ReceivedCoin{FromUser: coinOperation.FromUser, Amount: coinOperation.Amount}
	}

	coinOperations, err = u.CoinRepository.GetUserTransmissions(ctx, userId, coin.Sent)
	if err != nil {
		return nil, err
	}

	sent := make([]*models.SendCoinRequest, len(coinOperations))
	for i, coinOperation := range coinOperations {
		sent[i] = &models.SendCoinRequest{ToUser: coinOperation.ToUser, Amount: coinOperation.Amount}
	}

	history := &models.CoinHistory{Received: received, Sent: sent}

	return &models.InfoResponse{Coins: currentUser.Coins, Inventory: items, CoinHistory: history}, nil
}
