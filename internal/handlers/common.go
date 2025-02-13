package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
	"net/http"
)

func respondWithError(w http.ResponseWriter, l *slog.Logger, statusCode int, handlerName string, err error) {
	l.Error(err.Error(), "handler", handlerName)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	var errMessage string
	switch statusCode {
	case 400:
		errMessage = `{"errors": "bad request"}`
	case 401:
		errMessage = `{"errors": "invalid credentials"}`
	case 403:
		errMessage = `{"errors": "forbidden"}`
	default:
		errMessage = `{"errors": "internal server error"}`

	}
	_, _ = w.Write([]byte(errMessage))
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
		respondWithError(w, logger, http.StatusInternalServerError, "UserHandler", err)
		return 0, err
	}

	userId, ok := claims["user_id"].(float64)

	if !ok {
		respondWithError(w, logger, http.StatusInternalServerError, "UserHandler", errors.New("invalid jwt"))
		return 0, errors.New("invalid jwt")
	}

	return uint32(userId), nil
}
