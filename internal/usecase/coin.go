package usecase

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase/coin"
	"github.com/pvpender/avito-shop/internal/usecase/user"
)

type CoinUseCase struct {
	trManager *manager.Manager
	coin.CoinRepository
	user.UserRepository
}

func NewCoinUseCase(trManager *manager.Manager, userRepository user.UserRepository, coinRepository coin.CoinRepository) *CoinUseCase {
	return &CoinUseCase{
		trManager:      trManager,
		UserRepository: userRepository,
		CoinRepository: coinRepository,
	}
}

func (c *CoinUseCase) SendCoin(ctx context.Context, userId uint32, request *models.SendCoinRequest) error {
	currentUser, err := c.UserRepository.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	receiver, err := c.UserRepository.GetUserByUsername(ctx, request.ToUser)
	if err != nil {
		return err
	}

	var newAmount int32 = currentUser.Coins - request.Amount
	if newAmount < 0 {
		return &errors.InvalidAmount{}
	}

	transmission, err := models.CreateCoinOperationWithIds(currentUser.UserId, receiver.UserId, request.Amount)
	if err != nil {
		return err
	}

	err = c.trManager.Do(ctx, func(ctx context.Context) error {
		trErr := c.UserRepository.UpdateUserCoins(ctx, currentUser.UserId, newAmount)
		if trErr != nil {
			return trErr
		}

		trErr = c.UserRepository.UpdateUserCoins(ctx, receiver.UserId, request.Amount+request.Amount)
		if trErr != nil {
			return trErr
		}

		_, trErr = c.CoinRepository.CreateTransmission(ctx, transmission)
		if trErr != nil {
			return trErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
