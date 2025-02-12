package purchase

import (
	"context"
	"github.com/pvpender/avito-shop/internal/models"
)

type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, username string, itemType string) (int32, error)
	GetUserPurchases(ctx context.Context, username string) ([]*models.Item, error)
}

type PurchaseTransmissionRepository interface {
	CreatePurchaseTransmission(ctx context.Context, username string, itemType string) error
}
