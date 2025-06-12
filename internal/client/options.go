package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// RequestOption is a function that configures a RequestConfig.
type (
	RequestOption func(*RequestConfig)
	// RequestConfig holds the configuration for a request.
	RequestConfig struct {
		Timeout          time.Duration
		Headers          http.Header
		Query            url.Values
		BasicAuth        *BasicAuth
		BearerToken      string
		Retry            *RetryPolicy
		Context          context.Context
		Files            []FileUploadConfig
		ProgressCallback ProgressCallback
	}
)

// configPool is a pool of RequestConfig objects for reuse
var configPool = sync.Pool{
	New: func() any {
		return &RequestConfig{
			Headers: make(http.Header, 4),
			Query:   make(url.Values, 4),
		}
	},
}

// NewRequestConfig creates a new RequestConfig with default values.
func NewRequestConfig() *RequestConfig {
	return &RequestConfig{
		Headers: make(http.Header, 4),
		Query:   make(url.Values, 4),
	}
}

// AcquireConfig gets a RequestConfig from the pool.
func AcquireConfig() *RequestConfig {
	return configPool.Get().(*RequestConfig)
}

// ReleaseConfig returns a RequestConfig to the pool.
func ReleaseConfig(c *RequestConfig) {
	if c == nil {
		return
	}
	c.Reset()
	configPool.Put(c)
}

// Reset clears the config for reuse, keeping the underlying map allocations.
func (c *RequestConfig) Reset() {
	c.Timeout = 0
	// Clear headers but keep the map
	for k := range c.Headers {
		delete(c.Headers, k)
	}
	// Clear query but keep the map
	for k := range c.Query {
		delete(c.Query, k)
	}
	c.BasicAuth = nil
	c.BearerToken = ""
	c.Retry = nil
	c.Context = nil
	// Clear files slice but keep capacity
	if c.Files != nil {
		c.Files = c.Files[:0]
	}
	c.ProgressCallback = nil
}

// Apply applies all options to the config.
func (c *RequestConfig) Apply(opts ...RequestOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
}

// Clone creates a deep copy of the config.
func (c *RequestConfig) Clone() *RequestConfig {
	clone := &RequestConfig{
		Timeout:          c.Timeout,
		Headers:          c.Headers.Clone(),
		Query:            make(url.Values),
		BearerToken:      c.BearerToken,
		Context:          c.Context,
		ProgressCallback: c.ProgressCallback,
	}
	for k, v := range c.Query {
		clone.Query[k] = append([]string(nil), v...)
	}
	if c.BasicAuth != nil {
		clone.BasicAuth = &BasicAuth{
			Username: c.BasicAuth.Username,
			Password: c.BasicAuth.Password,
		}
	}
	if c.Retry != nil {
		retry := *c.Retry
		clone.Retry = &retry
	}
	if c.Files != nil {
		clone.Files = make([]FileUploadConfig, len(c.Files))
		copy(clone.Files, c.Files)
	}
	return clone
}

// Merge merges another config into this one.
// Values from other take precedence if set.
func (c *RequestConfig) Merge(other *RequestConfig) {
	if other == nil {
		return
	}
	if other.Timeout > 0 {
		c.Timeout = other.Timeout
	}
	for k, v := range other.Headers {
		c.Headers[k] = v
	}
	for k, v := range other.Query {
		c.Query[k] = v
	}
	if other.BasicAuth != nil {
		c.BasicAuth = other.BasicAuth
	}
	if other.BearerToken != "" {
		c.BearerToken = other.BearerToken
	}
	if other.Retry != nil {
		c.Retry = other.Retry
	}
	if other.Context != nil {
		c.Context = other.Context
	}
	if other.Files != nil {
		c.Files = append(c.Files, other.Files...)
	}
	if other.ProgressCallback != nil {
		c.ProgressCallback = other.ProgressCallback
	}
}

// ApplyToRequest applies the config to an HTTP request.
func (c *RequestConfig) ApplyToRequest(req *http.Request) {
	// Apply headers
	for k, v := range c.Headers {
		for _, val := range v {
			req.Header.Add(k, val)
		}
	}
	// Apply query parameters
	if len(c.Query) > 0 {
		q := req.URL.Query()
		for k, v := range c.Query {
			for _, val := range v {
				q.Add(k, val)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
	// Apply basic auth
	if c.BasicAuth != nil {
		auth := fmt.Sprintf("%s:%s", c.BasicAuth.Username, c.BasicAuth.Password)
		encoded := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))
	}
	// Apply bearer token
	if c.BearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.BearerToken))
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(d time.Duration) RequestOption {
	return func(c *RequestConfig) {
		c.Timeout = d
	}
}

// WithHeader sets a single header.
func WithHeader(key, value string) RequestOption {
	return func(c *RequestConfig) {
		if c.Headers == nil {
			c.Headers = make(http.Header)
		}
		c.Headers.Set(key, value)
	}
}

// WithHeaders sets multiple headers.
func WithHeaders(headers map[string]string) RequestOption {
	return func(c *RequestConfig) {
		if c.Headers == nil {
			c.Headers = make(http.Header)
		}
		for k, v := range headers {
			c.Headers.Set(k, v)
		}
	}
}

// WithQuery sets a single query parameter.
func WithQuery(key, value string) RequestOption {
	return func(c *RequestConfig) {
		if c.Query == nil {
			c.Query = make(url.Values)
		}
		c.Query.Set(key, value)
	}
}

// WithQueryParams sets multiple query parameters.
func WithQueryParams(params map[string]string) RequestOption {
	return func(c *RequestConfig) {
		if c.Query == nil {
			c.Query = make(url.Values)
		}
		for k, v := range params {
			c.Query.Set(k, v)
		}
	}
}

// WithBasicAuth sets basic authentication.
func WithBasicAuth(username, password string) RequestOption {
	return func(c *RequestConfig) {
		c.BasicAuth = &BasicAuth{
			Username: username,
			Password: password,
		}
	}
}

// WithBearerToken sets bearer token authentication.
func WithBearerToken(token string) RequestOption {
	return func(c *RequestConfig) {
		c.BearerToken = token
	}
}

// WithRetry sets the retry policy.
func WithRetry(policy RetryPolicy) RequestOption {
	return func(c *RequestConfig) {
		c.Retry = &policy
	}
}

// WithContext sets the request context.
func WithContext(ctx context.Context) RequestOption {
	return func(c *RequestConfig) {
		c.Context = ctx
	}
}

// WithContentType sets the Content-Type header.
func WithContentType(contentType string) RequestOption {
	return func(c *RequestConfig) {
		if c.Headers == nil {
			c.Headers = make(http.Header)
		}
		c.Headers.Set("Content-Type", contentType)
	}
}

// WithAccept sets the Accept header.
func WithAccept(accept string) RequestOption {
	return func(c *RequestConfig) {
		if c.Headers == nil {
			c.Headers = make(http.Header)
		}
		c.Headers.Set("Accept", accept)
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(userAgent string) RequestOption {
	return func(c *RequestConfig) {
		if c.Headers == nil {
			c.Headers = make(http.Header)
		}
		c.Headers.Set("User-Agent", userAgent)
	}
}

// ProgressCallback is called during file upload to report progress.
type (
	ProgressCallback func(uploaded, total int64)
	// FileUploadConfig holds file upload configuration.
	FileUploadConfig struct {
		FieldName string
		FileName  string
		FilePath  string
		Reader    any // io.Reader
		Size      int64
		Progress  ProgressCallback
	}
)

// WithFile adds a file for upload by path.
func WithFile(fieldName, filePath string) RequestOption {
	return func(c *RequestConfig) {
		if c.Files == nil {
			c.Files = make([]FileUploadConfig, 0)
		}
		c.Files = append(c.Files, FileUploadConfig{
			FieldName: fieldName,
			FilePath:  filePath,
		})
	}
}

// WithFileReader adds a file for upload from an io.Reader.
func WithFileReader(fieldName, fileName string, reader any, size int64) RequestOption {
	return func(c *RequestConfig) {
		if c.Files == nil {
			c.Files = make([]FileUploadConfig, 0)
		}
		c.Files = append(c.Files, FileUploadConfig{
			FieldName: fieldName,
			FileName:  fileName,
			Reader:    reader,
			Size:      size,
		})
	}
}

// WithProgress sets a progress callback for file uploads.
func WithProgress(callback ProgressCallback) RequestOption {
	return func(c *RequestConfig) {
		c.ProgressCallback = callback
	}
}
