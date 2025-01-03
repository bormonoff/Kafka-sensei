package consumer

import (
	"flag"

	"github.com/rs/zerolog/log"
)

type Config struct {
	KafkaConsumer KafkaConsumer
}

type KafkaConsumer struct {
	Broker           string
	Topic            string
	ConsumerGroup    string
	EnableAutoCommit bool
	PollPeriod       int
}

func ParseArgs() Config {
	var cfg Config

	flag.StringVar(&cfg.KafkaConsumer.Broker, "broker", "", "kafka consumer broker addresses url")
	flag.StringVar(&cfg.KafkaConsumer.Topic, "topic", "", "kafka consumer topic name")
	flag.StringVar(&cfg.KafkaConsumer.ConsumerGroup, "group", "", "kafka consumer group name")
	flag.IntVar(&cfg.KafkaConsumer.PollPeriod, "poll-period", 100, "kafka poll period. 100 by default.")
	flag.BoolVar(&cfg.KafkaConsumer.EnableAutoCommit, "auto-commit", true, "enable/disable auto commit. Enable by default.")

	flag.Parse()

	log.Info().
		Str("brokers", cfg.KafkaConsumer.Broker).
		Str("kafka_topic", cfg.KafkaConsumer.Topic).
		Str("kafka_consumer_group", cfg.KafkaConsumer.ConsumerGroup).
		Bool("enable_auto_commit", cfg.KafkaConsumer.EnableAutoCommit).
		Int("poll_period", cfg.KafkaConsumer.PollPeriod).Msg("parse cli args")

	return cfg
}
