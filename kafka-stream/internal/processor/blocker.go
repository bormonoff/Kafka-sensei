package processor

import (
	"context"
	"fmt"
	"slices"

	"github.com/IBM/sarama"
	"github.com/lovoo/goka"
	"github.com/rs/zerolog/log"

	"stream/internal/codec"
	"stream/internal/model"
)

var (
	blockGroup goka.Group = "block-users-consumer"
)

func block(ctx goka.Context, msg any) {
	blockEvent, ok := msg.(model.UserBlock)
	if !ok {
		log.Warn().Interface("value", blockEvent).Msg("type assertion failed, skip an event")
		return
	}

	key := ctx.Key()
	if key == "" {
		log.Warn().Interface("value", blockEvent).Int64("offset", ctx.Offset()).Msg("key is empty, skip an event")
		return
	}

	var userBlockList model.BlockList
	if val := ctx.Value(); val != nil {
		if userBlockList, ok = val.(model.BlockList); !ok {
			log.Warn().Interface("value", val).Msg("type assertion failed, unexpected event in the block table, skip an event")
			return
		}
	} else {
		userBlockList = model.BlockList{BlockList: make([]string, 0, 1)}
	}

	if !slices.Contains(userBlockList.BlockList, blockEvent.BlockUser) {
		userBlockList.BlockList = append(userBlockList.BlockList, blockEvent.BlockUser)
		ctx.SetValue(userBlockList)
		log.Debug().Str("user", key).Str("blocked", blockEvent.BlockUser).Msg("block handler processor: success block table update")
	}
}

func RunBlocker(ctx context.Context, brokers []string, inputTopic goka.Stream) error {
	g := goka.DefineGroup(blockGroup,
		goka.Input(inputTopic, codec.NewBlockEventCodec(), block),
		goka.Persist(codec.NewBlockListCodec()),
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
		return fmt.Errorf("create blocker handler processor: %w", err)
	}
	err = p.Run(ctx)
	if err != nil {
		return fmt.Errorf("run blocker handler processor: %w", err)
	}

	return nil
}
