package handlers

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	mock_auth "github.com/pvpender/avito-shop/internal/usecase/auth/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthHandler_Auth(t *testing.T) {
	type mockBehaviour func(ctx context.Context, s *mock_auth.MockAuthUseCase, request *models.AuthRequest)

	testTable := []struct {
		name              string
		inputBody         string
		expectedParseBody *models.AuthRequest
		mockBehaviour     mockBehaviour
		expectedCode      int
		expectedBody      string
	}{
		{
			name:              "valid input",
			inputBody:         `{"username":"Nic", "password":"123456"}`,
			expectedParseBody: &models.AuthRequest{Username: "Nic", Password: "123456"},
			mockBehaviour: func(ctx context.Context, s *mock_auth.MockAuthUseCase, request *models.AuthRequest) {
				s.EXPECT().Authenticate(gomock.Any(), request).Return(&models.AuthResponse{Token: "abc"}, nil)
			},
			expectedCode: 200,
			expectedBody: `{"token":"abc"}`,
		},
		{
			name:              "invalid input",
			inputBody:         "",
			expectedParseBody: nil,
			mockBehaviour:     func(ctx context.Context, s *mock_auth.MockAuthUseCase, request *models.AuthRequest) {},
			expectedCode:      400,
			expectedBody:      `{"errors":"bad request"}`,
		},
		{
			name:              "invalid password",
			inputBody:         `{"username":"Nic", "password":"-1"}`,
			expectedParseBody: &models.AuthRequest{Username: "Nic", Password: "-1"},
			mockBehaviour: func(ctx context.Context, s *mock_auth.MockAuthUseCase, request *models.AuthRequest) {
				s.EXPECT().Authenticate(gomock.Any(), request).Return(nil, &errors.InvalidCredentialsError{})
			},
			expectedCode: 401,
			expectedBody: `{"errors":"invalid credentials"}`,
		},
		{
			name:              "null username",
			inputBody:         `{"username":"", "password":"123456"}`,
			expectedParseBody: &models.AuthRequest{Username: "", Password: "123456"},
			mockBehaviour:     func(ctx context.Context, s *mock_auth.MockAuthUseCase, request *models.AuthRequest) {},
			expectedCode:      400,
			expectedBody:      `{"errors":"bad request"}`,
		},
		{
			name:              "null password",
			inputBody:         `{"username":"Nic", "password":""}`,
			expectedParseBody: &models.AuthRequest{Username: "Nic", Password: ""},
			mockBehaviour:     func(ctx context.Context, s *mock_auth.MockAuthUseCase, request *models.AuthRequest) {},
			expectedCode:      400,
			expectedBody:      `{"errors":"bad request"}`,
		},
		{
			name:              "internal server error",
			inputBody:         `{"username":"Nic", "password":"123456"}`,
			expectedParseBody: &models.AuthRequest{Username: "Nic", Password: "123456"},
			mockBehaviour: func(ctx context.Context, s *mock_auth.MockAuthUseCase, request *models.AuthRequest) {
				s.EXPECT().Authenticate(gomock.Any(), request).Return(nil, &errors.NilPointerError{})
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

			auth := mock_auth.NewMockAuthUseCase(c)
			tc.mockBehaviour(ctx, auth, tc.expectedParseBody)
			lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))

			handler := NewAuthHandler(auth, lgr)

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
