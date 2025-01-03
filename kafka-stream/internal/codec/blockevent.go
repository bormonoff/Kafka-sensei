package codec

import (
	"encoding/json"
	"fmt"

	"stream/internal/model"
)

type blockEventCodec struct{}

func NewBlockEventCodec() Codec {
	return blockEventCodec{}
}

func (c blockEventCodec) Encode(value any) (data []byte, err error) {
	if _, ok := value.(model.UserBlock); !ok {
		return nil, fmt.Errorf("blockEventCodec: invalid encode type: %v", value)
	}

	msg, err := json.Marshal(value)
	if err != nil {
		return msg, fmt.Errorf("blockEventCodec marshal: %w", err)
	}

	return msg, nil
}

func (c blockEventCodec) Decode(data []byte) (value any, err error) {
	res := model.UserBlock{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
