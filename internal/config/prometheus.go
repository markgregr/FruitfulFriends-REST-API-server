package config

type PrometheusConfig struct {
	Host string `env:"REST_SERVER_PROMETHEUS_HOST" env-required:"true"`
	Port int    `env:"REST_SERVER_PROMETHEUS_PORT" env-required:"true"`
}
