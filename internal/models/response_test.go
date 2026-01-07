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
		Status     string
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
		Status:     f.Status,
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
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"key": "value"}`)),
		Proto:      "HTTP/2.0",
	}
	finalURL := "https://example.com"
	resp, err := NewResponse(mockResp, finalURL)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "200 OK", resp.Status)
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
				Status     string
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
				Status:     "200 OK",
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
				Status     string
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

// Additional tests for coverage
func TestXML(t *testing.T) {
	type Person struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
	}
	t.Run("Success", func(t *testing.T) {
		r := &Response{body: []byte(`<Person><name>John</name><age>30</age></Person>`)}
		got, err := XML[Person](r)
		assert.NoError(t, err)
		assert.Equal(t, "John", got.Name)
		assert.Equal(t, 30, got.Age)
	})
	t.Run("Error", func(t *testing.T) {
		r := &Response{body: []byte(`<Person><name>John`)}
		_, err := XML[Person](r)
		assert.Error(t, err)
	})
}

func TestResponse_DecodeJSON(t *testing.T) {
	type Data struct {
		Key string `json:"key"`
	}
	t.Run("Success", func(t *testing.T) {
		r := &Response{body: []byte(`{"key": "value"}`)}
		var data Data
		err := r.DecodeJSON(&data)
		assert.NoError(t, err)
		assert.Equal(t, "value", data.Key)
	})
	t.Run("Error", func(t *testing.T) {
		r := &Response{body: []byte(`invalid json`)}
		var data Data
		err := r.DecodeJSON(&data)
		assert.Error(t, err)
	})
}

func TestResponse_DecodeXML(t *testing.T) {
	type Person struct {
		Name string `xml:"name"`
	}
	t.Run("Success", func(t *testing.T) {
		r := &Response{body: []byte(`<Person><name>John</name></Person>`)}
		var person Person
		err := r.DecodeXML(&person)
		assert.NoError(t, err)
		assert.Equal(t, "John", person.Name)
	})
	t.Run("Error", func(t *testing.T) {
		r := &Response{body: []byte(`<Person><name>John`)}
		var person Person
		err := r.DecodeXML(&person)
		assert.Error(t, err)
	})
}

func TestResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, true},
		{201, true},
		{204, true},
		{299, true},
		{300, false},
		{400, false},
		{500, false},
	}
	for _, tc := range tests {
		r := &Response{StatusCode: tc.statusCode}
		assert.Equal(t, tc.expected, r.IsSuccess(), "StatusCode: %d", tc.statusCode)
	}
}

func TestResponse_IsError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{300, false},
		{399, false},
		{400, true},
		{404, true},
		{500, true},
		{503, true},
	}
	for _, tc := range tests {
		r := &Response{StatusCode: tc.statusCode}
		assert.Equal(t, tc.expected, r.IsError(), "StatusCode: %d", tc.statusCode)
	}
}

func TestContentWrapper_DecodeUTF8(t *testing.T) {
	c := &ContentWrapper{body: []byte("utf8 content")}
	result, err := c.Decode("utf8")
	assert.NoError(t, err)
	assert.Equal(t, "utf8 content", result)
}

func TestContentWrapper_DecodeISO88591(t *testing.T) {
	c := &ContentWrapper{body: []byte{0xE9, 0xE8, 0xE7}}
	result, err := c.Decode("iso-8859-1")
	assert.NoError(t, err)
	assert.Equal(t, "éèç", result)
}

// Tests for buffer pool handling with different response sizes
func TestNewResponse_SmallContent(t *testing.T) {
	// Small content (< 4KB) should use small buffer pool
	smallBody := strings.Repeat("a", 1024) // 1KB
	mockResp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          io.NopCloser(strings.NewReader(smallBody)),
		ContentLength: int64(len(smallBody)),
		Proto:         "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, smallBody, resp.Text())
	assert.Equal(t, len(smallBody), len(resp.Bytes()))
}

func TestNewResponse_MediumContent(t *testing.T) {
	// Medium content with known length uses exact allocation
	mediumBody := strings.Repeat("b", 16*1024) // 16KB
	mockResp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          io.NopCloser(strings.NewReader(mediumBody)),
		ContentLength: int64(len(mediumBody)),
		Proto:         "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(mediumBody), len(resp.Bytes()))
}

// streamReader wraps a reader to simulate chunked/streaming responses
type streamReader struct {
	reader io.Reader
}

func (s *streamReader) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

func TestNewResponse_BufferPool_UnknownLength(t *testing.T) {
	// Unknown content length (-1) uses buffer pool
	body := strings.Repeat("s", 2*1024) // 2KB
	mockResp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          io.NopCloser(&streamReader{reader: strings.NewReader(body)}),
		ContentLength: -1, // Unknown
		Proto:         "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, body, resp.Text())
}

func TestNewResponse_BufferPool_ZeroLength(t *testing.T) {
	// Zero content length uses buffer pool
	body := "zero length content"
	mockResp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: 0, // Zero length triggers buffer pool
		Proto:         "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, body, resp.Text())
}

func TestNewResponse_BufferPool_ReadError(t *testing.T) {
	// Test error handling in buffer pool path
	mockResp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{},
		Body:          io.NopCloser(&errorReader{}),
		ContentLength: -1, // Unknown length triggers buffer pool
		Proto:         "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestNewResponse_LargeContent(t *testing.T) {
	// Large content with known length uses exact allocation
	largeBody := strings.Repeat("c", 128*1024) // 128KB
	mockResp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          io.NopCloser(strings.NewReader(largeBody)),
		ContentLength: int64(len(largeBody)),
		Proto:         "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(largeBody), len(resp.Bytes()))
}

func TestNewResponse_VeryLargeContent(t *testing.T) {
	// Very large content (> 256KB but < 10MB) should use direct allocation
	veryLargeBody := strings.Repeat("d", 512*1024) // 512KB
	mockResp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          io.NopCloser(strings.NewReader(veryLargeBody)),
		ContentLength: int64(len(veryLargeBody)),
		Proto:         "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(veryLargeBody), len(resp.Bytes()))
}

func TestNewResponse_UnknownContentLength(t *testing.T) {
	// Unknown content length (-1) should use small buffer pool
	body := "unknown length content"
	mockResp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: -1, // Unknown content length
		Proto:         "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, body, resp.Text())
}

func TestNewResponse_KnownExactContentLength(t *testing.T) {
	// Known content length should allocate exact size
	body := "exact size content"
	mockResp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Proto:         "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, body, resp.Text())
	assert.Equal(t, len(body), len(resp.Bytes()))
}

func TestNewResponseFast(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Proto:      "HTTP/1.1",
	}
	body := []byte(`{"key": "value"}`)
	finalURL := "https://example.com"

	resp := NewResponseFast(mockResp, finalURL, body)

	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Headers.Get("Content-Type"))
	assert.Equal(t, body, resp.Bytes())
	assert.Equal(t, finalURL, resp.GetURL())
	assert.Equal(t, "HTTP/1.1", resp.Proto)
}

func TestNewResponseFast_EmptyBody(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 204,
		Header:     http.Header{},
		Proto:      "HTTP/1.1",
	}
	resp := NewResponseFast(mockResp, "https://example.com", nil)

	assert.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode)
	assert.Nil(t, resp.Bytes())
}

// Tests for new status check methods
func TestResponse_IsRedirect(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{299, false},
		{300, true},
		{301, true},
		{302, true},
		{307, true},
		{399, true},
		{400, false},
		{500, false},
	}
	for _, tc := range tests {
		r := &Response{StatusCode: tc.statusCode}
		assert.Equal(t, tc.expected, r.IsRedirect(), "StatusCode: %d", tc.statusCode)
	}
}

func TestResponse_IsClientError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{300, false},
		{399, false},
		{400, true},
		{401, true},
		{403, true},
		{404, true},
		{499, true},
		{500, false},
		{503, false},
	}
	for _, tc := range tests {
		r := &Response{StatusCode: tc.statusCode}
		assert.Equal(t, tc.expected, r.IsClientError(), "StatusCode: %d", tc.statusCode)
	}
}

func TestResponse_IsServerError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{300, false},
		{400, false},
		{499, false},
		{500, true},
		{502, true},
		{503, true},
		{504, true},
		{599, true},
	}
	for _, tc := range tests {
		r := &Response{StatusCode: tc.statusCode}
		assert.Equal(t, tc.expected, r.IsServerError(), "StatusCode: %d", tc.statusCode)
	}
}

func TestResponse_ContentType(t *testing.T) {
	tests := []struct {
		name     string
		headers  http.Header
		expected string
	}{
		{
			name:     "JSON content type",
			headers:  http.Header{"Content-Type": []string{"application/json"}},
			expected: "application/json",
		},
		{
			name:     "HTML content type",
			headers:  http.Header{"Content-Type": []string{"text/html; charset=utf-8"}},
			expected: "text/html; charset=utf-8",
		},
		{
			name:     "No content type",
			headers:  http.Header{},
			expected: "",
		},
		{
			name:     "Nil headers",
			headers:  nil,
			expected: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := &Response{Headers: tc.headers}
			assert.Equal(t, tc.expected, r.ContentType())
		})
	}
}

func TestResponse_Status(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 404,
		Status:     "404 Not Found",
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("")),
		Proto:      "HTTP/1.1",
	}
	resp, err := NewResponse(mockResp, "https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "404 Not Found", resp.Status)
}

// Tests for DecodeError
func TestDecodeError_Error(t *testing.T) {
	err := &DecodeError{
		ContentType: "application/json",
		Err:         io.EOF,
	}
	expected := "decode error for content type application/json: EOF"
	assert.Equal(t, expected, err.Error())
}

func TestDecodeError_Unwrap(t *testing.T) {
	innerErr := io.EOF
	err := &DecodeError{
		ContentType: "application/xml",
		Err:         innerErr,
	}
	assert.Equal(t, innerErr, err.Unwrap())
}

// Tests for CreateMockResponse
func TestCreateMockResponse(t *testing.T) {
	t.Run("with headers", func(t *testing.T) {
		headers := http.Header{"Content-Type": []string{"application/json"}}
		resp := CreateMockResponse(200, []byte(`{"key":"value"}`), headers)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "OK", resp.Status)
		assert.Equal(t, "application/json", resp.Headers.Get("Content-Type"))
		assert.Equal(t, `{"key":"value"}`, resp.Text())
		assert.Equal(t, "HTTP/1.1", resp.Proto)
		assert.Empty(t, resp.GetURL())
		assert.Nil(t, resp.Raw())
	})

	t.Run("with nil headers", func(t *testing.T) {
		resp := CreateMockResponse(404, []byte("not found"), nil)

		assert.Equal(t, 404, resp.StatusCode)
		assert.Equal(t, "Not Found", resp.Status)
		assert.NotNil(t, resp.Headers)
		assert.Equal(t, "not found", resp.Text())
	})

	t.Run("with empty body", func(t *testing.T) {
		resp := CreateMockResponse(204, nil, nil)

		assert.Equal(t, 204, resp.StatusCode)
		assert.Equal(t, "No Content", resp.Status)
		assert.Empty(t, resp.Text())
		assert.Nil(t, resp.Bytes())
	})
}
