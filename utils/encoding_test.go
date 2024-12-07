package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLEncode(t *testing.T) {
	t.Run("Encode URL", func(t *testing.T) {
		input := "hello world"
		expected := "hello+world"
		assert.Equal(t, expected, URLEncode(input))
	})
}

func TestURLDecode(t *testing.T) {
	t.Run("Decode URL", func(t *testing.T) {
		input := "hello+world"
		expected := "hello world"
		output, err := URLDecode(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("Invalid Decode", func(t *testing.T) {
		input := "%hello"
		_, err := URLDecode(input)
		assert.Error(t, err)
	})
}

func TestBase64Encode(t *testing.T) {
	t.Run("Encode Base64", func(t *testing.T) {
		input := "hello world"
		expected := "aGVsbG8gd29ybGQ="
		assert.Equal(t, expected, Base64Encode(input))
	})
}

func TestBase64Decode(t *testing.T) {
	t.Run("Decode Base64", func(t *testing.T) {
		input := "aGVsbG8gd29ybGQ="
		expected := []byte("hello world")
		output, err := Base64Decode(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("Invalid Decode", func(t *testing.T) {
		input := "aGVsbG8gd29ybGQ"
		_, err := Base64Decode(input)
		assert.Error(t, err)
	})
}

func TestHexEncode(t *testing.T) {
	t.Run("Encode Hex", func(t *testing.T) {
		input := "hello world"
		expected := "68656c6c6f20776f726c64"
		assert.Equal(t, expected, HexEncode(input))
	})
}

func TestHexDecode(t *testing.T) {
	t.Run("Decode Hex", func(t *testing.T) {
		input := "68656c6c6f20776f726c64"
		expected := []byte("hello world")
		output, err := HexDecode(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("Invalid Decode", func(t *testing.T) {
		input := "68656c6c6f20776f726c6g"
		_, err := HexDecode(input)
		assert.Error(t, err)
	})
}
