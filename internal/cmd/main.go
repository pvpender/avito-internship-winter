package main

import (
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/jwtauth/v5"
	"github.com/pvpender/avito-shop/config"
	"github.com/pvpender/avito-shop/database"
	server "github.com/pvpender/avito-shop/internal/server"
	"log/slog"
	"os"
)

const (
	configFile = "config.yaml"
	configType = "yaml"
)

func main() {
	lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	cfg, err := config.LoadConfig(configFile, configType)
	if err != nil {
		lgr.Error(err.Error())
	}

	pgDB, err := database.NewPgPool(cfg)
	if err != nil {
		lgr.Error(err.Error())
	}
	defer pgDB.Close()
	trManager := manager.Must(trmpgx.NewDefaultFactory(pgDB))
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	tokenAuth := jwtauth.New("HS256", []byte(cfg.Auth.Secret), nil)

	s := server.NewServer(cfg, tokenAuth, pgDB, trManager, &builder, lgr)
	if sErr := s.Run(); sErr != nil {
		lgr.Error(sErr.Error())
	}
}
