package usecase

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase/common/mocks/mock"
	mock_item "github.com/pvpender/avito-shop/internal/usecase/item/mocks"
	mock_purchase "github.com/pvpender/avito-shop/internal/usecase/purchase/mocks"
	mock_user "github.com/pvpender/avito-shop/internal/usecase/user/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestPurchaseUseCase_CreatePurchase(t *testing.T) {
	type mockBehaviour func(
		ctx context.Context,
		u *mock_user.MockUserRepository,
		p *mock_purchase.MockPurchaseRepository,
		i *mock_item.MockItemRepository,
		userId uint32,
		itemType string,
	)

	testTable := []struct {
		name          string
		userId        uint32
		itemType      string
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:     "success",
			itemType: "cup",
			userId:   1,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				p *mock_purchase.MockPurchaseRepository,
				i *mock_item.MockItemRepository,
				userId uint32,
				itemType string,
			) {
				i.EXPECT().GetItemByType(gomock.Any(), itemType).Return(&models.Purchase{ItemId: 1, ItemType: itemType, Price: 10}, nil)
				u.EXPECT().GetUserById(gomock.Any(), userId).Return(&models.User{UserId: 1, Username: "Nic", Coins: 1000}, nil)
				p.EXPECT().CreatePurchase(gomock.Any(), userId, uint32(1)).Return(int32(1), nil)
				u.EXPECT().UpdateUserCoins(gomock.Any(), userId, int32(990)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "invalid itemType",
			itemType: "crujka",
			userId:   1,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				p *mock_purchase.MockPurchaseRepository,
				i *mock_item.MockItemRepository,
				userId uint32,
				itemType string,
			) {
				i.EXPECT().GetItemByType(gomock.Any(), itemType).Return(nil, pgx.ErrNoRows)
			},
			expectedError: pgx.ErrNoRows,
		},
		{
			name:     "invalid userId",
			itemType: "cup",
			userId:   0,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				p *mock_purchase.MockPurchaseRepository,
				i *mock_item.MockItemRepository,
				userId uint32,
				itemType string,
			) {
				i.EXPECT().GetItemByType(gomock.Any(), itemType).Return(&models.Purchase{ItemId: 1, ItemType: itemType, Price: 10}, nil)
				u.EXPECT().GetUserById(gomock.Any(), userId).Return(nil, pgx.ErrNoRows)
			},
			expectedError: pgx.ErrNoRows,
		},
		{
			name:     "not enough coins",
			itemType: "cup",
			userId:   1,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				p *mock_purchase.MockPurchaseRepository,
				i *mock_item.MockItemRepository,
				userId uint32,
				itemType string,
			) {
				i.EXPECT().GetItemByType(gomock.Any(), itemType).Return(&models.Purchase{ItemId: 1, ItemType: itemType, Price: 10}, nil)
				u.EXPECT().GetUserById(gomock.Any(), userId).Return(&models.User{UserId: 1, Username: "Nic", Coins: 5}, nil)
			},
			expectedError: &errors.PurchaseError{},
		},
		{
			name:     "transaction error",
			itemType: "cup",
			userId:   1,
			mockBehaviour: func(
				ctx context.Context,
				u *mock_user.MockUserRepository,
				p *mock_purchase.MockPurchaseRepository,
				i *mock_item.MockItemRepository,
				userId uint32,
				itemType string,
			) {
				i.EXPECT().GetItemByType(gomock.Any(), itemType).Return(&models.Purchase{ItemId: 1, ItemType: itemType, Price: 10}, nil)
				u.EXPECT().GetUserById(gomock.Any(), userId).Return(&models.User{UserId: 1, Username: "Nic", Coins: 1000}, nil)
				p.EXPECT().CreatePurchase(gomock.Any(), userId, uint32(1)).Return(int32(1), nil)
				u.EXPECT().UpdateUserCoins(gomock.Any(), userId, int32(990)).Return(&errors.NilPointerError{})
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
			purchase := mock_purchase.NewMockPurchaseRepository(c)
			item := mock_item.NewMockItemRepository(c)
			tc.mockBehaviour(ctx, user, purchase, item, tc.userId, tc.itemType)

			trManager := &mock.MockTransactionManager{}

			useCase := NewPurchaseUseCase(trManager, purchase, user, item)

			err := useCase.CreatePurchase(ctx, tc.userId, tc.itemType)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
