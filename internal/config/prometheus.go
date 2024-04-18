package config

type Prometheus struct {
	Host string `env:"REST_SERVER_PROMETHEUS_HOST" env-default:"0.0.0.0"`
	Port int    `env:"REST_SERVER_PROMETHEUS_PORT" env-default:"8082"`
}
