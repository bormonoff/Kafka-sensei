package main

import (
	"flag"
	"fmt"

	"github.com/lovoo/goka"
	"github.com/rs/zerolog/log"

	"stream/internal/codec"
	"stream/internal/config"
	"stream/internal/model"
)

var (
	user      = flag.String("user", "", "user")
	blockUser = flag.String("block", "", "user to block")
)

func main() {
	if err := runEmitter(); err != nil {
		log.Fatal().Err(err).Msg("unhandeled run error")
	}
}

func runEmitter() error {
	flag.Parse()
	if *user == "" {
		return fmt.Errorf("args: can't ban empty user")
	}

	cfg, err := config.ParseFromEnv()
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	emitter, err := goka.NewEmitter(cfg.Kafka.Brokers, goka.Stream(cfg.Kafka.BlockedUsersStream), codec.NewBlockEventCodec())
	if err != nil {
		return fmt.Errorf("create emitter: %w", err)
	}
	defer emitter.Finish()

	ev := model.UserBlock{
		BlockUser: *blockUser,
	}

	err = emitter.EmitSync(*user, ev)
	if err != nil {
		return fmt.Errorf("emit sync: %w", err)
	}
	log.Info().Msg(fmt.Sprintf("%s blocked: %s", *user, ev.BlockUser))
	return nil
}
