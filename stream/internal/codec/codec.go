package codec

type Codec interface {
	Encode(value any) (data []byte, err error)
	Decode(data []byte) (value any, err error)
}
