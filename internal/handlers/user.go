package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/pvpender/avito-shop/internal/usecase"
)

type UserHandler struct {
	userUS  *usecase.UserUseCase
	jwtAuth *jwtauth.JWTAuth
	logger  *slog.Logger
}

func NewUserHandler(us *usecase.UserUseCase, jwtAuth *jwtauth.JWTAuth, l *slog.Logger) *UserHandler {
	return &UserHandler{userUS: us, jwtAuth: jwtAuth, logger: l}
}

func (uh *UserHandler) Info(w http.ResponseWriter, r *http.Request) {
	uh.logger.Info("Info called")

	userId, err := getUserIdFromJwt(r.Context(), w, uh.logger, "UserHandler")
	if err != nil {
		return
	}

	info, err := uh.userUS.GetInfo(r.Context(), userId)
	if err != nil {
		respondWithError(w, uh.logger, http.StatusInternalServerError, "UserHandler", err)
		return
	}

	respondWithJSON(w, info, uh.logger, "UserHandler")
}
