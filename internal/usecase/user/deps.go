package user

import (
	"context"
	"github.com/pvpender/avito-shop/internal/models"
)

type UserRepository interface {
	GetUserById(ctx context.Context, userId uint32) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, request models.AuthRequest) (*models.User, error)
	UpdateUserCoins(ctx context.Context, username string, coins int32) error
}
