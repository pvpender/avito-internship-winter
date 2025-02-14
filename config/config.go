package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   HTTPServer
	Postgres PostgresConfig
	Auth     AuthConfig
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

type AuthConfig struct {
	Secret string
}

func LoadConfig(filename string, configType string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType(configType)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
