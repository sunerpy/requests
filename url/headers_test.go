package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeaders(t *testing.T) {
	t.Run("Initialize Headers", func(t *testing.T) {
		headers := NewHeaders()
		assert.NotNil(t, headers)
		assert.NotNil(t, headers.headers)
		assert.Equal(t, 0, len(headers.headers))
	})
}

func TestHeaders_Add(t *testing.T) {
	t.Run("Add single header", func(t *testing.T) {
		headers := NewHeaders()
		headers.Add("Content-Type", "application/json")
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
	})
	t.Run("Add multiple headers with the same key", func(t *testing.T) {
		headers := NewHeaders()
		headers.Add("Accept", "text/html")
		headers.Add("Accept", "application/json")
		expected := []string{"text/html", "application/json"}
		assert.Equal(t, expected, headers.headers["Accept"])
	})
}

func TestHeaders_Set(t *testing.T) {
	t.Run("Set a header", func(t *testing.T) {
		headers := NewHeaders()
		headers.Set("Authorization", "Bearer token")
		assert.Equal(t, "Bearer token", headers.Get("Authorization"))
	})
	t.Run("Overwrite existing header", func(t *testing.T) {
		headers := NewHeaders()
		headers.Set("Authorization", "Bearer old-token")
		headers.Set("Authorization", "Bearer new-token")
		assert.Equal(t, "Bearer new-token", headers.Get("Authorization"))
	})
}

func TestHeaders_Get(t *testing.T) {
	t.Run("Get an existing header", func(t *testing.T) {
		headers := NewHeaders()
		headers.Set("User-Agent", "GoClient/1.0")
		assert.Equal(t, "GoClient/1.0", headers.Get("User-Agent"))
	})
	t.Run("Get a non-existing header", func(t *testing.T) {
		headers := NewHeaders()
		assert.Empty(t, headers.Get("Non-Existing"))
	})
}
