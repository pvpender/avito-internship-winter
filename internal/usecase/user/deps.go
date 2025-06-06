package user

import (
	"context"

	"github.com/pvpender/avito-shop/internal/models"
)

//go:generate mockgen --source=deps.go --destination=mocks/mock.go

type UserRepository interface {
	GetUserById(ctx context.Context, userId uint32) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, request *models.AuthRequest) (int32, error)
	UpdateUserCoins(ctx context.Context, userId uint32, coins int32) error
}

type UserUseCase interface {
	CreateUser(ctx context.Context, user *models.AuthRequest) error
	GetInfo(ctx context.Context, userId uint32) (*models.InfoResponse, error)
}
