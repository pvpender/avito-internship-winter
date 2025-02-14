package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	errInt "github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	jwtAuth *jwtauth.JWTAuth
	user.UserRepository
	logger *slog.Logger
}

func NewAuthUseCase(jwtAuth *jwtauth.JWTAuth, userRepository user.UserRepository, logger *slog.Logger) *AuthUseCase {
	return &AuthUseCase{jwtAuth: jwtAuth, UserRepository: userRepository, logger: logger}
}

func (auc *AuthUseCase) Authenticate(ctx context.Context, request *models.AuthRequest) (*models.AuthResponse, error) {
	authUser, err := auc.UserRepository.GetUserByUsername(ctx, request.Username)
	if errors.Is(err, pgx.ErrNoRows) {
		hashPass, uErr := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if uErr != nil {
			return nil, err
		}

		request.Password = string(hashPass)

		id, uErr := auc.CreateUser(ctx, request)
		if uErr != nil {
			return nil, uErr
		}

		_, tokenAuth, uErr := auc.jwtAuth.Encode(map[string]interface{}{"user_id": id})
		if uErr != nil {
			return nil, uErr
		}

		return &models.AuthResponse{Token: tokenAuth}, nil
	}

	if err != nil {
		return nil, err
	}

	if auc.CheckPasswordHash(request.Password, authUser.Password) {
		_, tokenAuth, uErr := auc.jwtAuth.Encode(map[string]interface{}{"user_id": authUser.UserId})
		if uErr != nil {
			return nil, uErr
		}

		return &models.AuthResponse{Token: tokenAuth}, nil
	}

	return nil, &errInt.InvalidCredentialsError{}
}

func (auc *AuthUseCase) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
