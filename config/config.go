package config

import (
	"github.com/pvpender/avito-shop/config/errors"
	"os"
)

const (
	ServerPortAlias = "SERVER_PORT"
	DbHostAlias     = "DATABASE_HOST"
	DbPortAlias     = "DATABASE_PORT"
	DbNameAlias     = "DATABASE_NAME"
	DbUserAlias     = "DATABASE_USER"
	DbPasswordAlias = "DATABASE_PASSWORD"
)

type Config struct {
	HTTPServer
	PostgresConfig
}

type HTTPServer struct {
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func LoadConfig() (*Config, error) {
	port, exists := os.LookupEnv(ServerPortAlias)
	if exists != true {
		return nil, &errors.LoadEnvError{EnvName: ServerPortAlias}
	}

	dbHost, exists := os.LookupEnv(DbHostAlias)
	if exists != true {
		return nil, &errors.LoadEnvError{EnvName: DbHostAlias}
	}

	dbPort, exists := os.LookupEnv(DbPortAlias)
	if exists != true {
		return nil, &errors.LoadEnvError{EnvName: DbPortAlias}
	}

	dbName, exists := os.LookupEnv(DbNameAlias)
	if exists != true {
		return nil, &errors.LoadEnvError{EnvName: DbNameAlias}
	}

	dbUser, exists := os.LookupEnv(DbUserAlias)
	if exists != true {
		return nil, &errors.LoadEnvError{EnvName: DbUserAlias}
	}

	dbPassword, exists := os.LookupEnv(DbPasswordAlias)
	if exists != true {
		return nil, &errors.LoadEnvError{EnvName: DbPasswordAlias}
	}

	return &Config{
		HTTPServer{
			Port: port,
		},
		PostgresConfig{
			Host:     dbHost,
			Port:     dbPort,
			Database: dbName,
			Username: dbUser,
			Password: dbPassword,
		},
	}, nil
}
