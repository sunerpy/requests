package requests

import (
	"github.com/sunerpy/requests/internal/client"
	"github.com/sunerpy/requests/internal/models"
)

// Type aliases for client package types - allows users to import only the main package
// Result wraps both parsed response data and response metadata.
// Use Data() to access the parsed data and Response() for metadata.
type (
	Result[T any] = client.Result[T]
	// Response represents an HTTP response.
	// This is the unified response type used throughout the library.
	Response = models.Response
	// RequestBuilder provides a fluent interface for building HTTP requests.
	RequestBuilder = client.RequestBuilder
	// RequestOption is a function that modifies request configuration.
	RequestOption = client.RequestOption
	// HTTPClient is the interface for executing HTTP requests.
	HTTPClient = client.HTTPClient
	// Middleware is the interface for request/response middleware.
	Middleware = client.Middleware
	// Handler is a function that processes a request and returns a response.
	Handler = client.Handler
	// MiddlewareChain manages a chain of middleware.
	MiddlewareChain = client.MiddlewareChain
	// Hooks provides request/response lifecycle hooks.
	Hooks = client.Hooks
	// RetryPolicy defines retry behavior.
	RetryPolicy = client.RetryPolicy
	// RetryExecutor executes requests with retry logic.
	RetryExecutor = client.RetryExecutor
	// BasicAuth holds basic authentication credentials.
	BasicAuth = client.BasicAuth
)

// Interface type aliases - allows external packages to use these interfaces
type (
	// Client is the core HTTP client interface.
	// It provides Do, DoWithContext, and Clone methods.
	Client = client.Client
	// Session extends Client with session management capabilities.
	// It provides configuration methods like WithTimeout, WithRetry, WithMiddleware, etc.
	Session = client.Session
	// MiddlewareFunc is a function adapter for Middleware interface.
	MiddlewareFunc = client.MiddlewareFunc
)

// Error types
type (
	RequestError    = client.RequestError
	ResponseError   = client.ResponseError
	TimeoutError    = client.TimeoutError
	ConnectionError = client.ConnectionError
	DecodeError     = client.DecodeError
	EncodeError     = client.EncodeError
	RetryError      = client.RetryError
)

// ============================================================================
// RequestBuilder Constructors
// ============================================================================
// NewRequestBuilder creates a new RequestBuilder with the specified method and URL.
func NewRequestBuilder(method Method, rawURL string) *RequestBuilder {
	return client.NewRequest(method, rawURL)
}

// NewGet creates a new GET request builder.
func NewGet(rawURL string) *RequestBuilder {
	return client.NewGet(rawURL)
}

// NewPost creates a new POST request builder.
func NewPost(rawURL string) *RequestBuilder {
	return client.NewPost(rawURL)
}

// NewPut creates a new PUT request builder.
func NewPut(rawURL string) *RequestBuilder {
	return client.NewPut(rawURL)
}

// NewDelete creates a new DELETE request builder.
func NewDeleteBuilder(rawURL string) *RequestBuilder {
	return client.NewDelete(rawURL)
}

// NewPatch creates a new PATCH request builder.
func NewPatch(rawURL string) *RequestBuilder {
	return client.NewPatch(rawURL)
}

// ============================================================================
// Request Options
// ============================================================================
// WithTimeout sets the request timeout.
var (
	WithTimeout = client.WithTimeout
	// WithHeader adds a single header to the request.
	WithHeader = client.WithHeader
	// WithHeaders adds multiple headers to the request.
	WithHeaders = client.WithHeaders
	// WithQuery adds a single query parameter.
	WithQuery = client.WithQuery
	// WithQueryParams adds multiple query parameters.
	WithQueryParams = client.WithQueryParams
	// WithBasicAuth adds basic authentication.
	WithBasicAuth = client.WithBasicAuth
	// WithBearerToken adds bearer token authentication.
	WithBearerToken = client.WithBearerToken
	// WithContext sets the request context.
	WithContext = client.WithContext
	// WithContentType sets the Content-Type header.
	WithContentType = client.WithContentType
	// WithAccept sets the Accept header.
	WithAccept = client.WithAccept
	// ============================================================================
	// Middleware and Hooks
	// ============================================================================
	// NewMiddlewareChain creates a new middleware chain.
	NewMiddlewareChain = client.NewMiddlewareChain
	// NewHooks creates a new Hooks instance.
	NewHooks = client.NewHooks
	// NewMetricsHook creates a new metrics hook.
	NewMetricsHook = client.NewMetricsHook
	// HeaderMiddleware creates middleware that adds headers.
	HeaderMiddleware = client.HeaderMiddleware
	// HooksMiddleware creates middleware from hooks.
	HooksMiddleware = client.HooksMiddleware
	// ============================================================================
	// Retry
	// ============================================================================
	// NewRetryExecutor creates a new retry executor.
	NewRetryExecutor = client.NewRetryExecutor
	// NoRetryPolicy returns a policy that never retries.
	NoRetryPolicy = client.NoRetryPolicy
	// LinearRetryPolicy returns a policy with linear backoff.
	LinearRetryPolicy = client.LinearRetryPolicy
	// ExponentialRetryPolicy returns a policy with exponential backoff.
	ExponentialRetryPolicy = client.ExponentialRetryPolicy
	// RetryOn5xx returns a condition that retries on 5xx errors.
	RetryOn5xx = client.RetryOn5xx
	// RetryOnNetworkError returns a condition that retries on network errors.
	RetryOnNetworkError = client.RetryOnNetworkError
	// RetryOnStatusCodes returns a condition that retries on specific status codes.
	RetryOnStatusCodes = client.RetryOnStatusCodes
	// CombineRetryConditions combines multiple retry conditions.
	CombineRetryConditions = client.CombineRetryConditions
)

// ============================================================================
// Generic Request Execution
// ============================================================================
// DoJSON executes a request builder and parses JSON response into Result[T].
func DoJSON[T any](b *RequestBuilder) (Result[T], error) {
	return client.DoJSON[T](b)
}

// DoXML executes a request builder and parses XML response into Result[T].
func DoXML[T any](b *RequestBuilder) (Result[T], error) {
	return client.DoXML[T](b)
}

// ============================================================================
// Error Helpers
// ============================================================================
// IsTimeout checks if an error is a timeout error.
var (
	IsTimeout = client.IsTimeout
	// IsConnectionError checks if an error is a connection error.
	IsConnectionError = client.IsConnectionError
	// IsResponseError checks if an error is a response error.
	IsResponseError = client.IsResponseError
	// IsTemporary checks if an error is temporary.
	IsTemporary = client.IsTemporary
)
