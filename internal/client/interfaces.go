// Package client provides the core HTTP client interfaces and implementations.
//
// INTERNAL PACKAGE: This package is intended for internal use by the requests library.
// Users should import the main "github.com/sunerpy/requests" package instead, which
// re-exports all necessary types and functions.
//
// The main package provides:
//   - Type aliases: Request, Response, Result[T], RequestBuilder, etc.
//   - HTTP methods: Get, Post, Put, Delete, Patch, Head, Options
//   - Generic methods: GetJSON[T], PostJSON[T], etc.
//   - Builder constructors: NewGet, NewPost, NewPut, NewDelete, NewPatch
//   - Options: WithTimeout, WithHeader, WithQuery, etc.
//
// Example usage (recommended):
//
//	import "github.com/sunerpy/requests"
//
//	// Simple GET request
//	resp, err := requests.Get("https://api.example.com/users")
//
//	// GET with JSON parsing
//	result, err := requests.GetJSON[User]("https://api.example.com/users/1")
//	user := result.Data()
//
//	// Using RequestBuilder
//	req, err := requests.NewPost("https://api.example.com/users").
//	    WithJSON(userData).
//	    WithHeader("Authorization", "Bearer token").
//	    Build()
package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sunerpy/requests/internal/models"
)

// Client is the core HTTP client interface.
type Client interface {
	// Do executes an HTTP request and returns the response.
	Do(req *Request) (*models.Response, error)
	// DoWithContext executes an HTTP request with context.
	DoWithContext(ctx context.Context, req *Request) (*models.Response, error)
	// Clone creates a copy of the client.
	Clone() Client
}

// Session extends Client with session management capabilities.
type Session interface {
	Client
	// Configuration methods - all return Session for chaining
	WithBaseURL(base string) Session
	WithTimeout(d time.Duration) Session
	WithProxy(proxyURL string) Session
	WithHeader(key, value string) Session
	WithHeaders(headers map[string]string) Session
	WithBasicAuth(username, password string) Session
	WithBearerToken(token string) Session
	WithHTTP2(enabled bool) Session
	WithKeepAlive(enabled bool) Session
	WithMaxIdleConns(n int) Session
	WithIdleTimeout(d time.Duration) Session
	WithRetry(policy RetryPolicy) Session
	WithMiddleware(m Middleware) Session
	WithCookieJar(jar http.CookieJar) Session
	// Resource management
	Close() error
	// Clear resets the session to default state
	Clear() Session
}

// Method represents an HTTP method.
type Method string

// HTTP methods as constants.
const (
	MethodGet     Method = "GET"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodDelete  Method = "DELETE"
	MethodPatch   Method = "PATCH"
	MethodHead    Method = "HEAD"
	MethodOptions Method = "OPTIONS"
	MethodConnect Method = "CONNECT"
	MethodTrace   Method = "TRACE"
)

// String returns the string representation of the method.
func (m Method) String() string {
	return string(m)
}

// IsValid checks if the method is a valid HTTP method.
func (m Method) IsValid() bool {
	switch m {
	case MethodGet, MethodPost, MethodPut, MethodDelete, MethodPatch,
		MethodHead, MethodOptions, MethodConnect, MethodTrace:
		return true
	default:
		return false
	}
}

// IsIdempotent returns true if the method is idempotent.
func (m Method) IsIdempotent() bool {
	switch m {
	case MethodGet, MethodPut, MethodDelete, MethodHead, MethodOptions, MethodTrace:
		return true
	default:
		return false
	}
}

// IsSafe returns true if the method is safe (read-only).
func (m Method) IsSafe() bool {
	switch m {
	case MethodGet, MethodHead, MethodOptions, MethodTrace:
		return true
	default:
		return false
	}
}

// HasRequestBody returns true if the method typically has a request body.
func (m Method) HasRequestBody() bool {
	switch m {
	case MethodPost, MethodPut, MethodPatch:
		return true
	default:
		return false
	}
}

// Request represents an HTTP request.
type Request struct {
	Method  Method
	URL     *url.URL
	Headers http.Header
	Body    io.Reader
	Context context.Context
	// Internal fields for retry support
	bodyBytes []byte
	files     []FileUpload
	form      url.Values
}

// Clone creates a deep copy of the Request.
// Modifications to the original Request will not affect the clone, and vice versa.
func (r *Request) Clone() *Request {
	if r == nil {
		return nil
	}
	clone := &Request{
		Method:  r.Method,
		Context: r.Context,
	}
	// Deep copy URL
	if r.URL != nil {
		clonedURL := *r.URL
		if r.URL.User != nil {
			clonedURL.User = url.UserPassword(r.URL.User.Username(), "")
			if pwd, ok := r.URL.User.Password(); ok {
				clonedURL.User = url.UserPassword(r.URL.User.Username(), pwd)
			}
		}
		clone.URL = &clonedURL
	}
	// Deep copy Headers
	if r.Headers != nil {
		clone.Headers = r.Headers.Clone()
	}
	// Deep copy bodyBytes
	if r.bodyBytes != nil {
		clone.bodyBytes = make([]byte, len(r.bodyBytes))
		copy(clone.bodyBytes, r.bodyBytes)
	}
	// Deep copy form values
	if r.form != nil {
		clone.form = make(url.Values)
		for k, v := range r.form {
			clone.form[k] = append([]string(nil), v...)
		}
	}
	// Deep copy files (note: Reader cannot be cloned, only metadata)
	if r.files != nil {
		clone.files = make([]FileUpload, len(r.files))
		copy(clone.files, r.files)
	}
	// Body (io.Reader) cannot be deep copied - it's shared
	// For retry scenarios, use bodyBytes instead
	clone.Body = r.Body
	return clone
}

// AddHeader adds a header value to the request.
// Multiple values can be added for the same key.
func (r *Request) AddHeader(key, value string) {
	if r.Headers == nil {
		r.Headers = http.Header{}
	}
	r.Headers.Add(key, value)
}

// SetHeader sets a header value on the request.
// This replaces any existing values for the key.
func (r *Request) SetHeader(key, value string) {
	if r.Headers == nil {
		r.Headers = http.Header{}
	}
	r.Headers.Set(key, value)
}

// FileUpload represents a file to be uploaded.
type FileUpload struct {
	FieldName string
	FileName  string
	Reader    io.Reader
}

// Middleware defines the request/response middleware interface.
type Middleware interface {
	// Process handles the request and calls next to continue the chain.
	Process(req *Request, next Handler) (*models.Response, error)
}

// Handler is the next handler in the middleware chain.
type (
	Handler func(*Request) (*models.Response, error)
	// MiddlewareFunc is a function adapter for Middleware interface.
	MiddlewareFunc func(req *Request, next Handler) (*models.Response, error)
)

// Process implements the Middleware interface.
func (f MiddlewareFunc) Process(req *Request, next Handler) (*models.Response, error) {
	return f(req, next)
}

// RetryPolicy defines the retry strategy.
type RetryPolicy struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	Jitter          float64
	RetryIf         func(resp *models.Response, err error) bool
}

// DefaultRetryPolicy returns a sensible default retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.1,
		RetryIf:         DefaultRetryCondition,
	}
}

// DefaultRetryCondition is the default condition for retrying requests.
func DefaultRetryCondition(resp *models.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp != nil {
		// Retry on 5xx server errors and 429 Too Many Requests
		return resp.StatusCode >= 500 || resp.StatusCode == 429
	}
	return false
}

// BasicAuth holds basic authentication credentials.
type BasicAuth struct {
	Username string
	Password string
}
