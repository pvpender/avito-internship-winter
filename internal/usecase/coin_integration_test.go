//go:build integration

package usecase

import (
	"context"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/config"
	"github.com/pvpender/avito-shop/database"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	configFile = "../../config_test.yaml"
	configType = "yaml"
)

type CoinIntegrationTestSuit struct {
	suite.Suite
	pool    *pgxpool.Pool
	useCase *CoinUseCase
	builder *squirrel.StatementBuilderType
}

func (suite *CoinIntegrationTestSuit) SetupSuite() {
	cfg, err := config.LoadConfig(configFile, configType)

	require.NoError(suite.T(), err)

	pgDB, err := database.NewPgPool(cfg)

	require.NoError(suite.T(), err)

	trManager := manager.Must(trmpgx.NewDefaultFactory(pgDB))
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	suite.pool = pgDB

	userRepo := repositories.NewPgUserRepository(pgDB, trmpgx.DefaultCtxGetter, &builder)
	coinRepo := repositories.NewPgCoinRepository(pgDB, trmpgx.DefaultCtxGetter, &builder)

	suite.useCase = NewCoinUseCase(trManager, userRepo, coinRepo)
	suite.builder = &builder
}

func (suite *CoinIntegrationTestSuit) TearDownSuite() {
	suite.pool.Close()
}

func (suite *CoinIntegrationTestSuit) SetupTest() {
	tx, err := suite.pool.Begin(context.Background())
	require.NoError(suite.T(), err)
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		} else {
			_ = tx.Commit(context.Background())
		}
	}()

	_, err = tx.Exec(context.Background(), "DELETE FROM "+repositories.CoinTableName+" *")
	require.NoError(suite.T(), err)
	_, err = tx.Exec(context.Background(), "UPDATE "+repositories.UserTableName+" SET coins = 10000 WHERE id = 1")
	require.NoError(suite.T(), err)
	_, err = tx.Exec(context.Background(), "UPDATE "+repositories.UserTableName+" SET coins = 10000 WHERE id = 2")
	require.NoError(suite.T(), err)
}

func (suite *CoinIntegrationTestSuit) TestSuccessCoinSend() {
	err := suite.useCase.SendCoin(context.Background(), uint32(1), &models.SendCoinRequest{ToUser: "Jo", Amount: 100})
	require.NoError(suite.T(), err)

	tx, err := suite.pool.Begin(context.Background())
	require.NoError(suite.T(), err)

	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		} else {
			_ = tx.Commit(context.Background())
		}
	}()

	query, args, _ := suite.builder.Select("coins").
		From(repositories.UserTableName).
		Where(squirrel.Eq{"id": 1}).
		ToSql()
	row := tx.QueryRow(context.Background(), query, args...)
	require.NoError(suite.T(), err)

	var coinsSender, coinsReceiver int32
	err = row.Scan(&coinsSender)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), int32(10000-100), coinsSender)

	query, args, _ = suite.builder.Select("coins").
		From(repositories.UserTableName).
		Where(squirrel.Eq{"id": 2}).
		ToSql()

	row = tx.QueryRow(context.Background(), query, args...)
	require.NoError(suite.T(), err)

	err = row.Scan(&coinsReceiver)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), int32(10000+100), coinsReceiver)
}

func (suite *CoinIntegrationTestSuit) TestSendNotComplete() {
	err := suite.useCase.SendCoin(context.Background(), uint32(1), &models.SendCoinRequest{ToUser: "DONTEXISTS", Amount: 100})
	assert.Equal(suite.T(), pgx.ErrNoRows, err)

	tx, err := suite.pool.Begin(context.Background())
	require.NoError(suite.T(), err)

	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		} else {
			_ = tx.Commit(context.Background())
		}
	}()

	query, args, _ := suite.builder.Select("coins").
		From(repositories.UserTableName).
		Where(squirrel.Eq{"id": 1}).
		ToSql()
	row := tx.QueryRow(context.Background(), query, args...)
	require.NoError(suite.T(), err)

	var coinsSender int32
	err = row.Scan(&coinsSender)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), int32(10000), coinsSender)
}

func TestPurchaseIntegrationSuite(t *testing.T) {
	suite.Run(t, new(CoinIntegrationTestSuit))
}
