package client

import "net/http"

// Result encapsulates both the parsed response data and response metadata.
// It provides a clean API that returns only 2 values (Result[T], error) instead of 3.
type Result[T any] struct {
	data     T
	response *Response
}

// NewResult creates a new Result with the given data and response.
func NewResult[T any](data T, response *Response) Result[T] {
	return Result[T]{
		data:     data,
		response: response,
	}
}

// Data returns the parsed response body as type T.
func (r Result[T]) Data() T {
	return r.data
}

// Response returns the underlying Response object.
func (r Result[T]) Response() *Response {
	return r.response
}

// StatusCode returns the HTTP status code.
func (r Result[T]) StatusCode() int {
	if r.response == nil {
		return 0
	}
	return r.response.StatusCode
}

// Headers returns the response headers.
func (r Result[T]) Headers() http.Header {
	if r.response == nil {
		return nil
	}
	return r.response.Headers
}

// IsSuccess returns true if the status code is 2xx.
func (r Result[T]) IsSuccess() bool {
	code := r.StatusCode()
	return code >= 200 && code < 300
}

// IsError returns true if the status code is 4xx or 5xx.
func (r Result[T]) IsError() bool {
	code := r.StatusCode()
	return code >= 400
}

// IsClientError returns true if the status code is 4xx.
func (r Result[T]) IsClientError() bool {
	code := r.StatusCode()
	return code >= 400 && code < 500
}

// IsServerError returns true if the status code is 5xx.
func (r Result[T]) IsServerError() bool {
	code := r.StatusCode()
	return code >= 500
}

// Cookies returns the response cookies.
func (r Result[T]) Cookies() []*http.Cookie {
	if r.response == nil {
		return nil
	}
	return r.response.Cookies
}

// ContentType returns the Content-Type header value.
func (r Result[T]) ContentType() string {
	if r.response == nil {
		return ""
	}
	return r.response.ContentType()
}

// Text returns the response body as a string.
func (r Result[T]) Text() string {
	if r.response == nil {
		return ""
	}
	return r.response.Text()
}

// Bytes returns the response body as bytes.
func (r Result[T]) Bytes() []byte {
	if r.response == nil {
		return nil
	}
	return r.response.Bytes()
}

// URL returns the final URL after redirects.
func (r Result[T]) URL() string {
	if r.response == nil {
		return ""
	}
	return r.response.GetURL()
}
