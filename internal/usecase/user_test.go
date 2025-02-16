package usecase

import (
	"context"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase/coin"
	mock_coin "github.com/pvpender/avito-shop/internal/usecase/coin/mocks"
	mock_purchase "github.com/pvpender/avito-shop/internal/usecase/purchase/mocks"
	mock_user "github.com/pvpender/avito-shop/internal/usecase/user/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	type mockBehaviour func(
		ctx context.Context,
		u *mock_user.MockUserRepository,
		request *models.AuthRequest,
	)

	testTable := []struct {
		name          string
		request       *models.AuthRequest
		behaviour     mockBehaviour
		expectedError error
	}{
		{
			name:    "success",
			request: &models.AuthRequest{Username: "test", Password: "test"},
			behaviour: func(ctx context.Context, u *mock_user.MockUserRepository, request *models.AuthRequest) {
				u.EXPECT().CreateUser(gomock.Any(), request).Return(int32(1), nil)
			},
			expectedError: nil,
		},
		{
			name:          "nil pointer",
			request:       nil,
			behaviour:     func(ctx context.Context, u *mock_user.MockUserRepository, request *models.AuthRequest) {},
			expectedError: &errors.NilPointerError{},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, &chi.Context{})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user := mock_user.NewMockUserRepository(ctrl)
			tc.behaviour(ctx, user, tc.request)

			useCase := NewUserUseCase(user, nil, nil)
			err := useCase.CreateUser(ctx, tc.request)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUserUseCase_GetInfo(t *testing.T) {
	type mockBehaviour func(
		ctx context.Context,
		u *mock_user.MockUserRepository,
		c *mock_coin.MockCoinRepository,
		p *mock_purchase.MockPurchaseRepository,
		userId uint32,
	)

	testTable := []struct {
		name          string
		userId        uint32
		mockBehaviour mockBehaviour
		response      *models.InfoResponse
		expectedError error
	}{
		{
			name:   "success",
			userId: 1,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				c *mock_coin.MockCoinRepository,
				p *mock_purchase.MockPurchaseRepository,
				userId uint32,
			) {
				u.EXPECT().GetUserById(gomock.Any(), userId).Return(&models.User{Username: "Nic", Coins: 1000}, nil)
				purList := make([]*models.Item, 0)
				p.EXPECT().GetUserPurchases(gomock.Any(), userId).Return(purList, nil)

				var rTr []*models.CoinOperationWithUsernames
				sTr := []*models.CoinOperationWithUsernames{{FromUser: "Nic", ToUser: "Jo", Amount: 100}}
				c.EXPECT().GetUserTransmissions(gomock.Any(), userId, coin.Received).Return(rTr, nil)
				c.EXPECT().GetUserTransmissions(gomock.Any(), userId, coin.Sent).Return(sTr, nil)
			},
			response: &models.InfoResponse{
				Coins:     1000,
				Inventory: make([]*models.Item, 0),
				CoinHistory: &models.CoinHistory{
					Received: make([]*models.ReceivedCoin, 0),
					Sent:     []*models.SendCoinRequest{{ToUser: "Jo", Amount: 100}},
				},
			},
		},
		{
			name:   "invalid user id",
			userId: 0,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				c *mock_coin.MockCoinRepository,
				p *mock_purchase.MockPurchaseRepository,
				userId uint32,
			) {
				u.EXPECT().GetUserById(gomock.Any(), userId).Return(nil, pgx.ErrNoRows)
			},
			response:      nil,
			expectedError: pgx.ErrNoRows,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, &chi.Context{})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user := mock_user.NewMockUserRepository(ctrl)
			purchase := mock_purchase.NewMockPurchaseRepository(ctrl)
			coinRepo := mock_coin.NewMockCoinRepository(ctrl)

			useCase := NewUserUseCase(user, purchase, coinRepo)
			tc.mockBehaviour(ctx, user, coinRepo, purchase, tc.userId)

			info, err := useCase.GetInfo(ctx, tc.userId)

			assert.Equal(t, tc.response, info)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
