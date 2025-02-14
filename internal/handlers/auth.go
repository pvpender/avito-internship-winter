package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errInt "github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase"
)

type AuthHandler struct {
	authUS *usecase.AuthUseCase
	logger *slog.Logger
}

func NewAuthHandler(authUS *usecase.AuthUseCase, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{authUS: authUS, logger: logger}
}

func (handler *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Auth called")

	var request *models.AuthRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || !request.Validate() {
		respondWithError(w, handler.logger, http.StatusBadRequest, "AuthHandler", err)
		return
	}

	response, err := handler.authUS.Authenticate(r.Context(), request)
	if err != nil {
		if errors.Is(err, &errInt.InvalidCredentialsError{}) {
			respondWithError(w, handler.logger, http.StatusUnauthorized, "AuthHandler", err)
			return
		}

		respondWithError(w, handler.logger, http.StatusInternalServerError, "AuthHandler", err)

		return
	}

	respondWithJSON(w, response, handler.logger, "AuthHandler")
}
