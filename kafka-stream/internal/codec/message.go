package codec

import (
	"encoding/json"
	"fmt"

	"stream/internal/model"
)

type messageCodec struct{}

func NewMessageCodec() Codec {
	return messageCodec{}
}

func (c messageCodec) Encode(value any) (data []byte, err error) {
	if _, ok := value.(model.Message); !ok {
		return nil, fmt.Errorf("messageCodec: invalid encode type: %v", value)
	}

	msg, err := json.Marshal(value)
	if err != nil {
		return msg, fmt.Errorf("messageCodec marshal: %w", err)
	}

	return msg, nil
}

func (c messageCodec) Decode(data []byte) (value any, err error) {
	res := model.Message{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, fmt.Errorf("messageCodec unmasrshal: %w", err)
	}
	return res, nil
}
