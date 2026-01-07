package client

import (
	"errors"
	"fmt"
	"time"

	"github.com/sunerpy/requests/internal/models"
)

// Sentinel errors
var (
	ErrNilResponse        = errors.New("response is nil")
	ErrNilRequest         = errors.New("request is nil")
	ErrMissingURL         = errors.New("URL is required")
	ErrMissingMethod      = errors.New("method is required")
	ErrInvalidMethod      = errors.New("invalid HTTP method")
	ErrInvalidURL         = errors.New("invalid URL")
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")
	ErrPanic              = errors.New("panic recovered in middleware")
)

// RequestError represents an error during request building or sending.
type RequestError struct {
	Op  string // Operation name (e.g., "Build", "Send", "Encode")
	URL string // URL of the request
	Err error  // Underlying error
}

// Error implements the error interface.
func (e *RequestError) Error() string {
	if e.URL != "" {
		return fmt.Sprintf("request error [%s] %s: %v", e.Op, e.URL, e.Err)
	}
	return fmt.Sprintf("request error [%s]: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error.
func (e *RequestError) Unwrap() error {
	return e.Err
}

// Is reports whether target matches this error.
func (e *RequestError) Is(target error) bool {
	t, ok := target.(*RequestError)
	if !ok {
		return false
	}
	return e.Op == t.Op || t.Op == ""
}

// ResponseError represents an error from the server response.
type ResponseError struct {
	StatusCode int
	Status     string
	Body       []byte
	Response   *models.Response
	Err        error
}

// Error implements the error interface.
func (e *ResponseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("response error [%d %s]: %v", e.StatusCode, e.Status, e.Err)
	}
	return fmt.Sprintf("response error [%d %s]", e.StatusCode, e.Status)
}

// Unwrap returns the underlying error.
func (e *ResponseError) Unwrap() error {
	return e.Err
}

// Is reports whether target matches this error.
func (e *ResponseError) Is(target error) bool {
	t, ok := target.(*ResponseError)
	if !ok {
		return false
	}
	return e.StatusCode == t.StatusCode || t.StatusCode == 0
}

// TimeoutError represents a timeout error.
type TimeoutError struct {
	Op       string        // Operation that timed out
	URL      string        // URL of the request
	Duration time.Duration // How long we waited
	Err      error         // Underlying error
}

// Error implements the error interface.
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout error [%s] %s after %v: %v", e.Op, e.URL, e.Duration, e.Err)
}

// Unwrap returns the underlying error.
func (e *TimeoutError) Unwrap() error {
	return e.Err
}

// Timeout returns true, implementing net.Error interface.
func (e *TimeoutError) Timeout() bool {
	return true
}

// Temporary returns true, implementing net.Error interface.
func (e *TimeoutError) Temporary() bool {
	return true
}

// Is reports whether target matches this error.
func (e *TimeoutError) Is(target error) bool {
	t, ok := target.(*TimeoutError)
	if !ok {
		return false
	}
	return e.Op == t.Op || t.Op == ""
}

// ConnectionError represents a network connection error.
type ConnectionError struct {
	Op  string // Operation (e.g., "Dial", "TLS", "DNS")
	URL string // URL of the request
	Err error  // Underlying error
}

// Error implements the error interface.
func (e *ConnectionError) Error() string {
	return fmt.Sprintf("connection error [%s] %s: %v", e.Op, e.URL, e.Err)
}

// Unwrap returns the underlying error.
func (e *ConnectionError) Unwrap() error {
	return e.Err
}

// Is reports whether target matches this error.
func (e *ConnectionError) Is(target error) bool {
	t, ok := target.(*ConnectionError)
	if !ok {
		return false
	}
	return e.Op == t.Op || t.Op == ""
}

// Temporary returns true if the error is temporary.
func (e *ConnectionError) Temporary() bool {
	return true
}

// RetryError represents an error after all retries have been exhausted.
type RetryError struct {
	Attempts int   // Number of attempts made
	LastErr  error // The last error encountered
}

// Error implements the error interface.
func (e *RetryError) Error() string {
	return fmt.Sprintf("max retries exceeded after %d attempts: %v", e.Attempts, e.LastErr)
}

// Unwrap returns the underlying error.
func (e *RetryError) Unwrap() error {
	return e.LastErr
}

// Is reports whether target matches this error.
func (e *RetryError) Is(target error) bool {
	if errors.Is(target, ErrMaxRetriesExceeded) {
		return true
	}
	_, ok := target.(*RetryError)
	return ok
}

// DecodeError represents an error during response decoding.
type DecodeError struct {
	ContentType string // Content type being decoded
	Err         error  // Underlying error
}

// Error implements the error interface.
func (e *DecodeError) Error() string {
	return fmt.Sprintf("decode error [%s]: %v", e.ContentType, e.Err)
}

// Unwrap returns the underlying error.
func (e *DecodeError) Unwrap() error {
	return e.Err
}

// Is reports whether target matches this error.
func (e *DecodeError) Is(target error) bool {
	t, ok := target.(*DecodeError)
	if !ok {
		return false
	}
	return e.ContentType == t.ContentType || t.ContentType == ""
}

// EncodeError represents an error during request encoding.
type EncodeError struct {
	ContentType string // Content type being encoded
	Err         error  // Underlying error
}

// Error implements the error interface.
func (e *EncodeError) Error() string {
	return fmt.Sprintf("encode error [%s]: %v", e.ContentType, e.Err)
}

// Unwrap returns the underlying error.
func (e *EncodeError) Unwrap() error {
	return e.Err
}

// Is reports whether target matches this error.
func (e *EncodeError) Is(target error) bool {
	t, ok := target.(*EncodeError)
	if !ok {
		return false
	}
	return e.ContentType == t.ContentType || t.ContentType == ""
}

// IsTimeout checks if an error is a timeout error.
func IsTimeout(err error) bool {
	var timeoutErr *TimeoutError
	if errors.As(err, &timeoutErr) {
		return true
	}
	// Also check for standard library timeout errors
	type timeout interface {
		Timeout() bool
	}
	if t, ok := err.(timeout); ok {
		return t.Timeout()
	}
	return false
}

// IsTemporary checks if an error is temporary and can be retried.
func IsTemporary(err error) bool {
	type temporary interface {
		Temporary() bool
	}
	if t, ok := err.(temporary); ok {
		return t.Temporary()
	}
	return false
}

// IsConnectionError checks if an error is a connection error.
func IsConnectionError(err error) bool {
	var connErr *ConnectionError
	return errors.As(err, &connErr)
}

// IsResponseError checks if an error is a response error.
func IsResponseError(err error) bool {
	var respErr *ResponseError
	return errors.As(err, &respErr)
}
