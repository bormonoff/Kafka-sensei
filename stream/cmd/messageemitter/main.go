package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lovoo/goka"
	"github.com/rs/zerolog/log"

	"stream/internal/codec"
	"stream/internal/config"
	"stream/internal/model"
)

func main() {
	if err := runEmitter(); err != nil {
		log.Fatal().Err(err).Msg("unhandeled run error")
	}
}

func runEmitter() error {
	cfg, err := config.ParseFromEnv()
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	emitter, err := goka.NewEmitter(cfg.Kafka.Brokers, goka.Stream(cfg.Kafka.MessagesStream), codec.NewMessageCodec())
	if err != nil {
		return fmt.Errorf("create emitter: %w", err)
	}
	defer emitter.Finish()

	err = start(emitter)
	if err != nil {
		return fmt.Errorf("start emit: %w", err)
	}
	return nil
}

var (
	str1 = []string{"Hi", "Hello", "Sup"}
	str2 = []string{"dear", "my", "lovely"}
	str3 = []string{"soulmate", "bro", "bormon_off"}
)

func start(emitter *goka.Emitter) error {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for range t.C {
		userId := fmt.Sprintf("user-%d", rand.Intn(5))

		fakeMsg := model.Message{
			Msg:       fmt.Sprintf("%s %s %s", str1[rand.Intn(3)], str2[rand.Intn(3)], str3[rand.Intn(3)]),
			CreatedAt: time.Now(),
		}

		err := emitter.EmitSync(userId, fakeMsg)
		if err != nil {
			return fmt.Errorf("emit sync: %w", err)
		}
		log.Info().Msg(fmt.Sprintf("%s writes: %s", userId, fakeMsg.Msg))
	}

	return nil
}
