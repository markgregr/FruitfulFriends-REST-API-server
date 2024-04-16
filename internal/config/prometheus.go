package config

type PrometheusConfig struct {
	Listen string `env:"REST_SERVER_PROMETHEUS_LISTEN"`
}
