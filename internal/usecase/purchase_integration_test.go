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

type PurchaseIntegrationTestSuit struct {
	suite.Suite
	pool    *pgxpool.Pool
	useCase *PurchaseUseCase
	builder *squirrel.StatementBuilderType
}

func (suite *PurchaseIntegrationTestSuit) SetupSuite() {
	cfg, err := config.LoadConfig(configFile, configType)

	require.NoError(suite.T(), err)

	pgDB, err := database.NewPgPool(cfg)

	require.NoError(suite.T(), err)

	trManager := manager.Must(trmpgx.NewDefaultFactory(pgDB))
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	suite.pool = pgDB

	purchaseRepo := repositories.NewPgPurchaseRepository(pgDB, trmpgx.DefaultCtxGetter, &builder)
	userRepo := repositories.NewPgUserRepository(pgDB, trmpgx.DefaultCtxGetter, &builder)
	itemRepo := repositories.NewPgItemRepository(pgDB, trmpgx.DefaultCtxGetter, &builder)

	suite.useCase = NewPurchaseUseCase(trManager, purchaseRepo, userRepo, itemRepo)
	suite.builder = &builder
}

func (suite *PurchaseIntegrationTestSuit) TearDownSuite() {
	suite.pool.Close()
}

func (suite *PurchaseIntegrationTestSuit) SetupTest() {
	tx, err := suite.pool.Begin(context.Background())
	require.NoError(suite.T(), err)
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		} else {
			_ = tx.Commit(context.Background())
		}
	}()

	_, err = tx.Exec(context.Background(), "DELETE FROM "+repositories.PurchaseTableName+" *")
	require.NoError(suite.T(), err)
	_, err = tx.Exec(context.Background(), "UPDATE "+repositories.UserTableName+" SET coins = 10000 WHERE id = 1")
	require.NoError(suite.T(), err)
}

func (suite *PurchaseIntegrationTestSuit) TestSuccessPurchase() {
	err := suite.useCase.CreatePurchase(context.Background(), uint32(1), "umbrella")
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

	query, args, _ := suite.builder.Select("item_type", "count(item_type) as count").
		From(repositories.PurchaseTableName).
		LeftJoin("users on users.id = user_id").
		LeftJoin("merch on merch.id = item_id").
		Where(squirrel.Eq{"user_id": uint32(1)}).
		Where(squirrel.Eq{"item_type": "umbrella"}).
		GroupBy("item_type").
		ToSql()

	rows, err := tx.Query(context.Background(), query, args...)
	require.NoError(suite.T(), err)

	defer rows.Close()

	items := make([]*models.Item, 0)

	for rows.Next() {
		newItem := &models.Item{}

		err = rows.Scan(&newItem.ItemType, &newItem.Quantity)
		require.NoError(suite.T(), err)

		items = append(items, newItem)
	}

	assert.Equal(suite.T(), 1, len(items))
	assert.Equal(suite.T(), "umbrella", items[0].ItemType)
	assert.Equal(suite.T(), int32(1), items[0].Quantity)

	query, args, _ = suite.builder.Select("coins").
		From(repositories.UserTableName).
		Where(squirrel.Eq{"id": 1}).
		ToSql()
	row := tx.QueryRow(context.Background(), query, args...)
	require.NoError(suite.T(), err)

	var coins int32
	err = row.Scan(&coins)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(10000-200), coins)
}

func (suite *PurchaseIntegrationTestSuit) TestPurchaseNotComplete() {
	err := suite.useCase.CreatePurchase(context.Background(), uint32(1), "zontik")
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

	var coins int32
	err = row.Scan(&coins)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), int32(10000), coins)
}

func TestCoinIntegrationSuite(t *testing.T) {
	suite.Run(t, new(PurchaseIntegrationTestSuit))
}
