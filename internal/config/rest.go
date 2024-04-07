package config

import "time"

type HTTPServer struct {
	Host        string        `env:"REST_SERVER_HOST" env-required:"true"`
	Port        int           `env:"REST_SERVER_PORT" env-required:"true"`
	Timeout     time.Duration `env:"REST_SERVER_TIMEOUT" env-required:"true"`
	IdleTimeout time.Duration `env:"REST_SERVER_IDLE_TIMEOUT" env-required:"true"`
	AllowOrigin string        `env:"REST_SERVER_ALLOW_ORIGIN" envDefault:"*"`
}
