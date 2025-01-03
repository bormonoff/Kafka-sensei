package processor

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/lovoo/goka"
	"github.com/rs/zerolog/log"

	"stream/internal/codec"
	"stream/internal/model"
)

var (
	filterGroup goka.Group = "message-filter"
)

func filter(outputTopic goka.Stream) func(ctx goka.Context, msg any) {
	return func(ctx goka.Context, msg any) {
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

		newMsg, isCensored := censor(m.Msg)
		if isCensored {
			log.Debug().Str("user", key).Str("input message", m.Msg).Str("output message", newMsg).Msg("filter processor: success message censor")
		}

		m.Msg = newMsg
		ctx.Emit(outputTopic, key, m)
	}
}

func RunFilter(ctx context.Context, brokers []string, inputTopic goka.Stream, outputTopic goka.Stream) error {
	g := goka.DefineGroup(filterGroup,
		goka.Input(inputTopic, codec.NewMessageCodec(), filter(outputTopic)),
		goka.Output(outputTopic, codec.NewMessageCodec()),
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
		return fmt.Errorf("create filter processor: %w", err)
	}
	err = p.Run(ctx)
	if err != nil {
		return fmt.Errorf("run filter processor: %w", err)
	}
	return nil
}

var censorWords = []string{"Hi", "bro"}

func censor(msg string) (string, bool) {
	// The algorithm is for conceptual understanding only
	res := msg
	for _, val := range censorWords {
		res = strings.ReplaceAll(res, val, "***")
	}

	if res != msg {
		return res, true
	}

	return res, false
}
