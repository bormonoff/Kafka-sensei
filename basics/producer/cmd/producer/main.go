package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"

	"producer/internal/config"
	"producer/internal/model"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("unhandeled run error")
	}
}

func run() error {
	cfg, err := config.ParseFromEnv()
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"acks":              cfg.Acks,
		"retries":           3,
	})
	if err != nil {
		return fmt.Errorf("create producer: %w", err)
	}

	msg, err := json.Marshal(model.Square{
		X: 10,
		Y: 20,
	})
	if err != nil {
		return fmt.Errorf("marshal square: %w", err)
	}

	err = publishMsg(p, cfg.Topic, msg)
	if err != nil {
		return fmt.Errorf("publish kafka msg: %w", err)
	}

	p.Close()
	log.Info().Str("message", string(msg)).Msg("message has been successfully sent")

	return nil
}

func publishMsg(p *kafka.Producer, topic string, msg []byte) error {
	deliveryChan := make(chan kafka.Event)
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("produce msg fail: %w", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery: %w", err)
	}

	close(deliveryChan)
	return nil
}
