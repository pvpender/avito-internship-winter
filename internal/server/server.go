package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pvpender/avito-shop/config"
	"github.com/pvpender/avito-shop/internal/server/errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	config *config.Config
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewServer(config *config.Config, db *pgxpool.Pool, logger *slog.Logger) *Server {
	return &Server{config, db, logger}
}

func (server *Server) Run() error {
	r := chi.NewRouter()

	go func() {
		server.logger.With(
			slog.String("port", server.config.Server.Port),
		).Info("Server running on port")
		if err := http.ListenAndServe(server.config.Server.Port, r); err != nil {
			server.logger.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	select {
	case <-ctx.Done():
		server.logger.Warn("Server shutting down")
	}

	return &errors.ShutdownError{}
}
