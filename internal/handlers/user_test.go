package handlers

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/middleware"
	"github.com/pvpender/avito-shop/internal/models"
	mock_user "github.com/pvpender/avito-shop/internal/usecase/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUserHandler_Info(t *testing.T) {
	type mockBehaviour func(ctx context.Context, s *mock_user.MockUserUseCase, userId uint32)

	testTable := []struct {
		name          string
		userId        uint32
		mockBehaviour mockBehaviour
		expectedCode  int
		expectedBody  string
	}{
		{
			name:   "valid input",
			userId: 1,
			mockBehaviour: func(ctx context.Context, s *mock_user.MockUserUseCase, userId uint32) {
				s.EXPECT().GetInfo(gomock.Any(), userId).Return(&models.InfoResponse{
					Coins:     1000,
					Inventory: make([]*models.Item, 0),
					CoinHistory: &models.CoinHistory{
						Received: make([]*models.ReceivedCoin, 0),
						Sent:     make([]*models.SendCoinRequest, 0),
					},
				}, nil)
			},
			expectedCode: 200,
			expectedBody: `{"coins":1000,"inventory":[],"coinHistory":{"received":[],"sent":[]}}`,
		},
		{
			name:   "internal server error",
			userId: 1,
			mockBehaviour: func(ctx context.Context, s *mock_user.MockUserUseCase, userId uint32) {
				s.EXPECT().GetInfo(gomock.Any(), userId).Return(nil, &errors.NilPointerError{})
			},
			expectedCode: 500,
			expectedBody: `{"errors":"internal server error"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, &chi.Context{})
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_user.NewMockUserUseCase(c)
			tc.mockBehaviour(ctx, user, tc.userId)
			lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))
			tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
			_, token, uErr := tokenAuth.Encode(map[string]interface{}{"user_id": 1})

			require.NoError(t, uErr)

			handler := NewUserHandler(user, tokenAuth, lgr)

			r := chi.NewRouter()
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(middleware.Authenticator(tokenAuth))
			r.Get("/api/info", handler.Info)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/info", bytes.NewBufferString(""))
			req.Header.Set("Authorization", "Bearer "+token)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
