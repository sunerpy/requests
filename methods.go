package requests

import (
	"io"
	"strings"

	"github.com/sunerpy/requests/codec"
	"github.com/sunerpy/requests/internal/client"
	"github.com/sunerpy/requests/url"
)

const (
	contentKey      = "Content-Type"
	jsonContentType = "application/json"
)

// ============================================================================
// Basic HTTP Methods (with variadic options)
// ============================================================================
// Get sends a GET request and returns the response.
func Get(baseURL string, opts ...client.RequestOption) (*Response, error) {
	config := client.AcquireConfig()
	defer client.ReleaseConfig(config)
	config.Apply(opts...)
	// Build URL with query parameters from options
	u, err := url.BuildURLWithQuery(baseURL, config.Query)
	if err != nil {
		return nil, err
	}
	req, err := NewGet(u).Build()
	if err != nil {
		return nil, err
	}
	applyConfigToRequest(req, config)
	return defaultSess.Do(req)
}

// Post sends a POST request with optional body and returns the response.
func Post(baseURL string, body any, opts ...client.RequestOption) (*Response, error) {
	config := client.AcquireConfig()
	defer client.ReleaseConfig(config)
	config.Apply(opts...)
	reader, contentType := prepareBodyForRequest(body)
	builder := NewPost(baseURL)
	if reader != nil {
		builder = builder.WithBody(reader)
	}
	req, err := builder.Build()
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.AddHeader(contentKey, contentType)
	}
	applyConfigToRequest(req, config)
	return defaultSess.Do(req)
}

// Put sends a PUT request with optional body and returns the response.
func Put(baseURL string, body any, opts ...client.RequestOption) (*Response, error) {
	config := client.AcquireConfig()
	defer client.ReleaseConfig(config)
	config.Apply(opts...)
	reader, contentType := prepareBodyForRequest(body)
	builder := NewPut(baseURL)
	if reader != nil {
		builder = builder.WithBody(reader)
	}
	req, err := builder.Build()
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.AddHeader(contentKey, contentType)
	}
	applyConfigToRequest(req, config)
	return defaultSess.Do(req)
}

// Delete sends a DELETE request and returns the response.
func Delete(baseURL string, opts ...client.RequestOption) (*Response, error) {
	config := client.AcquireConfig()
	defer client.ReleaseConfig(config)
	config.Apply(opts...)
	u, err := url.BuildURLWithQuery(baseURL, config.Query)
	if err != nil {
		return nil, err
	}
	req, err := NewDeleteBuilder(u).Build()
	if err != nil {
		return nil, err
	}
	applyConfigToRequest(req, config)
	return defaultSess.Do(req)
}

// Patch sends a PATCH request with optional body and returns the response.
func Patch(baseURL string, body any, opts ...client.RequestOption) (*Response, error) {
	config := client.AcquireConfig()
	defer client.ReleaseConfig(config)
	config.Apply(opts...)
	reader, contentType := prepareBodyForRequest(body)
	builder := NewPatch(baseURL)
	if reader != nil {
		builder = builder.WithBody(reader)
	}
	req, err := builder.Build()
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.AddHeader(contentKey, contentType)
	}
	applyConfigToRequest(req, config)
	return defaultSess.Do(req)
}

// Head sends a HEAD request and returns the response.
func Head(baseURL string, opts ...client.RequestOption) (*Response, error) {
	config := client.AcquireConfig()
	defer client.ReleaseConfig(config)
	config.Apply(opts...)
	u, err := url.BuildURLWithQuery(baseURL, config.Query)
	if err != nil {
		return nil, err
	}
	req, err := NewRequestBuilder(MethodHead, u).Build()
	if err != nil {
		return nil, err
	}
	applyConfigToRequest(req, config)
	return defaultSess.Do(req)
}

// Options sends an OPTIONS request and returns the response.
func Options(baseURL string, opts ...client.RequestOption) (*Response, error) {
	config := client.AcquireConfig()
	defer client.ReleaseConfig(config)
	config.Apply(opts...)
	u, err := url.BuildURLWithQuery(baseURL, config.Query)
	if err != nil {
		return nil, err
	}
	req, err := NewRequestBuilder(MethodOptions, u).Build()
	if err != nil {
		return nil, err
	}
	applyConfigToRequest(req, config)
	return defaultSess.Do(req)
}

// DefaultSession returns the default session.
func DefaultSession() Session {
	return defaultSess
}

// ============================================================================
// Generic HTTP Methods - return (Result[T], error)
// ============================================================================
// GetJSON sends a GET request and returns Result[T] with parsed JSON response.
func GetJSON[T any](baseURL string, opts ...client.RequestOption) (Result[T], error) {
	var zero Result[T]
	resp, err := Get(baseURL, opts...)
	if err != nil {
		return zero, err
	}
	var result T
	if err := codec.Unmarshal(resp.Bytes(), &result); err != nil {
		return zero, &client.DecodeError{ContentType: jsonContentType, Err: err}
	}
	return client.NewResult(result, resp), nil
}

// PostJSON sends a POST request with JSON body and returns Result[T].
func PostJSON[T any](baseURL string, data any, opts ...client.RequestOption) (Result[T], error) {
	var zero Result[T]
	jsonData, err := codec.Marshal(data)
	if err != nil {
		return zero, &client.EncodeError{ContentType: jsonContentType, Err: err}
	}
	resp, err := Post(baseURL, strings.NewReader(string(jsonData)), append(opts, client.WithContentType(jsonContentType))...)
	if err != nil {
		return zero, err
	}
	var result T
	if err := codec.Unmarshal(resp.Bytes(), &result); err != nil {
		return zero, &client.DecodeError{ContentType: jsonContentType, Err: err}
	}
	return client.NewResult(result, resp), nil
}

// PutJSON sends a PUT request with JSON body and returns Result[T].
func PutJSON[T any](baseURL string, data any, opts ...client.RequestOption) (Result[T], error) {
	var zero Result[T]
	jsonData, err := codec.Marshal(data)
	if err != nil {
		return zero, &client.EncodeError{ContentType: jsonContentType, Err: err}
	}
	resp, err := Put(baseURL, strings.NewReader(string(jsonData)), append(opts, client.WithContentType(jsonContentType))...)
	if err != nil {
		return zero, err
	}
	var result T
	if err := codec.Unmarshal(resp.Bytes(), &result); err != nil {
		return zero, &client.DecodeError{ContentType: jsonContentType, Err: err}
	}
	return client.NewResult(result, resp), nil
}

// DeleteJSON sends a DELETE request and returns Result[T].
func DeleteJSON[T any](baseURL string, opts ...client.RequestOption) (Result[T], error) {
	var zero Result[T]
	resp, err := Delete(baseURL, opts...)
	if err != nil {
		return zero, err
	}
	var result T
	if err := codec.Unmarshal(resp.Bytes(), &result); err != nil {
		return zero, &client.DecodeError{ContentType: jsonContentType, Err: err}
	}
	return client.NewResult(result, resp), nil
}

// PatchJSON sends a PATCH request with JSON body and returns Result[T].
func PatchJSON[T any](baseURL string, data any, opts ...client.RequestOption) (Result[T], error) {
	var zero Result[T]
	jsonData, err := codec.Marshal(data)
	if err != nil {
		return zero, &client.EncodeError{ContentType: jsonContentType, Err: err}
	}
	resp, err := Patch(baseURL, strings.NewReader(string(jsonData)), append(opts, client.WithContentType(jsonContentType))...)
	if err != nil {
		return zero, err
	}
	var result T
	if err := codec.Unmarshal(resp.Bytes(), &result); err != nil {
		return zero, &client.DecodeError{ContentType: jsonContentType, Err: err}
	}
	return client.NewResult(result, resp), nil
}

// ============================================================================
// Convenience Functions
// ============================================================================
// GetString sends a GET request and returns the response body as a string.
func GetString(baseURL string, opts ...client.RequestOption) (string, error) {
	resp, err := Get(baseURL, opts...)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}

// GetBytes sends a GET request and returns the response body as bytes.
func GetBytes(baseURL string, opts ...client.RequestOption) ([]byte, error) {
	resp, err := Get(baseURL, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Bytes(), nil
}

// ============================================================================
// Helper Functions
// ============================================================================
// prepareBodyForRequest converts body to io.Reader and returns content type.
// This is a wrapper around client.PrepareBody that ignores errors for backward compatibility.
func prepareBodyForRequest(body any) (io.Reader, string) {
	reader, contentType, _ := client.PrepareBody(body)
	return reader, contentType
}

// applyConfigToRequest applies RequestConfig to a Request.
func applyConfigToRequest(req *Request, config *client.RequestConfig) {
	for k, vs := range config.Headers {
		for _, v := range vs {
			req.AddHeader(k, v)
		}
	}
}
