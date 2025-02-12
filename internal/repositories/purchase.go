package repositories

import (
	"context"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase/item"
	"github.com/pvpender/avito-shop/internal/usecase/purchase"
	"github.com/pvpender/avito-shop/internal/usecase/user"
)

type PgPurchaseRepository struct {
	db      *pgxpool.Pool
	getter  *trmpgx.CtxGetter
	builder *squirrel.StatementBuilderType
}

func (p *PgPurchaseRepository) CreatePurchase(ctx context.Context, username string, itemType string) (int32, error) {
	query, args, err := p.builder.Insert("purchase_history").
		Columns("username", "item_type").
		Values(username, itemType).
		ToSql()

	if err != nil {
		return 0, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.db)

	var id int32
	err = conn.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *PgPurchaseRepository) GetUserPurchases(ctx context.Context, username string) ([]*models.Item, error) {
	query, args, err := p.builder.Select("item_type", "count(item_type) as count)").
		From("purchase_history").
		Where(squirrel.Eq{"username": username}).
		GroupBy("item_id").
		ToSql()

	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.db)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]*models.Item, 0)
	for rows.Next() {
		newItem := &models.Item{}
		err = rows.Scan(&newItem.ItemType, &newItem.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, newItem)
	}

	return items, nil
}

type PgPurchaseTransactionRepository struct {
	trManager *manager.Manager
	purchase.PurchaseRepository
	user.UserRepository
	item.ItemRepository
}

func (p *PgPurchaseTransactionRepository) CreatePurchaseTransmission(ctx context.Context, username string, itemType string) error {
	purchasedItem, err := p.ItemRepository.GetItemByType(ctx, itemType)
	if purchasedItem == nil {
		return err
	}

	updatableUser, err := p.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	err = p.trManager.Do(ctx, func(ctx context.Context) error {
		if _, errTr := p.PurchaseRepository.CreatePurchase(ctx, username, itemType); err != nil {
			return errTr
		}

		if errTr := p.UserRepository.UpdateUserCoins(ctx, username, updatableUser.Coins-purchasedItem.Price); err != nil {
			return errTr
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
