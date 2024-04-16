package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Env        string `env:"REST_SERVER_ENV" env-required:"true"`
	AppID      int32  `env:"REST_SERVER_APP_ID" env-required:"true"`
	LogsPath   string `env:"REST_SERVER_LOGS_PATH_FILE" env-required:"true"`
	HTTPServer HTTPServer
	Clients    Clients
	Prometheus PrometheusConfig
}

func MustLoad() *Config {
	var cfg Config
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("unable to load .env file: %v", err)
		}
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error parsing environment variables: %v", err)
	}
	return &cfg
}
