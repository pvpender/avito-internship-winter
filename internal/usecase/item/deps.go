package item

import (
	"context"

	"github.com/pvpender/avito-shop/internal/models"
)

type ItemRepository interface {
	GetItemByType(ctx context.Context, itemType string) (*models.Purchase, error)
}
