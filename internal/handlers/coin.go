package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	errInt "github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/models"
	"github.com/pvpender/avito-shop/internal/usecase/coin"
)

type CoinHandler struct {
	coinUS  coin.CoinUseCase
	jwtAuth *jwtauth.JWTAuth
	logger  *slog.Logger
}

func NewCoinHandler(coinUS coin.CoinUseCase, jwtAuth *jwtauth.JWTAuth, logger *slog.Logger) *CoinHandler {
	return &CoinHandler{coinUS: coinUS, jwtAuth: jwtAuth, logger: logger}
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

	err = ch.coinUS.SendCoin(r.Context(), userId, request)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, &errInt.InvalidAmountError{}) {
			respondWithError(w, ch.logger, http.StatusBadRequest, "SendCoin", err)
			return
		}

		respondWithError(w, ch.logger, http.StatusInternalServerError, "SendCoin", err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
