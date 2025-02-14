package handlers

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/pvpender/avito-shop/internal/handlers"
	"github.com/pvpender/avito-shop/internal/models"
	mock_auth "github.com/pvpender/avito-shop/internal/usecase/auth/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler_CreateUser(t *testing.T) {
	type mockBehaviour func(ctx context.Context, s *mock_auth.MockAuthUseCase, request models.AuthRequest)

	testTable := []struct {
		name          string
		inputBody     string
		mockBehaviour mockBehaviour
		expectedCode  int
		expectedBody  string
	}{
		{
			name:      "valid input",
			inputBody: `{"username:"Nic", "password":"123456"}`,
			mockBehaviour: func(ctx context.Context, s *mock_auth.MockAuthUseCase, request models.AuthRequest) {
				s.EXPECT().Authenticate(ctx, request).Return(&models.AuthResponse{Token: "abc"}, nil)
			},
			expectedCode: 200,
			expectedBody: `{"token":"abc"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_auth.NewMockAuthUseCase(c)
			tc.mockBehaviour(ctx, auth, models.AuthRequest{Username: "Nic", Password: "12345"})
			lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))

			handler := handlers.NewAuthHandler(auth, lgr)

			r := chi.NewRouter()
			r.Post("/api/auth", handler.Auth)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/auth", bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}

}
