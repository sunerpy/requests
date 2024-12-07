package requests

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunerpy/requests/requests"
	"github.com/sunerpy/requests/url"
)

func TestGet(t *testing.T) {
	t.Run("Valid GET Request", func(t *testing.T) {
		params := url.NewURLParams()
		params.Set("key", "value")
		resp, err := Get("https://httpbin.org/get", params)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		var result map[string]interface{}
		err = json.Unmarshal(resp.Bytes(), &result)
		assert.NoError(t, err)
		assert.Contains(t, result["args"].(map[string]interface{}), "key")
	})
}

func TestPost(t *testing.T) {
	t.Run("Valid POST Request", func(t *testing.T) {
		form := url.NewForm()
		form.Set("key", "value")
		resp, err := Post("https://httpbin.org/post", form)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		var result map[string]interface{}
		err = json.Unmarshal(resp.Bytes(), &result)
		assert.NoError(t, err)
		assert.Contains(t, result["form"].(map[string]interface{}), "key")
	})
}

func TestPut(t *testing.T) {
	t.Run("Valid PUT Request", func(t *testing.T) {
		form := url.NewForm()
		form.Set("key", "value")
		resp, err := Put("https://httpbin.org/put", form)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		var result map[string]interface{}
		err = json.Unmarshal(resp.Bytes(), &result)
		assert.NoError(t, err)
		assert.Contains(t, result["form"].(map[string]interface{}), "key")
	})
}

func TestDelete(t *testing.T) {
	t.Run("Valid DELETE Request", func(t *testing.T) {
		params := url.NewURLParams()
		params.Set("key", "value")
		resp, err := Delete("https://httpbin.org/delete", params)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestPatch(t *testing.T) {
	t.Run("Valid PATCH Request", func(t *testing.T) {
		form := url.NewForm()
		form.Set("key", "value")
		resp, err := Patch("https://httpbin.org/patch", form)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		var result map[string]interface{}
		err = json.Unmarshal(resp.Bytes(), &result)
		assert.NoError(t, err)
		assert.Contains(t, result["form"].(map[string]interface{}), "key")
	})
}

func TestNewRequest(t *testing.T) {
	t.Run("Create New Request", func(t *testing.T) {
		params := url.NewURLParams()
		params.Set("key", "value")
		body := strings.NewReader("body content")
		req, err := NewRequest("POST", "https://example.com", params, body)
		assert.NoError(t, err)
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "https://example.com?key=value", req.URL.String())
		assert.NotNil(t, req.Body)
	})
	t.Run("Invalid URL", func(t *testing.T) {
		_, err := NewRequest("GET", "://invalid-url", nil, nil)
		assert.Error(t, err)
	})
}

func TestSetHTTP2Enabled(t *testing.T) {
	t.Run("Enable HTTP/2", func(t *testing.T) {
		SetHTTP2Enabled(true)
		assert.True(t, requests.IsHTTP2Enabled())
	})
	t.Run("Disable HTTP/2", func(t *testing.T) {
		SetHTTP2Enabled(false)
		assert.False(t, requests.IsHTTP2Enabled())
	})
}

func TestNewsession(t *testing.T) {
	t.Run("Create New Session", func(t *testing.T) {
		session := NewSession()
		assert.NotNil(t, session)
	})
}
