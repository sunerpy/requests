package models

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
)

// Response 表示 HTTP 响应
type Response struct {
	StatusCode int
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
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Cookies:    resp.Cookies(),
		body:       data,
		finalURL:   finalURL,
		Proto:      resp.Proto,
		rawResp:    resp, // 保存原始响应
	}, nil
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

// JSON 使用泛型解析 JSON
func JSON[T any](r *Response) (T, error) {
	var result T
	err := sonic.Unmarshal(r.body, &result)
	return result, err
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
