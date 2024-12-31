package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Brokers string `envconfig:"BROKERS" required:"true"`
	Topic   string `envconfig:"TOPIC" required:"true"`
	Acks    int    `envconfig:"ACKS" default:"-1"`
}

func ParseFromEnv() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("process environment variables: %w", err)
	}

	return cfg, nil
}
