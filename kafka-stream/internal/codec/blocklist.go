package codec

import (
	"encoding/json"
	"fmt"

	"stream/internal/model"
)

type blockListCodec struct{}

func NewBlockListCodec() Codec {
	return blockListCodec{}
}

func (c blockListCodec) Encode(value any) (data []byte, err error) {
	if _, ok := value.(model.BlockList); !ok {
		return nil, fmt.Errorf("blockListCodec: invalid encode type: %v", value)
	}

	msg, err := json.Marshal(value)
	if err != nil {
		return msg, fmt.Errorf("blockListCodec marshal: %w", err)
	}

	return msg, nil
}

func (c blockListCodec) Decode(data []byte) (value any, err error) {
	res := model.BlockList{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
