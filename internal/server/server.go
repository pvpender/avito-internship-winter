package server

import (
	"context"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/config"
	_ "github.com/pvpender/avito-shop/docs"
	"github.com/pvpender/avito-shop/internal/errors"
	"github.com/pvpender/avito-shop/internal/handlers"
	"github.com/pvpender/avito-shop/internal/middleware"
	"github.com/pvpender/avito-shop/internal/repositories"
	"github.com/pvpender/avito-shop/internal/usecase"
)

const (
	GracefulShutdownTimeOut = 5
	ServerTimeOut           = 3
	ServerMaxAge            = 300
)

type Server struct {
	config    *config.Config
	jwtAuth   *jwtauth.JWTAuth
	db        *pgxpool.Pool
	trManager *manager.Manager
	builder   *squirrel.StatementBuilderType
	logger    *slog.Logger
}

func NewServer(
	config *config.Config,
	jwtAuth *jwtauth.JWTAuth,
	db *pgxpool.Pool,
	trManager *manager.Manager,
	builder *squirrel.StatementBuilderType,
	logger *slog.Logger,
) *Server {
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

	ctx, shutdown := context.WithTimeout(context.Background(), GracefulShutdownTimeOut*time.Second)
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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           ServerMaxAge,
	}))

	r.Use(middleware2.Recoverer)
	r.Use(middleware2.AllowContentType("application/json"))
	r.Use(middleware2.RedirectSlashes)
	r.Use(middleware2.RequestID)
	r.Use(middleware2.CleanPath)
	r.Use(middleware2.NoCache)

	r.Group(func(r chi.Router) {
		r.Post("/api/auth", authHandler.Auth)
		r.Get("/swagger/*", httpSwagger.WrapHandler)
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
