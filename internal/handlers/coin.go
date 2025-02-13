package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase"
	"log/slog"
	"net/http"
)

type CoinHandler struct {
	purchaseUS *usecase.CoinUseCase
	jwtAuth    *jwtauth.JWTAuth
	logger     *slog.Logger
}

func NewCoinHandler(purchaseUS *usecase.CoinUseCase, jwtAuth *jwtauth.JWTAuth, logger *slog.Logger) *CoinHandler {
	return &CoinHandler{purchaseUS: purchaseUS, jwtAuth: jwtAuth, logger: logger}
}

func (ch *CoinHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	ch.logger.Info("SendCoin called")
	userId, err := getUserIdFromJwt(r.Context(), w, ch.logger, "CoinHandler")
	if err != nil {
		respondWithError(w, ch.logger, http.StatusInternalServerError, "CoinHandler", err)
		return
	}

	var request *models.SendCoinRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil || !request.Validate() {
		respondWithError(w, ch.logger, http.StatusBadRequest, "SendCoin", err)
		return
	}

	err = ch.purchaseUS.SendCoin(r.Context(), userId, request)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, ch.logger, http.StatusBadRequest, "SendCoin", err)
			return
		}

		respondWithError(w, ch.logger, http.StatusInternalServerError, "SendCoin", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
