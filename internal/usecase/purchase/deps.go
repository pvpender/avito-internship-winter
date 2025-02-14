package purchase

import (
	"context"

	"github.com/pvpender/avito-shop/internal/models"
)

type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, userId uint32, itemId uint32) (int32, error)
	GetUserPurchases(ctx context.Context, userId uint32) ([]*models.Item, error)
}
