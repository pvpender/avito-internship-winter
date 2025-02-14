package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
)

func respondWithError(w http.ResponseWriter, l *slog.Logger, statusCode int, handlerName string, err error) {
	l.Error(err.Error(), "handler", handlerName)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errMessage := &models.ErrorResponse{}

	switch statusCode {
	case http.StatusBadRequest:
		errMessage.Errors = "bad request"
	case http.StatusUnauthorized:
		errMessage.Errors = "invalid credentials"
	case http.StatusForbidden:
		errMessage.Errors = "forbidden"
	default:
		errMessage.Errors = "internal server error"
	}

	response, _ := json.Marshal(errMessage)
	_, _ = w.Write(response)
}

func respondWithJSON(w http.ResponseWriter, payload interface{}, logger *slog.Logger, handlerName string) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, logger, http.StatusInternalServerError, handlerName, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func getUserIdFromJwt(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, handlerName string) (uint32, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		respondWithError(w, logger, http.StatusInternalServerError, handlerName, err)
		return 0, err
	}

	userId, ok := claims["user_id"].(float64)

	if !ok {
		respondWithError(w, logger, http.StatusInternalServerError, handlerName, &errors.InvalidJWTError{})
		return 0, &errors.InvalidJWTError{}
	}

	return uint32(userId), nil
}