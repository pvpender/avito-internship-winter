package main

import (
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

	s := server.NewServer(cfg, pgDB, lgr)
	if err := s.Run(); err != nil {
		lgr.Error(err.Error())
	}
}
