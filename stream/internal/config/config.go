package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" required:"true"`

	ClientName string `envconfig:"CLIENT_NAME" required:"true"`

	Kafka KafkaConfig
}

type KafkaConfig struct {
	Brokers                []string `envconfig:"BROKERS" required:"true"`
	MessagesStream         string   `envconfig:"MESSAGES_TOPIC" required:"true"`
	FilteredMessagesStream string   `envconfig:"FILTERED_MESSAGES_TOPIC" required:"true"`
	BlockedUsersStream     string   `envconfig:"BLOCKED_USERS_TOPIC" required:"true"`
	BlockedUsersTable      string   `envconfig:"BLOCKED_USERS_TABLE" required:"true"`
}

func ParseFromEnv() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("process environment variables: %w", err)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.LogLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return cfg, nil
}
