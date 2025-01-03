package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/lovoo/goka"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"stream/internal/config"
	"stream/internal/processor"
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

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return processor.RunFilter(ctx, cfg.Kafka.Brokers, goka.Stream(cfg.Kafka.MessagesStream), goka.Stream(cfg.Kafka.FilteredMessagesStream))
	})

	g.Go(func() error {
		return processor.RunBlocker(ctx, cfg.Kafka.Brokers, goka.Stream(cfg.Kafka.BlockedUsersStream))
	})

	g.Go(func() error {
		return processor.RunClientSender(ctx, cfg.Kafka.Brokers, goka.Stream(cfg.Kafka.FilteredMessagesStream), cfg.ClientName)
	})

	return g.Wait()
}
