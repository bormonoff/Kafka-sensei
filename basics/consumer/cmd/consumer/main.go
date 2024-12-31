package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"

	"consumer/internal/config"
	"consumer/internal/model"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("unhandeled run error")
	}
}

func run() error {
	cfg := config.ParseArgs()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.KafkaConsumer.Broker,
		"group.id":           cfg.KafkaConsumer.ConsumerGroup,
		"enable.auto.commit": cfg.KafkaConsumer.EnableAutoCommit,
		"auto.offset.reset":  "earliest",
		"session.timeout.ms": 6000,
	})
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	err = c.SubscribeTopics([]string{cfg.KafkaConsumer.Topic}, nil)
	if err != nil {
		return fmt.Errorf("subscribe on topics: %w", err)
	}

	log.Info().Msg("consumer has been started...")
	err = pollEvents(done, c, cfg.KafkaConsumer.EnableAutoCommit, cfg.KafkaConsumer.PollPeriod)
	if err != nil {
		return fmt.Errorf("poll event: %w", err)
	}

	c.Close()

	return nil
}

func pollEvents(done chan (os.Signal), c *kafka.Consumer, enableAutoCommit bool, pollPeriod int) error {
	for {
		select {
		case <-done:
			return nil
		default:
			ev := c.Poll(pollPeriod)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				value := model.Square{}
				err := json.Unmarshal(e.Value, &value)
				if err != nil {
					return fmt.Errorf("unmarshal kafka event: %w", err)
				} else {
					log.Info().Interface("data", value).Msg("get event")
				}

			case kafka.Error:
				return fmt.Errorf("unexpected kafka error code: %s", e.Code().String())
			}

			if !enableAutoCommit {
				c.CommitMessage(ev.(*kafka.Message))
			}
		}
	}
}
