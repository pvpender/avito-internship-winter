package usecase

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	mock_coin "github.com/pvpender/avito-shop/internal/usecase/coin/mocks"
	"github.com/pvpender/avito-shop/internal/usecase/common/mocks/mock"
	mock_user "github.com/pvpender/avito-shop/internal/usecase/user/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCoinUseCase_SendCoin(t *testing.T) {
	type mockBehaviour func(
		ctx context.Context,
		u *mock_user.MockUserRepository,
		c *mock_coin.MockCoinRepository,
		userId uint32,
		request *models.SendCoinRequest,

	)

	testTable := []struct {
		name          string
		userId        uint32
		input         *models.SendCoinRequest
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:   "success",
			input:  &models.SendCoinRequest{ToUser: "Jo", Amount: 100},
			userId: 1,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				c *mock_coin.MockCoinRepository,
				userId uint32,
				request *models.SendCoinRequest,
			) {
				u.EXPECT().GetUserById(gomock.Any(), userId).Return(&models.User{UserId: 1, Username: "Nic", Coins: 1000}, nil)
				u.EXPECT().GetUserByUsername(gomock.Any(), request.ToUser).Return(&models.User{UserId: 2, Username: "Jo", Coins: 1000}, nil)
				u.EXPECT().UpdateUserCoins(gomock.Any(), userId, int32(900)).Return(nil)
				u.EXPECT().UpdateUserCoins(gomock.Any(), uint32(2), int32(1100)).Return(nil)
				c.EXPECT().CreateTransmission(gomock.Any(), &models.CoinOperationWithIds{FromUser: 1, ToUser: 2, Amount: 100}).Return(int32(1), nil)
			},
			expectedError: nil,
		},
		{
			name:   "not enough coins",
			input:  &models.SendCoinRequest{ToUser: "Jo", Amount: 1000000},
			userId: 1,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				c *mock_coin.MockCoinRepository,
				userId uint32,
				request *models.SendCoinRequest,
			) {
				u.EXPECT().GetUserById(gomock.Any(), userId).Return(&models.User{UserId: 1, Username: "Nic", Coins: 1000}, nil)
				u.EXPECT().GetUserByUsername(gomock.Any(), request.ToUser).Return(&models.User{UserId: 2, Username: "Jo", Coins: 1000}, nil)
			},
			expectedError: &errors.InvalidAmountError{},
		},
		{
			name:   "failed transaction",
			input:  &models.SendCoinRequest{ToUser: "Jo", Amount: 100},
			userId: 1,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				c *mock_coin.MockCoinRepository,
				userId uint32, request *models.SendCoinRequest,
			) {
				u.EXPECT().GetUserById(gomock.Any(), userId).Return(&models.User{UserId: 1, Username: "Nic", Coins: 1000}, nil)
				u.EXPECT().GetUserByUsername(gomock.Any(), request.ToUser).Return(&models.User{UserId: 2, Username: "Jo", Coins: 1000}, nil)
				u.EXPECT().UpdateUserCoins(gomock.Any(), userId, int32(900)).Return(&errors.NilPointerError{})
			},
			expectedError: &errors.NilPointerError{},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, &chi.Context{})

			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_user.NewMockUserRepository(c)
			coin := mock_coin.NewMockCoinRepository(c)
			tc.mockBehaviour(ctx, user, coin, tc.userId, tc.input)

			trManager := &mock.MockTransactionManager{}

			useCase := NewCoinUseCase(trManager, user, coin)
			err := useCase.SendCoin(ctx, uint32(1), tc.input)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
