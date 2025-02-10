package main

import (
	"github.com/pvpender/avito-shop/config"
	server "github.com/pvpender/avito-shop/internal/server"
	"log/slog"
	"os"
)

func main() {
	lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	cfg, err := config.LoadConfig()
	if err != nil {
		lgr.Error(err.Error())
	}

	s := server.NewServer(cfg, lgr)
	if err := s.Run(); err != nil {
		lgr.Error(err.Error())
	}
}
