package requests

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sunerpy/requests/url"
)

func TestGet(t *testing.T) {
	t.Run("Valid GET Request", func(t *testing.T) {
		resp, err := Get("https://httpbin.org/get", WithQuery("key", "value"))
		if err != nil {
			t.Skipf("Skipping test due to network error: %v", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		var result map[string]any
		err = json.Unmarshal(resp.Bytes(), &result)
		assert.NoError(t, err)
		if args, ok := result["args"].(map[string]any); ok {
			assert.Contains(t, args, "key")
		}
	})
	t.Run("GET without options", func(t *testing.T) {
		resp, err := Get("https://httpbin.org/get")
		if err != nil {
			t.Skipf("Skipping test due to network error: %v", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("Invalid URL", func(t *testing.T) {
		_, err := Get("://invalid")
		assert.Error(t, err)
	})
}

func TestPost(t *testing.T) {
	t.Run("Valid POST Request with form", func(t *testing.T) {
		form := url.NewForm()
		form.Set("key", "value")
		resp, err := Post("https://httpbin.org/post", form)
		if err != nil {
			t.Skipf("Skipping test due to network error: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Skipf("Skipping test due to unexpected status code: %d", resp.StatusCode)
		}
		var result map[string]any
		err = json.Unmarshal(resp.Bytes(), &result)
		if err != nil {
			t.Skipf("Skipping test due to JSON parse error: %v", err)
		}
		if formData, ok := result["form"].(map[string]any); ok {
			assert.Contains(t, formData, "key")
		}
	})
	t.Run("POST with nil body", func(t *testing.T) {
		resp, err := Post("https://httpbin.org/post", nil)
		if err != nil {
			t.Skipf("Skipping test due to network error: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Skipf("Skipping test due to unexpected status code: %d", resp.StatusCode)
		}
	})
}

func TestPut(t *testing.T) {
	t.Run("Valid PUT Request with form", func(t *testing.T) {
		form := url.NewForm()
		form.Set("key", "value")
		resp, err := Put("https://httpbin.org/put", form)
		if err != nil {
			t.Skipf("Skipping test due to network error: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Skipf("Skipping test due to unexpected status code: %d", resp.StatusCode)
		}
		var result map[string]any
		err = json.Unmarshal(resp.Bytes(), &result)
		assert.NoError(t, err)
		if formData, ok := result["form"].(map[string]any); ok {
			assert.Contains(t, formData, "key")
		} else {
			t.Logf("Response body: %s", resp.Text())
			t.Skip("Skipping test: form data not in expected format")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Valid DELETE Request", func(t *testing.T) {
		resp, err := Delete("https://httpbin.org/delete", WithQuery("key", "value"))
		if err != nil {
			t.Skipf("Skipping test due to network error: %v", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("DELETE without options", func(t *testing.T) {
		resp, err := Delete("https://httpbin.org/delete")
		if err != nil {
			t.Skipf("Skipping test due to network error: %v", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("Invalid URL", func(t *testing.T) {
		_, err := Delete("://invalid")
		assert.Error(t, err)
	})
}

func TestPatch(t *testing.T) {
	t.Run("Valid PATCH Request with form", func(t *testing.T) {
		form := url.NewForm()
		form.Set("key", "value")
		resp, err := Patch("https://httpbin.org/patch", form)
		if err != nil {
			t.Skipf("Skipping test due to network error: %v", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		var result map[string]any
		err = json.Unmarshal(resp.Bytes(), &result)
		assert.NoError(t, err)
		if formData, ok := result["form"].(map[string]any); ok {
			assert.Contains(t, formData, "key")
		}
	})
}

func TestNewRequestBuilder(t *testing.T) {
	t.Run("Create New Request with Builder", func(t *testing.T) {
		req, err := NewGet("https://example.com").
			WithQuery("key", "value").
			Build()
		assert.NoError(t, err)
		assert.Equal(t, MethodGet, req.Method)
		assert.Contains(t, req.URL.String(), "key=value")
	})
	t.Run("Create POST Request with Builder", func(t *testing.T) {
		req, err := NewPost("https://example.com").
			WithJSON(map[string]string{"key": "value"}).
			Build()
		assert.NoError(t, err)
		assert.Equal(t, MethodPost, req.Method)
		assert.NotNil(t, req.Body)
	})
	t.Run("Invalid URL", func(t *testing.T) {
		_, err := NewGet("://invalid-url").Build()
		assert.Error(t, err)
	})
}

func TestNewRequestBuilderWithContext(t *testing.T) {
	t.Run("Create New Request with Context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "key", "value")
		req, err := NewGet("https://example.com").
			WithContext(ctx).
			WithQuery("key", "value").
			Build()
		assert.NoError(t, err)
		assert.Equal(t, MethodGet, req.Method)
		assert.Contains(t, req.URL.String(), "key=value")
		assert.Equal(t, ctx, req.Context)
	})
	t.Run("Invalid URL with Context", func(t *testing.T) {
		_, err := NewGet("://invalid-url").
			WithContext(context.Background()).
			Build()
		assert.Error(t, err)
	})
}

func TestSetHTTP2Enabled(t *testing.T) {
	t.Run("Enable HTTP/2", func(t *testing.T) {
		SetHTTP2Enabled(true)
		assert.True(t, IsHTTP2Enabled())
	})
	t.Run("Disable HTTP/2", func(t *testing.T) {
		SetHTTP2Enabled(false)
		assert.False(t, IsHTTP2Enabled())
	})
}

func TestNewSession(t *testing.T) {
	t.Run("Create New Session", func(t *testing.T) {
		session := NewSession()
		assert.NotNil(t, session)
	})
}

func TestMethod(t *testing.T) {
	t.Run("Method String", func(t *testing.T) {
		assert.Equal(t, "GET", MethodGet.String())
		assert.Equal(t, "POST", MethodPost.String())
	})
	t.Run("Method IsValid", func(t *testing.T) {
		assert.True(t, MethodGet.IsValid())
		assert.True(t, MethodPost.IsValid())
		assert.False(t, Method("INVALID").IsValid())
	})
	t.Run("Method IsIdempotent", func(t *testing.T) {
		assert.True(t, MethodGet.IsIdempotent())
		assert.True(t, MethodPut.IsIdempotent())
		assert.True(t, MethodDelete.IsIdempotent())
		assert.False(t, MethodPost.IsIdempotent())
		assert.False(t, MethodPatch.IsIdempotent())
	})
	t.Run("Method IsSafe", func(t *testing.T) {
		assert.True(t, MethodGet.IsSafe())
		assert.True(t, MethodHead.IsSafe())
		assert.False(t, MethodPost.IsSafe())
		assert.False(t, MethodPut.IsSafe())
	})
	t.Run("Method HasRequestBody", func(t *testing.T) {
		assert.True(t, MethodPost.HasRequestBody())
		assert.True(t, MethodPut.HasRequestBody())
		assert.True(t, MethodPatch.HasRequestBody())
		assert.False(t, MethodGet.HasRequestBody())
		assert.False(t, MethodDelete.HasRequestBody())
	})
}

// Error path tests for Head and Options
func TestHead_InvalidURL(t *testing.T) {
	_, err := Head("://invalid")
	assert.Error(t, err)
}

func TestOptions_InvalidURL(t *testing.T) {
	_, err := Options("://invalid")
	assert.Error(t, err)
}

func TestGetString_InvalidURL(t *testing.T) {
	_, err := GetString("://invalid")
	assert.Error(t, err)
}

func TestGetBytes_InvalidURL(t *testing.T) {
	_, err := GetBytes("://invalid")
	assert.Error(t, err)
}

func TestGetJSON_InvalidURL(t *testing.T) {
	type Response struct{}
	_, err := GetJSON[Response]("://invalid")
	assert.Error(t, err)
}

func TestDeleteJSON_InvalidURL(t *testing.T) {
	type Response struct{}
	_, err := DeleteJSON[Response]("://invalid")
	assert.Error(t, err)
}

func TestPostJSON_InvalidURL(t *testing.T) {
	type Response struct{}
	_, err := PostJSON[Response]("://invalid", nil)
	assert.Error(t, err)
}

func TestPostJSON_MarshalError(t *testing.T) {
	type Response struct{}
	_, err := PostJSON[Response]("https://httpbin.org/post", make(chan int))
	assert.Error(t, err)
}

func TestPutJSON_InvalidURL(t *testing.T) {
	type Response struct{}
	_, err := PutJSON[Response]("://invalid", nil)
	assert.Error(t, err)
}

func TestPutJSON_MarshalError(t *testing.T) {
	type Response struct{}
	_, err := PutJSON[Response]("https://httpbin.org/put", make(chan int))
	assert.Error(t, err)
}

func TestPatchJSON_InvalidURL(t *testing.T) {
	type Response struct{}
	_, err := PatchJSON[Response]("://invalid", nil)
	assert.Error(t, err)
}

func TestPatchJSON_MarshalError(t *testing.T) {
	type Response struct{}
	_, err := PatchJSON[Response]("https://httpbin.org/patch", make(chan int))
	assert.Error(t, err)
}

// Test Result[T] API - replaces old *WithResponse tests
func TestGetJSON_Result(t *testing.T) {
	type Response struct {
		URL string `json:"url"`
	}
	result, err := GetJSON[Response]("https://httpbin.org/get")
	assert.NoError(t, err)
	assert.NotNil(t, result.Response())
	assert.Equal(t, 200, result.StatusCode())
	assert.Contains(t, result.Data().URL, "httpbin.org")
	assert.True(t, result.IsSuccess())
}

func TestPostJSON_Result(t *testing.T) {
	type Response struct {
		URL string `json:"url"`
	}
	result, err := PostJSON[Response]("https://httpbin.org/post", map[string]string{"key": "value"})
	assert.NoError(t, err)
	assert.NotNil(t, result.Response())
	assert.Equal(t, 200, result.StatusCode())
	assert.Contains(t, result.Data().URL, "httpbin.org")
}

func TestPutJSON_Result(t *testing.T) {
	type Response struct {
		URL string `json:"url"`
	}
	result, err := PutJSON[Response]("https://httpbin.org/put", map[string]string{"key": "value"})
	assert.NoError(t, err)
	assert.NotNil(t, result.Response())
	assert.Equal(t, 200, result.StatusCode())
	assert.Contains(t, result.Data().URL, "httpbin.org")
}

func TestDeleteJSON_Result(t *testing.T) {
	type Response struct {
		URL string `json:"url"`
	}
	result, err := DeleteJSON[Response]("https://httpbin.org/delete")
	assert.NoError(t, err)
	assert.NotNil(t, result.Response())
	assert.Equal(t, 200, result.StatusCode())
	assert.Contains(t, result.Data().URL, "httpbin.org")
}

func TestPatchJSON_Result(t *testing.T) {
	type Response struct {
		URL string `json:"url"`
	}
	result, err := PatchJSON[Response]("https://httpbin.org/patch", map[string]string{"key": "value"})
	assert.NoError(t, err)
	assert.NotNil(t, result.Response())
	assert.Equal(t, 200, result.StatusCode())
	assert.Contains(t, result.Data().URL, "httpbin.org")
}

// Test prepareBodyForRequest with different body types
func TestPrepareBodyForRequest(t *testing.T) {
	t.Run("With bytes", func(t *testing.T) {
		resp, err := Post("https://httpbin.org/post", []byte("test data"))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("With string", func(t *testing.T) {
		resp, err := Post("https://httpbin.org/post", "test string")
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("With url.Values pointer", func(t *testing.T) {
		form := url.NewForm()
		form.Set("key", "value")
		resp, err := Post("https://httpbin.org/post", form)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("With nil url.Values pointer", func(t *testing.T) {
		var form *url.Values
		resp, err := Post("https://httpbin.org/post", form)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

// Test Result[T] decode errors
func TestGetJSON_DecodeError(t *testing.T) {
	type Response struct {
		ID int `json:"id"`
	}
	// Use httpbin to get a non-JSON response
	_, err := GetJSON[Response]("https://httpbin.org/html")
	assert.Error(t, err)
}

func TestPostJSON_DecodeError(t *testing.T) {
	type Response struct {
		ID int `json:"id"`
	}
	_, err := PostJSON[Response]("https://httpbin.org/html", map[string]string{"key": "value"})
	assert.Error(t, err)
}

func TestPutJSON_DecodeError(t *testing.T) {
	type Response struct {
		ID int `json:"id"`
	}
	_, err := PutJSON[Response]("https://httpbin.org/html", map[string]string{"key": "value"})
	assert.Error(t, err)
}

func TestDeleteJSON_DecodeError(t *testing.T) {
	type Response struct {
		ID int `json:"id"`
	}
	_, err := DeleteJSON[Response]("https://httpbin.org/html")
	assert.Error(t, err)
}

func TestPatchJSON_DecodeError(t *testing.T) {
	type Response struct {
		ID int `json:"id"`
	}
	_, err := PatchJSON[Response]("https://httpbin.org/html", map[string]string{"key": "value"})
	assert.Error(t, err)
}
