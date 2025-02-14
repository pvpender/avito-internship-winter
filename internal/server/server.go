package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/config"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/handlers"
	"github.com/pvpender/avito-shop/internal/middleware"
	"github.com/pvpender/avito-shop/internal/repositories"
	"github.com/pvpender/avito-shop/internal/usecase"
)

const (
	GrasefultShutdownTimeOut = 5
	ServerTimeOut            = 3
)

type Server struct {
	config    *config.Config
	jwtAuth   *jwtauth.JWTAuth
	db        *pgxpool.Pool
	trManager *manager.Manager
	builder   *squirrel.StatementBuilderType
	logger    *slog.Logger
}

func NewServer(config *config.Config, jwtAuth *jwtauth.JWTAuth, db *pgxpool.Pool, trManager *manager.Manager, builder *squirrel.StatementBuilderType, logger *slog.Logger) *Server {
	return &Server{config, jwtAuth, db, trManager, builder, logger}
}

func (server *Server) Run() error {
	r := chi.NewRouter()
	_ = server.PrepareHandlers(r)

	go func() {
		server.logger.With(
			slog.String("port", server.config.Server.Port),
		).Info("Server running on port")

		serv := &http.Server{
			Addr:              server.config.Server.Port,
			Handler:           r,
			ReadHeaderTimeout: ServerTimeOut * time.Second,
		}
		if err := serv.ListenAndServe(); err != nil {
			server.logger.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), GrasefultShutdownTimeOut*time.Second)
	defer shutdown()

	<-ctx.Done()
	server.logger.Warn("Server shutting down")

	return &errors.ShutdownError{}
}

func (server *Server) PrepareHandlers(r *chi.Mux) error {
	coinRepo := repositories.NewPgCoinRepository(server.db, trmpgx.DefaultCtxGetter, server.builder)
	itemRepo := repositories.NewPgItemRepository(server.db, trmpgx.DefaultCtxGetter, server.builder)
	purchaseRepo := repositories.NewPgPurchaseRepository(server.db, trmpgx.DefaultCtxGetter, server.builder)
	userRepo := repositories.NewPgUserRepository(server.db, trmpgx.DefaultCtxGetter, server.builder)

	authUS := usecase.NewAuthUseCase(server.jwtAuth, userRepo, server.logger)
	coinUS := usecase.NewCoinUseCase(server.trManager, userRepo, coinRepo)
	purchaseUS := usecase.NewPurchaseUseCase(server.trManager, purchaseRepo, userRepo, itemRepo)
	userUS := usecase.NewUserUseCase(userRepo, purchaseRepo, coinRepo)

	authHandler := handlers.NewAuthHandler(authUS, server.logger)
	userHandler := handlers.NewUserHandler(userUS, server.jwtAuth, server.logger)
	purchaseHandler := handlers.NewPurchaseHandler(purchaseUS, server.jwtAuth, server.logger)
	coinHandler := handlers.NewCoinHandler(coinUS, server.jwtAuth, server.logger)

	r.Group(func(r chi.Router) {
		r.Post("/api/auth", authHandler.Auth)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(server.jwtAuth))
		r.Use(middleware.Authenticator(server.jwtAuth))

		r.Get("/api/info", userHandler.Info)
		r.Get("/api/buy/{item}", purchaseHandler.Purchase)
		r.Post("/api/sendCoin", coinHandler.SendCoin)
	})

	return nil
}
