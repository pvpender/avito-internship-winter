package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	errInt "github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/usecase"
)

type PurchaseHandler struct {
	purchaseUS *usecase.PurchaseUseCase
	jwtAuth    *jwtauth.JWTAuth
	logger     *slog.Logger
}

func NewPurchaseHandler(
	purchaseUS *usecase.PurchaseUseCase,
	jwtAuth *jwtauth.JWTAuth,
	logger *slog.Logger,
) *PurchaseHandler {
	return &PurchaseHandler{purchaseUS: purchaseUS, jwtAuth: jwtAuth, logger: logger}
}

func (ph *PurchaseHandler) Purchase(w http.ResponseWriter, r *http.Request) {
	ph.logger.Info("Purchase called")

	userId, err := getUserIdFromJwt(r.Context(), w, ph.logger, "PurchaseHandler")
	if err != nil {
		return
	}

	itemType := chi.URLParam(r, "item")

	err = ph.purchaseUS.CreatePurchase(r.Context(), userId, itemType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, &errInt.PurchaseError{}) {
			respondWithError(w, ph.logger, http.StatusBadRequest, "PurchaseHandler", err)
			return
		}

		respondWithError(w, ph.logger, http.StatusInternalServerError, "PurchaseHandler", err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
