package usecase

import (
	"context"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/usecase/common"
	"github.com/pvpender/avito-shop/internal/usecase/item"
	"github.com/pvpender/avito-shop/internal/usecase/purchase"
	"github.com/pvpender/avito-shop/internal/usecase/user"
)

type PurchaseUseCase struct {
	trManager common.TransactionManager
	purchase.PurchaseRepository
	user.UserRepository
	item.ItemRepository
}

func NewPurchaseUseCase(
	trManager common.TransactionManager,
	purchaseRepository purchase.PurchaseRepository,
	userRepository user.UserRepository,
	itemRepository item.ItemRepository,
) *PurchaseUseCase {
	return &PurchaseUseCase{
		trManager:          trManager,
		PurchaseRepository: purchaseRepository,
		UserRepository:     userRepository,
		ItemRepository:     itemRepository,
	}
}

func (p PurchaseUseCase) CreatePurchase(ctx context.Context, userId uint32, itemType string) error {
	purchasedItem, err := p.ItemRepository.GetItemByType(ctx, itemType)
	if purchasedItem == nil {
		return err
	}

	updatableUser, err := p.UserRepository.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	newAmount := updatableUser.Coins - purchasedItem.Price
	if newAmount < 0 {
		return &errors.PurchaseError{}
	}

	err = p.trManager.Do(ctx, func(ctx context.Context) error {
		if _, errTr := p.PurchaseRepository.CreatePurchase(ctx, userId, purchasedItem.ItemId); errTr != nil {
			return errTr
		}

		if errTr := p.UserRepository.UpdateUserCoins(ctx, userId, newAmount); errTr != nil {
			return errTr
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
