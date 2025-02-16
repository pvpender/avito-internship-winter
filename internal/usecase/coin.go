package usecase

import (
	"context"

	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase/coin"
	"github.com/pvpender/avito-shop/internal/usecase/common"
	"github.com/pvpender/avito-shop/internal/usecase/user"
)

type CoinUseCase struct {
	trManager common.TransactionManager
	coin.CoinRepository
	user.UserRepository
}

func NewCoinUseCase(
	trManager common.TransactionManager,
	userRepository user.UserRepository,
	coinRepository coin.CoinRepository,
) *CoinUseCase {
	return &CoinUseCase{
		trManager:      trManager,
		UserRepository: userRepository,
		CoinRepository: coinRepository,
	}
}

func (c *CoinUseCase) SendCoin(ctx context.Context, userId uint32, request *models.SendCoinRequest) error {
	if request == nil {
		return &errors.NilPointerError{}
	}

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
		return &errors.InvalidAmountError{}
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

		trErr = c.UserRepository.UpdateUserCoins(ctx, receiver.UserId, receiver.Coins+request.Amount)
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
