package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/internal/models"
)

const ItemTableName = "merch"

type PgItemRepository struct {
	db      *pgxpool.Pool
	getter  *trmpgx.CtxGetter
	builder *squirrel.StatementBuilderType
}

func NewPgItemRepository(
	db *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	builder *squirrel.StatementBuilderType,
) *PgItemRepository {
	return &PgItemRepository{db: db, getter: getter, builder: builder}
}

func (p *PgItemRepository) GetItemByType(ctx context.Context, itemType string) (*models.Purchase, error) {
	query, args, err := p.builder.Select("*").
		From(ItemTableName).
		Where(squirrel.Eq{"item_type": itemType}).
		ToSql()

	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.db)

	var item models.Purchase
	err = conn.QueryRow(ctx, query, args...).Scan(&item.ItemId, &item.ItemType, &item.Price)

	if err != nil {
		return nil, err
	}

	return &item, nil
}
