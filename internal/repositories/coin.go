package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase/coin"
)

const CoinTableName = "send_history"

type PgCoinRepository struct {
	db      *pgxpool.Pool
	getter  *trmpgx.CtxGetter
	builder *squirrel.StatementBuilderType
}

func NewPgCoinRepository(
	db *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	builder *squirrel.StatementBuilderType,
	) *PgCoinRepository {
	return &PgCoinRepository{db: db, getter: getter, builder: builder}
}

func (p *PgCoinRepository) CreateTransmission(
	ctx context.Context,
	request *models.CoinOperationWithIds
	) (int32, error) {
	query, args, err := p.builder.Insert(CoinTableName).
		Columns("from_user", "to_user", "amount").
		Values(request.FromUser, request.ToUser, request.Amount).
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

func (p *PgCoinRepository) GetUserTransmissions(
	ctx context.Context,
	userId uint32,
	transmissionType coin.TransmissionType,
	) ([]*models.CoinOperationWithUsernames, error) {
	var expr map[string]interface{}

	switch transmissionType {
	case coin.Sent:
		expr = squirrel.Eq{"from_user": userId}

	case coin.Received:
		expr = squirrel.Eq{"to_user": userId}

	default:
		return nil, &errors.InvalidTransmissionError{}
	}

	query, args, err := p.builder.Select("u.username", "u2.username", "amount").
		From(CoinTableName).
		LeftJoin("users u on u.id = send_history.from_user").
		LeftJoin("users u2 on u2.id = send_history.to_user").
		Where(expr).
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

	coinOperations := make([]*models.CoinOperationWithUsernames, 0)

	for rows.Next() {
		coinOperation := &models.CoinOperationWithUsernames{}

		err = rows.Scan(&coinOperation.FromUser, &coinOperation.ToUser, &coinOperation.Amount)
		if err != nil {
			return nil, err
		}

		coinOperations = append(coinOperations, coinOperation)
	}

	return coinOperations, nil
}
