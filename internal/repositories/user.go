package repositories

import (
	"context"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/internal/models"
)

type PgUserRepository struct {
	db      *pgxpool.Pool
	getter  *trmpgx.CtxGetter
	builder *squirrel.StatementBuilderType
}

func (p *PgUserRepository) GetUserById(ctx context.Context, userId uint32) (*models.User, error) {
	query, args, err := p.builder.Select("*").From("users").Where(squirrel.Eq{"id": userId}).ToSql()
	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.db)
	row := conn.QueryRow(ctx, query, args...)

	user := &models.User{}
	err = row.Scan(&user.UserId, &user.Username, &user.Password, &user.Coins)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PgUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query, args, err := p.builder.Select("*").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()

	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.db)
	row := conn.QueryRow(ctx, query, args...)

	user := &models.User{}
	err = row.Scan(&user.UserId, &user.Username, &user.Password, &user.Coins)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PgUserRepository) CreateUser(ctx context.Context, request models.AuthRequest) (*models.User, error) {
	query, args, err := p.builder.Insert("user").
		Columns("username", "password", "coins").
		Values(request.Username, request.Password, 1000).
		ToSql()

	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.db)
	user := &models.User{}

	err = conn.QueryRow(ctx, query, args...).Scan(&user.UserId)
	if err != nil {
		return nil, err
	}

	user.Username = request.Username
	user.Password = request.Password
	user.Coins = 1000

	return user, nil
}

func (p *PgUserRepository) UpdateUserCoins(ctx context.Context, username string, coins int32) error {
	query, args, err := p.builder.Update("user").Where(squirrel.Eq{"username": username}).ToSql()
	if err != nil {
		return err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.db)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
