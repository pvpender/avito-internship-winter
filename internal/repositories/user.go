package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/internal/models"
)

const (
	UserTableName    = "users"
	DefaultUserCoins = 1000
)

type PgUserRepository struct {
	db      *pgxpool.Pool
	getter  *trmpgx.CtxGetter
	builder *squirrel.StatementBuilderType
}

func NewPgUserRepository(
	db *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	builder *squirrel.StatementBuilderType,
) *PgUserRepository {
	return &PgUserRepository{db: db, getter: getter, builder: builder}
}

func (p *PgUserRepository) GetUserById(ctx context.Context, userId uint32) (*models.User, error) {
	query, args, err := p.builder.Select("*").From(UserTableName).Where(squirrel.Eq{"id": userId}).ToSql()
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
		From(UserTableName).
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

func (p *PgUserRepository) CreateUser(ctx context.Context, request *models.AuthRequest) (int32, error) {
	query, args, err := p.builder.Insert(UserTableName).
		Columns("username", "password", "coins").
		Values(request.Username, request.Password, DefaultUserCoins).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return -1, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.db)

	var userId int32

	err = conn.QueryRow(ctx, query, args...).Scan(&userId)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (p *PgUserRepository) UpdateUserCoins(ctx context.Context, userId uint32, coins int32) error {
	query, args, err := p.builder.
		Update(UserTableName).
		Where(squirrel.Eq{"id": userId}).
		Set("coins", coins).
		ToSql()

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
