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
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/middleware"
	mock_purchase "github.com/pvpender/avito-shop/internal/usecase/purchase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPurchaseHandler_CreatePurchase(t *testing.T) {
	type mockBehaviour func(ctx context.Context, s *mock_purchase.MockPurchaseUseCase, userId uint32, itemType string)

	testTable := []struct {
		name          string
		userId        uint32
		itemType      string
		mockBehaviour mockBehaviour
		expectedCode  int
		expectedBody  string
	}{
		{
			name:     "valid input",
			userId:   1,
			itemType: "umbrella",
			mockBehaviour: func(ctx context.Context, s *mock_purchase.MockPurchaseUseCase, userId uint32, itemType string) {
				s.EXPECT().CreatePurchase(gomock.Any(), userId, itemType).Return(nil)
			},
			expectedCode: 200,
			expectedBody: "",
		},
		{
			name:     "invalid input",
			userId:   1,
			itemType: "zontik",
			mockBehaviour: func(ctx context.Context, s *mock_purchase.MockPurchaseUseCase, userId uint32, itemType string) {
				s.EXPECT().CreatePurchase(gomock.Any(), userId, itemType).Return(pgx.ErrNoRows)
			},
			expectedCode: 400,
			expectedBody: `{"errors":"bad request"}`,
		},
		{
			name:     "incorrect coins amount",
			userId:   1,
			itemType: "umbrella",
			mockBehaviour: func(ctx context.Context, s *mock_purchase.MockPurchaseUseCase, userId uint32, itemType string) {
				s.EXPECT().CreatePurchase(gomock.Any(), userId, itemType).Return(&errors.PurchaseError{})
			},
			expectedCode: 400,
			expectedBody: `{"errors":"bad request"}`,
		},
		{
			name:     "internal server error",
			userId:   1,
			itemType: "umbrella",
			mockBehaviour: func(ctx context.Context, s *mock_purchase.MockPurchaseUseCase, userId uint32, itemType string) {
				s.EXPECT().CreatePurchase(gomock.Any(), userId, itemType).Return(&errors.NilPointerError{})
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

			purchase := mock_purchase.NewMockPurchaseUseCase(c)
			tc.mockBehaviour(ctx, purchase, tc.userId, tc.itemType)

			lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))
			tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
			_, token, uErr := tokenAuth.Encode(map[string]interface{}{"user_id": 1})

			require.NoError(t, uErr)

			handler := NewPurchaseHandler(purchase, tokenAuth, lgr)

			r := chi.NewRouter()
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(middleware.Authenticator(tokenAuth))
			r.Get("/api/buy/{item}", handler.Purchase)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/buy/"+tc.itemType, bytes.NewBufferString(""))
			req.Header.Set("Authorization", "Bearer "+token)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
