//go:build amd64 || arm64
// +build amd64 arm64

// Package codec provides encoder/decoder implementations for various formats.
package codec

import (
	"github.com/bytedance/sonic"
)

const (
	// ContentTypeJSON is the MIME type for JSON.
	ContentTypeJSON = "application/json"
)

// SonicCodec implements Codec using bytedance/sonic for high-performance JSON.
// Sonic is significantly faster than encoding/json on amd64 and arm64 platforms.
type SonicCodec struct {
	api sonic.API
}

// JSONCodec is an alias for SonicCodec for backward compatibility.
type JSONCodec = SonicCodec

// NewSonicCodec creates a new Sonic JSON codec with default settings.
func NewSonicCodec() *SonicCodec {
	return &SonicCodec{
		api: sonic.ConfigDefault,
	}
}

// NewJSONCodec creates a new JSON codec (alias for NewSonicCodec).
func NewJSONCodec() *JSONCodec {
	return NewSonicCodec()
}

// NewSonicCodecFastest creates a Sonic codec optimized for maximum speed.
// Note: This may sacrifice some compatibility for performance.
func NewSonicCodecFastest() *SonicCodec {
	return &SonicCodec{
		api: sonic.ConfigFastest,
	}
}

// NewSonicCodecStd creates a Sonic codec with standard library compatibility.
func NewSonicCodecStd() *SonicCodec {
	return &SonicCodec{
		api: sonic.ConfigStd,
	}
}

// Encode encodes the value to JSON bytes using Sonic.
func (c *SonicCodec) Encode(v any) ([]byte, error) {
	return c.api.Marshal(v)
}

// Decode decodes JSON bytes into the destination using Sonic.
func (c *SonicCodec) Decode(data []byte, v any) error {
	return c.api.Unmarshal(data, v)
}

// ContentType returns the JSON MIME type.
func (c *SonicCodec) ContentType() string {
	return ContentTypeJSON
}

// Marshal is a convenience function for JSON encoding using Sonic.
func Marshal(v any) ([]byte, error) {
	return sonic.Marshal(v)
}

// Unmarshal is a convenience function for JSON decoding using Sonic.
func Unmarshal(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}

// Codec instances
var (
	// Sonic is the default Sonic codec instance.
	Sonic = NewSonicCodec()
	// SonicFastest is the fastest Sonic codec instance.
	SonicFastest = NewSonicCodecFastest()
	// JSON is the default JSON codec instance (alias for Sonic).
	JSON = Sonic
)

func init() {
	// Register JSON codec with Sonic for better performance
	Register(ContentTypeJSON, Sonic)
	Register("text/json", Sonic)
	Register("application/json; charset=utf-8", Sonic)
}
