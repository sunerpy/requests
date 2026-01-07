package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/sunerpy/requests/internal/models"
)

// HTTPClient is the interface for executing HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient is the default HTTP client used for requests.
var DefaultHTTPClient HTTPClient = &http.Client{
	Timeout: 30 * time.Second,
}

// SetDefaultClient sets the default HTTP client.
func SetDefaultClient(client HTTPClient) {
	DefaultHTTPClient = client
}

// doRequest executes an HTTP request and returns the response.
func doRequest(ctx context.Context, method Method, rawURL string, body io.Reader, opts ...RequestOption) (*models.Response, error) {
	config := NewRequestConfig()
	config.Apply(opts...)
	// Parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, &RequestError{Op: "ParseURL", URL: rawURL, Err: err}
	}
	// Apply query parameters
	if len(config.Query) > 0 {
		q := parsedURL.Query()
		for k, vs := range config.Query {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
		parsedURL.RawQuery = q.Encode()
	}
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method.String(), parsedURL.String(), body)
	if err != nil {
		return nil, &RequestError{Op: "NewRequest", URL: rawURL, Err: err}
	}
	// Apply headers
	for k, vs := range config.Headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	// Apply auth
	if config.BasicAuth != nil {
		req.SetBasicAuth(config.BasicAuth.Username, config.BasicAuth.Password)
	}
	if config.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+config.BearerToken)
	}
	// Execute request
	httpResp, err := DefaultHTTPClient.Do(req)
	if err != nil {
		return nil, &RequestError{Op: "Do", URL: rawURL, Err: err}
	}
	return models.NewResponse(httpResp, parsedURL.String())
}

// ============================================================================
// Basic HTTP Methods - return (*models.Response, error)
// ============================================================================
// Get performs a GET request and returns the response.
func Get(rawURL string, opts ...RequestOption) (*models.Response, error) {
	return doRequest(context.Background(), MethodGet, rawURL, nil, opts...)
}

// GetWithContext performs a GET request with context.
func GetWithContext(ctx context.Context, rawURL string, opts ...RequestOption) (*models.Response, error) {
	return doRequest(ctx, MethodGet, rawURL, nil, opts...)
}

// Post performs a POST request and returns the response.
func Post(rawURL string, body any, opts ...RequestOption) (*models.Response, error) {
	return PostWithContext(context.Background(), rawURL, body, opts...)
}

// PostWithContext performs a POST request with context.
func PostWithContext(ctx context.Context, rawURL string, body any, opts ...RequestOption) (*models.Response, error) {
	reader, contentType, err := PrepareBody(body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		opts = append(opts, WithContentType(contentType))
	}
	return doRequest(ctx, MethodPost, rawURL, reader, opts...)
}

// Put performs a PUT request and returns the response.
func Put(rawURL string, body any, opts ...RequestOption) (*models.Response, error) {
	return PutWithContext(context.Background(), rawURL, body, opts...)
}

// PutWithContext performs a PUT request with context.
func PutWithContext(ctx context.Context, rawURL string, body any, opts ...RequestOption) (*models.Response, error) {
	reader, contentType, err := PrepareBody(body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		opts = append(opts, WithContentType(contentType))
	}
	return doRequest(ctx, MethodPut, rawURL, reader, opts...)
}

// Delete performs a DELETE request and returns the response.
func Delete(rawURL string, opts ...RequestOption) (*models.Response, error) {
	return doRequest(context.Background(), MethodDelete, rawURL, nil, opts...)
}

// DeleteWithContext performs a DELETE request with context.
func DeleteWithContext(ctx context.Context, rawURL string, opts ...RequestOption) (*models.Response, error) {
	return doRequest(ctx, MethodDelete, rawURL, nil, opts...)
}

// Patch performs a PATCH request and returns the response.
func Patch(rawURL string, body any, opts ...RequestOption) (*models.Response, error) {
	return PatchWithContext(context.Background(), rawURL, body, opts...)
}

// PatchWithContext performs a PATCH request with context.
func PatchWithContext(ctx context.Context, rawURL string, body any, opts ...RequestOption) (*models.Response, error) {
	reader, contentType, err := PrepareBody(body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		opts = append(opts, WithContentType(contentType))
	}
	return doRequest(ctx, MethodPatch, rawURL, reader, opts...)
}

// Head performs a HEAD request and returns the response.
func Head(rawURL string, opts ...RequestOption) (*models.Response, error) {
	return doRequest(context.Background(), MethodHead, rawURL, nil, opts...)
}

// HeadWithContext performs a HEAD request with context.
func HeadWithContext(ctx context.Context, rawURL string, opts ...RequestOption) (*models.Response, error) {
	return doRequest(ctx, MethodHead, rawURL, nil, opts...)
}

// Options performs an OPTIONS request and returns the response.
func Options(rawURL string, opts ...RequestOption) (*models.Response, error) {
	return doRequest(context.Background(), MethodOptions, rawURL, nil, opts...)
}

// OptionsWithContext performs an OPTIONS request with context.
func OptionsWithContext(ctx context.Context, rawURL string, opts ...RequestOption) (*models.Response, error) {
	return doRequest(ctx, MethodOptions, rawURL, nil, opts...)
}

// PrepareBody converts the body to an io.Reader and returns the content type.
// This is exported for use by the main package.
func PrepareBody(body any) (io.Reader, string, error) {
	if body == nil {
		return nil, "", nil
	}
	switch v := body.(type) {
	case io.Reader:
		return v, "", nil
	case []byte:
		return bytes.NewReader(v), "", nil
	case string:
		return bytes.NewReader([]byte(v)), "", nil
	case *url.Values:
		if v != nil {
			return bytes.NewReader([]byte(v.Encode())), "application/x-www-form-urlencoded", nil
		}
		return nil, "", nil
	case url.Values:
		return bytes.NewReader([]byte(v.Encode())), "application/x-www-form-urlencoded", nil
	default:
		// Check if it has an Encode() method (like custom url.Values types)
		// Use reflection to check for nil pointer before calling Encode()
		if encoder, ok := body.(interface{ Encode() string }); ok {
			// Check if the underlying value is nil using reflection
			rv := reflect.ValueOf(body)
			if rv.Kind() == reflect.Ptr && rv.IsNil() {
				return nil, "", nil
			}
			return bytes.NewReader([]byte(encoder.Encode())), "application/x-www-form-urlencoded", nil
		}
		// Assume JSON
		data, err := json.Marshal(v)
		if err != nil {
			return nil, "", &EncodeError{ContentType: "application/json", Err: err}
		}
		return bytes.NewReader(data), "application/json", nil
	}
}

// ============================================================================
// Generic HTTP Methods - return (Result[T], error)
// ============================================================================
// GetJSON performs a GET request and returns Result[T] with parsed JSON response.
func GetJSON[T any](rawURL string, opts ...RequestOption) (Result[T], error) {
	return GetJSONWithContext[T](context.Background(), rawURL, opts...)
}

// GetJSONWithContext performs a GET request with context and returns Result[T].
func GetJSONWithContext[T any](ctx context.Context, rawURL string, opts ...RequestOption) (Result[T], error) {
	var zero Result[T]
	resp, err := doRequest(ctx, MethodGet, rawURL, nil, opts...)
	if err != nil {
		return zero, err
	}
	data, err := models.JSON[T](resp)
	if err != nil {
		return zero, err
	}
	return NewResult(data, resp), nil
}

// PostJSON performs a POST request and returns Result[T] with parsed JSON response.
func PostJSON[T any](rawURL string, body any, opts ...RequestOption) (Result[T], error) {
	return PostJSONWithContext[T](context.Background(), rawURL, body, opts...)
}

// PostJSONWithContext performs a POST request with context and returns Result[T].
func PostJSONWithContext[T any](ctx context.Context, rawURL string, body any, opts ...RequestOption) (Result[T], error) {
	var zero Result[T]
	reader, contentType, err := PrepareBody(body)
	if err != nil {
		return zero, err
	}
	if contentType != "" {
		opts = append(opts, WithContentType(contentType))
	}
	resp, err := doRequest(ctx, MethodPost, rawURL, reader, opts...)
	if err != nil {
		return zero, err
	}
	data, err := models.JSON[T](resp)
	if err != nil {
		return zero, err
	}
	return NewResult(data, resp), nil
}

// PutJSON performs a PUT request and returns Result[T] with parsed JSON response.
func PutJSON[T any](rawURL string, body any, opts ...RequestOption) (Result[T], error) {
	return PutJSONWithContext[T](context.Background(), rawURL, body, opts...)
}

// PutJSONWithContext performs a PUT request with context and returns Result[T].
func PutJSONWithContext[T any](ctx context.Context, rawURL string, body any, opts ...RequestOption) (Result[T], error) {
	var zero Result[T]
	reader, contentType, err := PrepareBody(body)
	if err != nil {
		return zero, err
	}
	if contentType != "" {
		opts = append(opts, WithContentType(contentType))
	}
	resp, err := doRequest(ctx, MethodPut, rawURL, reader, opts...)
	if err != nil {
		return zero, err
	}
	data, err := models.JSON[T](resp)
	if err != nil {
		return zero, err
	}
	return NewResult(data, resp), nil
}

// DeleteJSON performs a DELETE request and returns Result[T] with parsed JSON response.
func DeleteJSON[T any](rawURL string, opts ...RequestOption) (Result[T], error) {
	return DeleteJSONWithContext[T](context.Background(), rawURL, opts...)
}

// DeleteJSONWithContext performs a DELETE request with context and returns Result[T].
func DeleteJSONWithContext[T any](ctx context.Context, rawURL string, opts ...RequestOption) (Result[T], error) {
	var zero Result[T]
	resp, err := doRequest(ctx, MethodDelete, rawURL, nil, opts...)
	if err != nil {
		return zero, err
	}
	data, err := models.JSON[T](resp)
	if err != nil {
		return zero, err
	}
	return NewResult(data, resp), nil
}

// PatchJSON performs a PATCH request and returns Result[T] with parsed JSON response.
func PatchJSON[T any](rawURL string, body any, opts ...RequestOption) (Result[T], error) {
	return PatchJSONWithContext[T](context.Background(), rawURL, body, opts...)
}

// PatchJSONWithContext performs a PATCH request with context and returns Result[T].
func PatchJSONWithContext[T any](ctx context.Context, rawURL string, body any, opts ...RequestOption) (Result[T], error) {
	var zero Result[T]
	reader, contentType, err := PrepareBody(body)
	if err != nil {
		return zero, err
	}
	if contentType != "" {
		opts = append(opts, WithContentType(contentType))
	}
	resp, err := doRequest(ctx, MethodPatch, rawURL, reader, opts...)
	if err != nil {
		return zero, err
	}
	data, err := models.JSON[T](resp)
	if err != nil {
		return zero, err
	}
	return NewResult(data, resp), nil
}

// ============================================================================
// Generic XML Methods - return (Result[T], error)
// ============================================================================
// GetXML performs a GET request and returns Result[T] with parsed XML response.
func GetXML[T any](rawURL string, opts ...RequestOption) (Result[T], error) {
	return GetXMLWithContext[T](context.Background(), rawURL, opts...)
}

// GetXMLWithContext performs a GET request with context and returns Result[T].
func GetXMLWithContext[T any](ctx context.Context, rawURL string, opts ...RequestOption) (Result[T], error) {
	var zero Result[T]
	resp, err := doRequest(ctx, MethodGet, rawURL, nil, opts...)
	if err != nil {
		return zero, err
	}
	data, err := models.XML[T](resp)
	if err != nil {
		return zero, err
	}
	return NewResult(data, resp), nil
}

// PostXML performs a POST request and returns Result[T] with parsed XML response.
func PostXML[T any](rawURL string, body any, opts ...RequestOption) (Result[T], error) {
	return PostXMLWithContext[T](context.Background(), rawURL, body, opts...)
}

// PostXMLWithContext performs a POST request with context and returns Result[T].
func PostXMLWithContext[T any](ctx context.Context, rawURL string, body any, opts ...RequestOption) (Result[T], error) {
	var zero Result[T]
	reader, contentType, err := PrepareBody(body)
	if err != nil {
		return zero, err
	}
	if contentType != "" {
		opts = append(opts, WithContentType(contentType))
	}
	resp, err := doRequest(ctx, MethodPost, rawURL, reader, opts...)
	if err != nil {
		return zero, err
	}
	data, err := models.XML[T](resp)
	if err != nil {
		return zero, err
	}
	return NewResult(data, resp), nil
}

// ============================================================================
// Convenience Methods
// ============================================================================
// GetString performs a GET request and returns the response body as string.
func GetString(rawURL string, opts ...RequestOption) (string, error) {
	resp, err := Get(rawURL, opts...)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}

// GetBytes performs a GET request and returns the response body as bytes.
func GetBytes(rawURL string, opts ...RequestOption) ([]byte, error) {
	resp, err := Get(rawURL, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Bytes(), nil
}
