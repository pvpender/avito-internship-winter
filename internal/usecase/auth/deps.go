package auth

import (
	"context"
	"github.com/pvpender/avito-shop/internal/models"
)

//go:generate mockgen --source=deps.go --destination=mocks/mock.go

type AuthUseCase interface {
	Authenticate(ctx context.Context, request *models.AuthRequest) (*models.AuthResponse, error)
	CheckPasswordHash(password, hash string) bool
}
