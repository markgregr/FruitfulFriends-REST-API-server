package config

import "time"

type GRPCClient struct {
	Address      string        `env:"REST_SERVER_GRPC_CLIENT_ADDRESS" env-required:"true"`
	Timeout      time.Duration `env:"REST_SERVER_GRPC_CLIENT_TIMEOUT"env-required:"true"`
	RetriesCount int           `env:"REST_SERVER_GRPC_CLIENT_RETRIES_COUNT" env-required:"true"`
	Insecure     bool          `env:"REST_SERVER_GRPC_CLIENT_INSECURE" env-default:"false"`
}
