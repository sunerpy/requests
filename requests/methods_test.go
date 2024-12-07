package requests

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	customurl "github.com/sunerpy/requests/url"
)

func newMockSession(response *http.Response, err error) Session {
	return &defaultSession{
		client: &http.Client{Transport: &MockTransport{Response: response, Err: err}},
	}
}

func TestGet_NormalRequest(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"key": "value"}`)),
	}
	defaultSess = newMockSession(mockResp, nil)
	params := customurl.NewValues()
	params.Set("query", "golang")
	resp, err := Get("https://example.com", params)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, `{"key": "value"}`, resp.Text())
}

func TestGet_ErrorInNewRequest_WithGomonkey(t *testing.T) {
	patches := gomonkey.ApplyFunc(NewRequest, func(method, rawURL string, params *customurl.Values, body io.Reader) (*Request, error) {
		return nil, errors.New("request creation failed")
	})
	defer patches.Reset()
	params := customurl.NewValues()
	params.Set("query", "golang")
	resp, err := Get("https://example.com", params)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "request creation failed", err.Error())
}

func TestPost(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 201,
		Body:       io.NopCloser(strings.NewReader(`{"status": "created"}`)),
	}
	defaultSess = newMockSession(mockResp, nil)
	form := customurl.NewValues()
	form.Set("name", "John Doe")
	resp, err := Post("https://example.com", form)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, `{"status": "created"}`, resp.Text())
}

func TestPut(t *testing.T) {
	t.Run("Normal Request", func(t *testing.T) {
		mockResp := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"status": "updated"}`)),
		}
		defaultSess = newMockSession(mockResp, nil)
		form := customurl.NewValues()
		form.Set("id", "123")
		form.Set("name", "Jane Doe")
		resp, err := Put("https://example.com", form)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, `{"status": "updated"}`, resp.Text())
	})
	t.Run("Error in NewRequest", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(NewRequest, func(method string, rawURL string, params *customurl.Values, body io.Reader) (*Request, error) {
			return nil, errors.New("request creation failed")
		})
		defer patches.Reset()
		form := customurl.NewValues()
		form.Set("id", "123")
		form.Set("name", "Jane Doe")
		resp, err := Put("https://example.com", form)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Normal Request", func(t *testing.T) {
		mockResp := &http.Response{
			StatusCode: 204,
			Body:       io.NopCloser(bytes.NewReader(nil)),
		}
		defaultSess = newMockSession(mockResp, nil)
		params := customurl.NewValues()
		params.Set("id", "123")
		resp, err := Delete("https://example.com", params)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 204, resp.StatusCode)
		assert.Equal(t, "", resp.Text())
	})
	t.Run("Invalid URL", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(customurl.BuildURL, func(base string, params *customurl.Values) (string, error) {
			return "", errors.New("request creation failed")
		})
		defer patches.Reset()
		params := customurl.NewValues()
		params.Set("id", "123")
		resp, err := Delete("https://example.com", params)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
	t.Run("Error in NewRequest", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(NewRequest, func(method string, rawURL string, params *customurl.Values, body io.Reader) (*Request, error) {
			return nil, errors.New("request creation failed")
		})
		defer patches.Reset()
		params := customurl.NewValues()
		params.Set("id", "123")
		resp, err := Delete("https://example.com", params)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestPatch(t *testing.T) {
	t.Run("Normal Request", func(t *testing.T) {
		mockResp := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"status": "patched"}`)),
		}
		defaultSess = newMockSession(mockResp, nil)
		form := customurl.NewValues()
		form.Set("id", "123")
		form.Set("name", "Updated Name")
		resp, err := Patch("https://example.com", form)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, `{"status": "patched"}`, resp.Text())
	})
	t.Run("Error in NewRequest", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(NewRequest, func(method string, rawURL string, params *customurl.Values, body io.Reader) (*Request, error) {
			return nil, errors.New("request creation failed")
		})
		defer patches.Reset()
		form := customurl.NewValues()
		form.Set("id", "123")
		form.Set("name", "Updated Name")
		resp, err := Patch("https://example.com", form)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGet_InvalidURL(t *testing.T) {
	resp, err := Get(":", nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestPost_InvalidURL(t *testing.T) {
	resp, err := Post(":", nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
