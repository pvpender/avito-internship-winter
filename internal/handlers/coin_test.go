package handlers

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pvpender/avito-shop/internal/middleware"
	"github.com/pvpender/avito-shop/internal/models"
	mock_coin "github.com/pvpender/avito-shop/internal/usecase/coin/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCoinHandler_SendCoin(t *testing.T) {
	type mockBehaviour func(ctx context.Context, s *mock_coin.MockCoinUseCase, userId uint32, request *models.SendCoinRequest)

	testTable := []struct {
		name              string
		inputBody         string
		expectedParseBody *models.SendCoinRequest
		userId            uint32
		mockBehaviour     mockBehaviour
		expectedCode      int
		expectedBody      string
	}{
		{
			name:              "valid input",
			inputBody:         `{"toUser":"Nic", "amount":1}`,
			expectedParseBody: &models.SendCoinRequest{ToUser: "Nic", Amount: 1},
			userId:            1,
			mockBehaviour: func(ctx context.Context, s *mock_coin.MockCoinUseCase, userId uint32, request *models.SendCoinRequest) {
				s.EXPECT().SendCoin(gomock.Any(), userId, request).Return(nil)
			},
			expectedCode: 200,
			expectedBody: "",
		},
		{
			name:              "invalid amount",
			inputBody:         `{"toUser":"Nic", "amount":-1}`,
			expectedParseBody: &models.SendCoinRequest{ToUser: "Nic", Amount: -1},
			userId:            1,
			mockBehaviour: func(ctx context.Context, s *mock_coin.MockCoinUseCase, userId uint32, request *models.SendCoinRequest) {
			},
			expectedCode: 400,
			expectedBody: `{"errors":"bad request"}`,
		},
		{
			name:              "invalid toUser",
			inputBody:         `{"toUser":"Nic", "amount":1}`,
			expectedParseBody: &models.SendCoinRequest{ToUser: "Nic", Amount: 1},
			userId:            1,
			mockBehaviour: func(ctx context.Context, s *mock_coin.MockCoinUseCase, userId uint32, request *models.SendCoinRequest) {
				s.EXPECT().SendCoin(gomock.Any(), userId, request).Return(pgx.ErrNoRows)
			},
			expectedCode: 400,
			expectedBody: `{"errors":"bad request"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, &chi.Context{})

			c := gomock.NewController(t)
			defer c.Finish()

			coin := mock_coin.NewMockCoinUseCase(c)
			tc.mockBehaviour(ctx, coin, tc.userId, tc.expectedParseBody)

			lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))
			tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
			_, token, uErr := tokenAuth.Encode(map[string]interface{}{"user_id": 1})

			require.NoError(t, uErr)

			handler := NewCoinHandler(coin, tokenAuth, lgr)

			r := chi.NewRouter()
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(middleware.Authenticator(tokenAuth))
			r.Post("/api/sendCoin", handler.SendCoin)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBufferString(tc.inputBody))
			req.Header.Set("Authorization", "Bearer "+token)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
