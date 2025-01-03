package processor

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/lovoo/goka"
	"github.com/rs/zerolog/log"

	"stream/internal/codec"
	"stream/internal/model"
)

var clientName string

func send(ctx goka.Context, msg any) {
	m, ok := msg.(model.Message)
	if !ok {
		log.Warn().Interface("value", m).Msg("type assertion failed, skip an event")
		return
	}

	key := ctx.Key()
	if key == "" {
		log.Warn().Interface("value", m).Int64("offset", ctx.Offset()).Msg("key is empty, skip an event")
		return
	}

	sender := ctx.Key()
	if sender == clientName {
		log.Debug().Str("message", m.Msg).Msg("skip a message because client is the sender")
		return
	}

	if bu := ctx.Lookup(goka.GroupTable(blockGroup), clientName); bu != nil {
		blockedUsers, ok := bu.(model.BlockList)
		if !ok {
			log.Warn().Interface("value", m).Msg("type assertion failed for ban-table stream, skip an event")
		}

		for _, val := range blockedUsers.BlockList {
			if val == sender {
				log.Debug().Str("message", m.Msg).Msg("skip an event for a blocked user")
				return
			}
		}
	}

	log.Info().Str("message", m.Msg).Msg(fmt.Sprintf("%s recieved a message from %s", clientName, sender))
}

func RunClientSender(ctx context.Context, brokers []string, inputTopic goka.Stream, client string) error {
	clientName = client
	g := goka.DefineGroup(goka.Group(client),
		goka.Input(inputTopic, codec.NewMessageCodec(), send),
		goka.Lookup(goka.GroupTable(blockGroup), codec.NewBlockListCodec()),
	)

	cfg := goka.DefaultConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	p, err := goka.NewProcessor(
		brokers,
		g,
		goka.WithConsumerGroupBuilder(
			goka.ConsumerGroupBuilderWithConfig(cfg),
		),
	)

	if err != nil {
		return fmt.Errorf("create client processor: %w", err)
	}
	err = p.Run(ctx)
	if err != nil {
		return fmt.Errorf("run client processor: %w", err)
	}

	return nil
}
