package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	mock_user "github.com/pvpender/avito-shop/internal/usecase/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUseCase_Authenticate(t *testing.T) {
	type mockBehaviour func(ctx context.Context, s *mock_user.MockUserRepository, request *models.AuthRequest)

	testTable := []struct {
		name          string
		input         *models.AuthRequest
		mockBehaviour mockBehaviour
		jwtInfo       map[string]interface{}
		expectedError error
	}{
		{
			name:  "success",
			input: &models.AuthRequest{Username: "test", Password: "test"},
			mockBehaviour: func(ctx context.Context, s *mock_user.MockUserRepository, request *models.AuthRequest) {
				s.EXPECT().GetUserByUsername(gomock.Any(), request.Username).Return(nil, pgx.ErrNoRows)
				s.EXPECT().CreateUser(gomock.Any(), request).Return(int32(1), nil)
			},
			jwtInfo:       map[string]interface{}{"user_id": 1},
			expectedError: nil,
		},
		{
			name:  "invalid credentials",
			input: &models.AuthRequest{Username: "test", Password: "test"},
			mockBehaviour: func(ctx context.Context, s *mock_user.MockUserRepository, request *models.AuthRequest) {
				s.EXPECT().GetUserByUsername(gomock.Any(), request.Username).Return(&models.User{Username: "test", Password: "adfsfadfgs"}, nil)
			},
			expectedError: &errors.InvalidCredentialsError{},
		},
		{
			name:  "success login",
			input: &models.AuthRequest{Username: "test", Password: "1"},
			mockBehaviour: func(ctx context.Context, s *mock_user.MockUserRepository, request *models.AuthRequest) {
				hashPass, _ := bcrypt.GenerateFromPassword([]byte("1"), bcrypt.DefaultCost)
				s.EXPECT().GetUserByUsername(gomock.Any(), request.Username).Return(&models.User{
					UserId:   1,
					Username: "test",
					Password: string(hashPass),
				}, nil)
			},
			jwtInfo:       map[string]interface{}{"user_id": 1},
			expectedError: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, &chi.Context{})

			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_user.NewMockUserRepository(c)
			tc.mockBehaviour(ctx, user, tc.input)

			lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))
			tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

			useCase := NewAuthUseCase(tokenAuth, user, lgr)
			response, err := useCase.Authenticate(ctx, tc.input)

			assert.Equal(t, tc.expectedError, err)

			if tc.expectedError == nil {
				_, token, errT := tokenAuth.Encode(tc.jwtInfo)
				require.NoError(t, errT)
				assert.Equal(t, token, response.Token)
			}
		})
	}
}
