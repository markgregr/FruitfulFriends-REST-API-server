package config

type KafkaClient struct {
	Broker   string `env:"REST_SERVER_KAFKA_BROKER" env-required:"true"`
	GroupID  string `env:"REST_SERVER_KAFKA_GROUP_ID" env-required:"true"`
	Topic    string `env:"REST_SERVER_KAFKA_TOPIC" env-required:"true"`
	Username string `env:"REST_SERVER_KAFKA_USERNAME" env-required:"true"`
	Password string `env:"REST_SERVER_KAFKA_PASSWORD" env-required:"true"`
}
