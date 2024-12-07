package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	t.Run("Valid GET Request", func(t *testing.T) {
		method := "GET"
		rawURL := "https://example.com"
		req, err := NewRequest(method, rawURL)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, method, req.Method)
		assert.Equal(t, rawURL, req.URL.String())
	})
	t.Run("Valid POST Request", func(t *testing.T) {
		method := "POST"
		rawURL := "https://example.com/api"
		req, err := NewRequest(method, rawURL)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, method, req.Method)
		assert.Equal(t, rawURL, req.URL.String())
	})
	t.Run("Invalid URL", func(t *testing.T) {
		method := "GET"
		rawURL := "://invalid-url"
		req, err := NewRequest(method, rawURL)
		assert.Error(t, err)
		assert.Nil(t, req)
	})
}
