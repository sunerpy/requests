//go:build !amd64 && !arm64
// +build !amd64,!arm64

// Package codec provides encoder/decoder implementations for various formats.
package codec

import (
	"encoding/json"
)

const (
	// ContentTypeJSON is the MIME type for JSON.
	ContentTypeJSON = "application/json"
)

// SonicCodec falls back to encoding/json on unsupported platforms.
type (
	SonicCodec struct{}
	// JSONCodec is an alias for SonicCodec for backward compatibility.
	JSONCodec = SonicCodec
)

// NewSonicCodec creates a new codec (falls back to encoding/json).
func NewSonicCodec() *SonicCodec {
	return &SonicCodec{}
}

// NewJSONCodec creates a new JSON codec (alias for NewSonicCodec).
func NewJSONCodec() *JSONCodec {
	return NewSonicCodec()
}

// NewSonicCodecFastest creates a codec (falls back to encoding/json).
func NewSonicCodecFastest() *SonicCodec {
	return &SonicCodec{}
}

// NewSonicCodecStd creates a codec (falls back to encoding/json).
func NewSonicCodecStd() *SonicCodec {
	return &SonicCodec{}
}

// Encode encodes the value to JSON bytes.
func (c *SonicCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Decode decodes JSON bytes into the destination.
func (c *SonicCodec) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// ContentType returns the JSON MIME type.
func (c *SonicCodec) ContentType() string {
	return ContentTypeJSON
}

// Marshal is a convenience function for JSON encoding.
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal is a convenience function for JSON decoding.
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// Codec instances
var (
	// Sonic is the default codec instance.
	Sonic = NewSonicCodec()
	// SonicFastest is the fastest codec instance (same as Sonic on fallback).
	SonicFastest = NewSonicCodec()
	// JSON is the default JSON codec instance (alias for Sonic).
	JSON = Sonic
)

func init() {
	// Register JSON codec
	Register(ContentTypeJSON, Sonic)
	Register("text/json", Sonic)
	Register("application/json; charset=utf-8", Sonic)
}
