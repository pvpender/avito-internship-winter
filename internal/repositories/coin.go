package repositories

import (
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/internal/usecase/coin"
	"github.com/pvpender/avito-shop/internal/usecase/user"
)

type PgCoinRepository struct {
	db      *pgxpool.Pool
	getter  *trmpgx.CtxGetter
	builder *squirrel.StatementBuilderType
}

type PgTransmissionRepository struct {
	trManager *manager.Manager
	coin.CoinRepository
	user.UserRepository
}
