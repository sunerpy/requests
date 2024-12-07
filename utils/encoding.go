package utils

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

func URLEncode(value string) string {
	return url.QueryEscape(value)
}

func URLDecode(value string) (string, error) {
	return url.QueryUnescape(value)
}

func Base64Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func Base64Decode(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

func HexEncode(value string) string {
	return hex.EncodeToString([]byte(value))
}

func HexDecode(value string) ([]byte, error) {
	return hex.DecodeString(value)
}
