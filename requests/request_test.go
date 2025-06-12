package requests

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	customurl "github.com/sunerpy/requests/url"
)

func TestNewRequest(t *testing.T) {
	t.Run("Valid URL with Params", func(t *testing.T) {
		params := customurl.NewValues()
		params.Set("query", "golang")
		req, err := NewRequest("GET", "https://example.com", params, nil)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "https://example.com?query=golang", req.URL.String())
		assert.Nil(t, req.Body)
		assert.Equal(t, params, req.Params)
	})
	t.Run("Valid URL without Params", func(t *testing.T) {
		req, err := NewRequest("POST", "https://example.com", nil, strings.NewReader("test body"))
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "https://example.com", req.URL.String())
		assert.NotNil(t, req.Body)
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "test body", string(body))
	})
	t.Run("Invalid URL", func(t *testing.T) {
		req, err := NewRequest("GET", "://invalid-url", nil, nil)
		assert.Error(t, err)
		assert.Nil(t, req)
	})
	t.Run("Invalid URL with ctx", func(t *testing.T) {
		req, err := NewRequestWithContext(context.Background(), "GET", "://invalid-url", nil, nil)
		assert.Error(t, err)
		assert.Nil(t, req)
	})
}

func TestRequest_AddHeader(t *testing.T) {
	req, _ := NewRequest("GET", "https://example.com", nil, nil)
	req.AddHeader("Authorization", "Bearer token")
	req.AddHeader("Content-Type", "application/json")
	assert.Equal(t, "Bearer token", req.Headers.Get("Authorization"))
	assert.Equal(t, "application/json", req.Headers.Get("Content-Type"))
}
