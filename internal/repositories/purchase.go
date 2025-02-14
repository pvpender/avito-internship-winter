package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/internal/models"
)

const PurchaseTableName = "purchases_history"

type PgPurchaseRepository struct {
	db      *pgxpool.Pool
	getter  *trmpgx.CtxGetter
	builder *squirrel.StatementBuilderType
}

func NewPgPurchaseRepository(db *pgxpool.Pool, getter *trmpgx.CtxGetter, builder *squirrel.StatementBuilderType) *PgPurchaseRepository {
	return &PgPurchaseRepository{db: db, getter: getter, builder: builder}
}

func (p *PgPurchaseRepository) CreatePurchase(ctx context.Context, userId uint32, itemId uint32) (int32, error) {
	query, args, err := p.builder.Insert(PurchaseTableName).
		Columns("user_id", "item_id").
		Values(userId, itemId).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return -1, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.db)

	var id int32

	err = conn.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (p *PgPurchaseRepository) GetUserPurchases(ctx context.Context, userId uint32) ([]*models.Item, error) {
	query, args, err := p.builder.Select("item_type", "count(item_type) as count").
		From(PurchaseTableName).
		LeftJoin("users on users.id = user_id").
		LeftJoin("merch on merch.id = item_id").
		Where(squirrel.Eq{"user_id": userId}).
		GroupBy("item_type").
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
