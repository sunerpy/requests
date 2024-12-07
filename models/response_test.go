package models

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Name      string
	Fields    any
	Args      any
	Want      any
	ExpectErr bool
}

func initResponse(fields any) *Response {
	f := fields.(struct {
		StatusCode int
		Headers    http.Header
		Cookies    []*http.Cookie
		body       []byte
		finalURL   string
		once       *sync.Once
		cachedErr  error
		Proto      string
		rawResp    *http.Response
	})
	return &Response{
		StatusCode: f.StatusCode,
		Headers:    f.Headers,
		Cookies:    f.Cookies,
		body:       f.body,
		finalURL:   f.finalURL,
		once:       f.once,
		cachedErr:  f.cachedErr,
		Proto:      f.Proto,
		rawResp:    f.rawResp,
	}
}

func TestNewResponse_Error(t *testing.T) {
	mockResp := &http.Response{
		Body: io.NopCloser(&errorReader{}),
	}
	finalURL := "https://example.com"
	resp, err := NewResponse(mockResp, finalURL)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, io.ErrUnexpectedEOF, err)
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestNewResponse(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"key": "value"}`)),
		Proto:      "HTTP/2.0",
	}
	finalURL := "https://example.com"
	resp, err := NewResponse(mockResp, finalURL)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Headers.Get("Content-Type"))
	assert.Equal(t, `{"key": "value"}`, resp.Text())
	assert.Equal(t, "HTTP/2.0", resp.Proto)
	assert.Equal(t, finalURL, resp.GetURL())
}

func TestResponseBytes(t *testing.T) {
	t.Run("Valid Bytes Output", func(t *testing.T) {
		response := &Response{
			body: []byte("Hello, world!"),
		}
		assert.Equal(t, []byte("Hello, world!"), response.Bytes())
	})
	t.Run("Empty Body", func(t *testing.T) {
		response := &Response{
			body: []byte{},
		}
		assert.Equal(t, []byte{}, response.Bytes())
	})
}

func TestResponse_Content(t *testing.T) {
	r := &Response{
		body: []byte("test content"),
	}
	content := r.Content()
	assert.NotNil(t, content)
	assert.Equal(t, []byte("test content"), content.body)
}

func TestResponse_Raw(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
	r := &Response{
		rawResp: mockResp,
	}
	raw := r.Raw()
	assert.NotNil(t, raw)
	assert.Equal(t, mockResp, raw)
	assert.Equal(t, 200, raw.StatusCode)
	assert.Equal(t, "application/json", raw.Header.Get("Content-Type"))
}

func TestResponse(t *testing.T) {
	testCases := []TestCase{
		{
			Name: "Test NewResponse",
			Fields: struct {
				StatusCode int
				Headers    http.Header
				Cookies    []*http.Cookie
				body       []byte
				finalURL   string
				once       *sync.Once
				cachedErr  error
				Proto      string
				rawResp    *http.Response
			}{
				StatusCode: 200,
				Headers:    http.Header{"Content-Type": []string{"application/json"}},
				body:       []byte(`{"key": "value"}`),
				finalURL:   "https://example.com",
			},
			Want:      "https://example.com",
			ExpectErr: false,
		},
		{
			Name: "Test Text",
			Fields: struct {
				StatusCode int
				Headers    http.Header
				Cookies    []*http.Cookie
				body       []byte
				finalURL   string
				once       *sync.Once
				cachedErr  error
				Proto      string
				rawResp    *http.Response
			}{
				body: []byte("plain text"),
			},
			Want:      "plain text",
			ExpectErr: false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			r := initResponse(tt.Fields)
			switch tt.Name {
			case "Test NewResponse":
				assert.Equal(t, tt.Want, r.GetURL())
			case "Test Text":
				assert.Equal(t, tt.Want, r.Text())
			}
		})
	}
}

func TestContentWrapper(t *testing.T) {
	testCases := []TestCase{
		{
			Name: "Decode UTF-8",
			Fields: struct {
				body []byte
			}{
				body: []byte("test content"),
			},
			Args:      "utf-8",
			Want:      "test content",
			ExpectErr: false,
		},
		{
			Name: "Decode Latin1",
			Fields: struct {
				body []byte
			}{
				body: []byte{0xE9, 0xE8, 0xE7},
			},
			Args:      "latin1",
			Want:      "éèç",
			ExpectErr: false,
		},
		{
			Name: "Unsupported Encoding",
			Fields: struct {
				body []byte
			}{
				body: []byte("unsupported content"),
			},
			Args:      "unsupported",
			Want:      nil,
			ExpectErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			c := &ContentWrapper{body: tt.Fields.(struct{ body []byte }).body}
			got, err := c.Decode(tt.Args.(string))
			if tt.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Want, got)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	type resultStruct struct {
		Key string `json:"key"`
	}
	testCases := []TestCase{
		{
			Name: "Success",
			Fields: struct {
				body []byte
			}{
				body: []byte(`{"key": "value"}`),
			},
			Want:      resultStruct{Key: "value"},
			ExpectErr: false,
		},
		{
			Name: "Error",
			Fields: struct {
				body []byte
			}{
				body: []byte(`{"key": `),
			},
			Want:      resultStruct{},
			ExpectErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			r := &Response{body: tt.Fields.(struct{ body []byte }).body}
			got, err := JSON[resultStruct](r)
			if tt.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Want, got)
			}
		})
	}
}
