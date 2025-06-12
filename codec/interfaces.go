// Package codec provides encoder/decoder interfaces and implementations.
package codec

// Encoder defines the interface for encoding data.
type Encoder interface {
	// Encode encodes the given value to bytes.
	Encode(v any) ([]byte, error)
	// ContentType returns the MIME type for this encoder.
	ContentType() string
}

// Decoder defines the interface for decoding data.
type Decoder interface {
	// Decode decodes the given bytes into the destination.
	Decode(data []byte, v any) error
}

// Codec combines Encoder and Decoder interfaces.
type Codec interface {
	Encoder
	Decoder
}

// EncoderFunc is a function adapter for Encoder interface.
type EncoderFunc struct {
	EncodeFunc      func(v any) ([]byte, error)
	ContentTypeFunc func() string
}

// Encode implements Encoder.
func (f EncoderFunc) Encode(v any) ([]byte, error) {
	return f.EncodeFunc(v)
}

// ContentType implements Encoder.
func (f EncoderFunc) ContentType() string {
	return f.ContentTypeFunc()
}

// DecoderFunc is a function adapter for Decoder interface.
type DecoderFunc func(data []byte, v any) error

// Decode implements Decoder.
func (f DecoderFunc) Decode(data []byte, v any) error {
	return f(data, v)
}
