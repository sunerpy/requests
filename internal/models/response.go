package models

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/sunerpy/requests/codec"
)

// Buffer pool sizes
const (
	smallBufSize = 4 * 1024         // 4KB for small responses
	largeBufSize = 10 * 1024 * 1024 // 10MB threshold
)

// Buffer pool for response reading
var smallBufPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, smallBufSize))
	},
}

// Response 表示 HTTP 响应
type Response struct {
	StatusCode int
	Status     string // HTTP status text (e.g., "200 OK")
	Headers    http.Header
	Cookies    []*http.Cookie
	body       []byte
	finalURL   string
	once       *sync.Once
	cachedErr  error
	Proto      string // 协议版本 (HTTP/1.1 或 HTTP/2.0)
	rawResp    *http.Response
}

// NewResponse 初始化 Response 对象并读取响应体
func NewResponse(resp *http.Response, finalURL string) (*Response, error) {
	defer resp.Body.Close()
	var data []byte
	var err error
	contentLen := resp.ContentLength
	if contentLen > 0 && contentLen < largeBufSize {
		// Known content length within reasonable size: allocate exact size
		data = make([]byte, contentLen)
		_, err = io.ReadFull(resp.Body, data)
	} else if contentLen == -1 || contentLen == 0 {
		// Unknown or zero content length: use buffer pool based on actual read size
		buf := smallBufPool.Get().(*bytes.Buffer)
		buf.Reset()
		_, err = buf.ReadFrom(resp.Body)
		if err == nil {
			// Copy data out of pooled buffer
			data = make([]byte, buf.Len())
			copy(data, buf.Bytes())
		}
		smallBufPool.Put(buf)
	} else {
		// Very large content (>= 10MB): direct allocation
		data, err = io.ReadAll(resp.Body)
	}
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Cookies:    resp.Cookies(),
		body:       data,
		finalURL:   finalURL,
		Proto:      resp.Proto,
		rawResp:    resp,
	}, nil
}

// NewResponseFast creates a Response without copying the body.
// Warning: The response body must not be modified after this call.
// Use this only when you know the body won't be modified.
func NewResponseFast(resp *http.Response, finalURL string, body []byte) *Response {
	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Cookies:    resp.Cookies(),
		body:       body,
		finalURL:   finalURL,
		Proto:      resp.Proto,
		rawResp:    resp,
	}
}

// Content 返回支持链式解码的封装对象
func (r *Response) Content() *ContentWrapper {
	return &ContentWrapper{body: r.body}
}

// Raw 返回原始 http.Response 对象
func (r *Response) Raw() *http.Response {
	return r.rawResp
}

// Text 返回响应的文本形式（使用 UTF-8 解码）
func (r *Response) Text() string {
	return string(r.body)
}

func (r *Response) Bytes() []byte {
	return r.body
}

// DecodeError represents an error during response body decoding.
type DecodeError struct {
	ContentType string
	Err         error
}

func (e *DecodeError) Error() string {
	return "decode error for content type " + e.ContentType + ": " + e.Err.Error()
}

func (e *DecodeError) Unwrap() error {
	return e.Err
}

// JSON 使用泛型解析 JSON
func JSON[T any](r *Response) (T, error) {
	var result T
	err := codec.Unmarshal(r.body, &result)
	if err != nil {
		return result, &DecodeError{ContentType: "application/json", Err: err}
	}
	return result, nil
}

// XML 使用泛型解析 XML
func XML[T any](r *Response) (T, error) {
	var result T
	err := xml.Unmarshal(r.body, &result)
	if err != nil {
		return result, &DecodeError{ContentType: "application/xml", Err: err}
	}
	return result, nil
}

func (r *Response) DecodeJSON(dest any) error {
	if err := codec.Unmarshal(r.body, dest); err != nil {
		return &DecodeError{ContentType: "application/json", Err: err}
	}
	return nil
}

func (r *Response) DecodeXML(dest any) error {
	if err := xml.Unmarshal(r.body, dest); err != nil {
		return &DecodeError{ContentType: "application/xml", Err: err}
	}
	return nil
}

// IsSuccess returns true if status code is 2xx.
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError returns true if status code is 4xx or 5xx.
func (r *Response) IsError() bool {
	return r.StatusCode >= 400
}

// IsRedirect returns true if status code is 3xx.
func (r *Response) IsRedirect() bool {
	return r.StatusCode >= 300 && r.StatusCode < 400
}

// IsClientError returns true if status code is 4xx.
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if status code is 5xx.
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500
}

// ContentType returns the Content-Type header value.
func (r *Response) ContentType() string {
	return r.Headers.Get("Content-Type")
}

// ContentWrapper 封装响应体的二进制数据，并提供链式解码功能
type ContentWrapper struct {
	body []byte
}

// Decode 使用指定的编码格式对内容进行解码
func (c *ContentWrapper) Decode(encoding string) (string, error) {
	switch strings.ToLower(encoding) {
	case "utf-8", "utf8":
		return string(c.body), nil
	case "latin1", "iso-8859-1":
		// 解码 Latin1 编码为 Unicode
		return decodeLatin1ToUTF8(c.body), nil
	default:
		return "", ErrUnsupportedEncoding
	}
}

// decodeLatin1ToUTF8 将 Latin1 编码的字节解码为 UTF-8
func decodeLatin1ToUTF8(data []byte) string {
	buf := bytes.Buffer{}
	for _, b := range data {
		buf.WriteRune(rune(b))
	}
	return buf.String()
}

// ErrUnsupportedEncoding 表示不支持的编码格式
var ErrUnsupportedEncoding = errors.New("unsupported encoding format")

func (r *Response) GetURL() string {
	return r.finalURL
}

// CreateMockResponse creates a Response for testing purposes.
func CreateMockResponse(statusCode int, body []byte, headers http.Header) *Response {
	if headers == nil {
		headers = make(http.Header)
	}
	return &Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Headers:    headers,
		Cookies:    nil,
		Proto:      "HTTP/1.1",
		body:       body,
		finalURL:   "",
		rawResp:    nil,
	}
}
