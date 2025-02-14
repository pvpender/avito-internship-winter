package item

import (
	"context"

	"github.com/pvpender/avito-shop/internal/models"
)

//go:generate mockgen --source=deps.go --destination=mocks/mock.go

type ItemRepository interface {
	GetItemByType(ctx context.Context, itemType string) (*models.Purchase, error)
}
