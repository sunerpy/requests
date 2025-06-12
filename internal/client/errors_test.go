package client

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

// Feature: http-client-refactor, Property 13: Error Types Support Unwrapping
// For any error returned by the Client, calling errors.Is() and errors.As()
// SHALL correctly identify the error type and allow access to wrapped errors.
func TestErrorTypes_Property13_Unwrapping(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	// Property: RequestError wrapping and unwrapping
	properties.Property("RequestError supports errors.Is and errors.As", prop.ForAll(
		func(op, url, errMsg string) bool {
			innerErr := errors.New(errMsg)
			reqErr := &RequestError{Op: op, URL: url, Err: innerErr}
			// errors.Is should work
			if !errors.Is(reqErr, innerErr) {
				return false
			}
			// errors.As should work
			var target *RequestError
			if !errors.As(reqErr, &target) {
				return false
			}
			if target.Op != op || target.URL != url {
				return false
			}
			// Unwrap should return inner error
			if reqErr.Unwrap() != innerErr {
				return false
			}
			return true
		},
		gen.AlphaString(),
		gen.AlphaString(),
		gen.AlphaString(),
	))
	// Property: ResponseError wrapping and unwrapping
	properties.Property("ResponseError supports errors.Is and errors.As", prop.ForAll(
		func(statusCode int, status, errMsg string) bool {
			innerErr := errors.New(errMsg)
			respErr := &ResponseError{
				StatusCode: statusCode,
				Status:     status,
				Body:       []byte(errMsg),
				Err:        innerErr,
			}
			// errors.Is should work for inner error
			if !errors.Is(respErr, innerErr) {
				return false
			}
			// errors.As should work
			var target *ResponseError
			if !errors.As(respErr, &target) {
				return false
			}
			if target.StatusCode != statusCode || target.Status != status {
				return false
			}
			// Unwrap should return inner error
			if respErr.Unwrap() != innerErr {
				return false
			}
			return true
		},
		gen.IntRange(100, 599),
		gen.AlphaString(),
		gen.AlphaString(),
	))
	// Property: TimeoutError wrapping and unwrapping
	properties.Property("TimeoutError supports errors.Is and errors.As", prop.ForAll(
		func(op, url string, durationMs int64, errMsg string) bool {
			innerErr := errors.New(errMsg)
			duration := time.Duration(durationMs) * time.Millisecond
			timeoutErr := &TimeoutError{
				Op:       op,
				URL:      url,
				Duration: duration,
				Err:      innerErr,
			}
			// errors.Is should work for inner error
			if !errors.Is(timeoutErr, innerErr) {
				return false
			}
			// errors.As should work
			var target *TimeoutError
			if !errors.As(timeoutErr, &target) {
				return false
			}
			if target.Op != op || target.URL != url || target.Duration != duration {
				return false
			}
			// Timeout() should return true
			if !timeoutErr.Timeout() {
				return false
			}
			// Unwrap should return inner error
			if timeoutErr.Unwrap() != innerErr {
				return false
			}
			return true
		},
		gen.AlphaString(),
		gen.AlphaString(),
		gen.Int64Range(0, 10000),
		gen.AlphaString(),
	))
	// Property: ConnectionError wrapping and unwrapping
	properties.Property("ConnectionError supports errors.Is and errors.As", prop.ForAll(
		func(op, url, errMsg string) bool {
			innerErr := errors.New(errMsg)
			connErr := &ConnectionError{Op: op, URL: url, Err: innerErr}
			// errors.Is should work for inner error
			if !errors.Is(connErr, innerErr) {
				return false
			}
			// errors.As should work
			var target *ConnectionError
			if !errors.As(connErr, &target) {
				return false
			}
			if target.Op != op || target.URL != url {
				return false
			}
			// Unwrap should return inner error
			if connErr.Unwrap() != innerErr {
				return false
			}
			return true
		},
		gen.AlphaString(),
		gen.AlphaString(),
		gen.AlphaString(),
	))
	// Property: RetryError wrapping and unwrapping
	properties.Property("RetryError supports errors.Is and errors.As", prop.ForAll(
		func(attempts int, errMsg string) bool {
			innerErr := errors.New(errMsg)
			retryErr := &RetryError{Attempts: attempts, LastErr: innerErr}
			// errors.Is should work for inner error
			if !errors.Is(retryErr, innerErr) {
				return false
			}
			// errors.Is should work for ErrMaxRetriesExceeded
			if !errors.Is(retryErr, ErrMaxRetriesExceeded) {
				return false
			}
			// errors.As should work
			var target *RetryError
			if !errors.As(retryErr, &target) {
				return false
			}
			if target.Attempts != attempts {
				return false
			}
			// Unwrap should return inner error
			if retryErr.Unwrap() != innerErr {
				return false
			}
			return true
		},
		gen.IntRange(1, 10),
		gen.AlphaString(),
	))
	// Property: Nested error unwrapping
	properties.Property("Nested errors can be unwrapped through chain", prop.ForAll(
		func(op, url, errMsg string) bool {
			// Create a chain: RequestError -> ConnectionError -> base error
			baseErr := errors.New(errMsg)
			connErr := &ConnectionError{Op: "Dial", URL: url, Err: baseErr}
			reqErr := &RequestError{Op: op, URL: url, Err: connErr}
			// Should be able to find ConnectionError in chain
			var targetConn *ConnectionError
			if !errors.As(reqErr, &targetConn) {
				return false
			}
			// Should be able to find base error in chain
			if !errors.Is(reqErr, baseErr) {
				return false
			}
			return true
		},
		gen.AlphaString(),
		gen.AlphaString(),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

// Unit tests for specific error scenarios
func TestRequestError(t *testing.T) {
	innerErr := errors.New("connection refused")
	err := &RequestError{
		Op:  "Send",
		URL: "https://example.com",
		Err: innerErr,
	}
	assert.Contains(t, err.Error(), "Send")
	assert.Contains(t, err.Error(), "https://example.com")
	assert.Contains(t, err.Error(), "connection refused")
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestResponseError(t *testing.T) {
	err := &ResponseError{
		StatusCode: 404,
		Status:     "Not Found",
		Body:       []byte("page not found"),
	}
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "Not Found")
	assert.Equal(t, 404, err.StatusCode)
	assert.Equal(t, []byte("page not found"), err.Body)
}

func TestTimeoutError(t *testing.T) {
	innerErr := errors.New("context deadline exceeded")
	err := &TimeoutError{
		Op:       "Read",
		URL:      "https://example.com",
		Duration: 5 * time.Second,
		Err:      innerErr,
	}
	assert.Contains(t, err.Error(), "Read")
	assert.Contains(t, err.Error(), "5s")
	assert.True(t, err.Timeout())
	assert.True(t, err.Temporary())
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestConnectionError(t *testing.T) {
	innerErr := errors.New("no route to host")
	err := &ConnectionError{
		Op:  "Dial",
		URL: "https://example.com",
		Err: innerErr,
	}
	assert.Contains(t, err.Error(), "Dial")
	assert.Contains(t, err.Error(), "no route to host")
	assert.True(t, err.Temporary())
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestRetryError(t *testing.T) {
	innerErr := errors.New("server unavailable")
	err := &RetryError{
		Attempts: 3,
		LastErr:  innerErr,
	}
	assert.Contains(t, err.Error(), "3 attempts")
	assert.Contains(t, err.Error(), "server unavailable")
	assert.True(t, errors.Is(err, ErrMaxRetriesExceeded))
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestDecodeError(t *testing.T) {
	innerErr := errors.New("invalid json")
	err := &DecodeError{
		ContentType: "application/json",
		Err:         innerErr,
	}
	assert.Contains(t, err.Error(), "application/json")
	assert.Contains(t, err.Error(), "invalid json")
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestEncodeError(t *testing.T) {
	innerErr := errors.New("unsupported type")
	err := &EncodeError{
		ContentType: "application/json",
		Err:         innerErr,
	}
	assert.Contains(t, err.Error(), "application/json")
	assert.Contains(t, err.Error(), "unsupported type")
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestIsTimeout(t *testing.T) {
	timeoutErr := &TimeoutError{Op: "Read", Err: errors.New("timeout")}
	assert.True(t, IsTimeout(timeoutErr))
	regularErr := errors.New("regular error")
	assert.False(t, IsTimeout(regularErr))
	// Test with standard library timeout interface
	stdTimeoutErr := &mockTimeoutError{isTimeout: true}
	assert.True(t, IsTimeout(stdTimeoutErr))
	stdNonTimeoutErr := &mockTimeoutError{isTimeout: false}
	assert.False(t, IsTimeout(stdNonTimeoutErr))
}

// mockTimeoutError implements the timeout interface
type mockTimeoutError struct {
	isTimeout bool
}

func (e *mockTimeoutError) Error() string {
	return "mock timeout error"
}

func (e *mockTimeoutError) Timeout() bool {
	return e.isTimeout
}

func TestIsConnectionError(t *testing.T) {
	connErr := &ConnectionError{Op: "Dial", Err: errors.New("refused")}
	assert.True(t, IsConnectionError(connErr))
	regularErr := errors.New("regular error")
	assert.False(t, IsConnectionError(regularErr))
}

func TestIsResponseError(t *testing.T) {
	respErr := &ResponseError{StatusCode: 500}
	assert.True(t, IsResponseError(respErr))
	regularErr := errors.New("regular error")
	assert.False(t, IsResponseError(regularErr))
}

func TestErrorIs_SameType(t *testing.T) {
	// Test that errors.Is works for same error types
	err1 := &RequestError{Op: "Send"}
	err2 := &RequestError{Op: "Send"}
	assert.True(t, errors.Is(err1, err2))
	// Different Op should still match if target Op is empty
	err3 := &RequestError{Op: ""}
	assert.True(t, errors.Is(err1, err3))
}

func TestWrappedErrorChain(t *testing.T) {
	// Create a complex error chain
	baseErr := errors.New("base error")
	connErr := &ConnectionError{Op: "Dial", URL: "https://example.com", Err: baseErr}
	reqErr := &RequestError{Op: "Send", URL: "https://example.com", Err: connErr}
	wrappedErr := fmt.Errorf("wrapped: %w", reqErr)
	// Should be able to find all errors in chain
	var targetReq *RequestError
	assert.True(t, errors.As(wrappedErr, &targetReq))
	var targetConn *ConnectionError
	assert.True(t, errors.As(wrappedErr, &targetConn))
	assert.True(t, errors.Is(wrappedErr, baseErr))
}

// Additional tests for Is() methods to improve coverage
func TestRequestError_Is(t *testing.T) {
	err := &RequestError{Op: "Send", URL: "https://example.com"}
	// Same Op should match
	target := &RequestError{Op: "Send"}
	if !err.Is(target) {
		t.Error("Should match same Op")
	}
	// Empty Op in target should match any
	emptyTarget := &RequestError{Op: ""}
	if !err.Is(emptyTarget) {
		t.Error("Empty Op target should match any")
	}
	// Different Op should not match
	diffTarget := &RequestError{Op: "Build"}
	if err.Is(diffTarget) {
		t.Error("Different Op should not match")
	}
	// Non-RequestError should not match
	if err.Is(errors.New("other")) {
		t.Error("Should not match non-RequestError")
	}
}

func TestRequestError_ErrorWithoutURL(t *testing.T) {
	err := &RequestError{Op: "Build", Err: errors.New("test")}
	msg := err.Error()
	if !strings.Contains(msg, "Build") {
		t.Error("Error message should contain Op")
	}
	if strings.Contains(msg, "https://") {
		t.Error("Error message should not contain URL when not set")
	}
}

func TestResponseError_Is(t *testing.T) {
	err := &ResponseError{StatusCode: 404, Status: "Not Found"}
	// Same status code should match
	target := &ResponseError{StatusCode: 404}
	if !err.Is(target) {
		t.Error("Should match same status code")
	}
	// Zero status code in target should match any
	zeroTarget := &ResponseError{StatusCode: 0}
	if !err.Is(zeroTarget) {
		t.Error("Zero status code target should match any")
	}
	// Different status code should not match
	diffTarget := &ResponseError{StatusCode: 500}
	if err.Is(diffTarget) {
		t.Error("Different status code should not match")
	}
	// Non-ResponseError should not match
	if err.Is(errors.New("other")) {
		t.Error("Should not match non-ResponseError")
	}
}

func TestTimeoutError_Is(t *testing.T) {
	err := &TimeoutError{Op: "Read", URL: "https://example.com"}
	// Same Op should match
	target := &TimeoutError{Op: "Read"}
	if !err.Is(target) {
		t.Error("Should match same Op")
	}
	// Empty Op in target should match any
	emptyTarget := &TimeoutError{Op: ""}
	if !err.Is(emptyTarget) {
		t.Error("Empty Op target should match any")
	}
	// Different Op should not match
	diffTarget := &TimeoutError{Op: "Write"}
	if err.Is(diffTarget) {
		t.Error("Different Op should not match")
	}
	// Non-TimeoutError should not match
	if err.Is(errors.New("other")) {
		t.Error("Should not match non-TimeoutError")
	}
}

func TestConnectionError_Is(t *testing.T) {
	err := &ConnectionError{Op: "Dial", URL: "https://example.com"}
	// Same Op should match
	target := &ConnectionError{Op: "Dial"}
	if !err.Is(target) {
		t.Error("Should match same Op")
	}
	// Empty Op in target should match any
	emptyTarget := &ConnectionError{Op: ""}
	if !err.Is(emptyTarget) {
		t.Error("Empty Op target should match any")
	}
	// Different Op should not match
	diffTarget := &ConnectionError{Op: "TLS"}
	if err.Is(diffTarget) {
		t.Error("Different Op should not match")
	}
	// Non-ConnectionError should not match
	if err.Is(errors.New("other")) {
		t.Error("Should not match non-ConnectionError")
	}
}

func TestDecodeError_Is(t *testing.T) {
	err := &DecodeError{ContentType: "application/json"}
	// Same content type should match
	target := &DecodeError{ContentType: "application/json"}
	if !err.Is(target) {
		t.Error("Should match same content type")
	}
	// Empty content type in target should match any
	emptyTarget := &DecodeError{ContentType: ""}
	if !err.Is(emptyTarget) {
		t.Error("Empty content type target should match any")
	}
	// Different content type should not match
	diffTarget := &DecodeError{ContentType: "application/xml"}
	if err.Is(diffTarget) {
		t.Error("Different content type should not match")
	}
	// Non-DecodeError should not match
	if err.Is(errors.New("other")) {
		t.Error("Should not match non-DecodeError")
	}
}

func TestEncodeError_Is(t *testing.T) {
	err := &EncodeError{ContentType: "application/json"}
	// Same content type should match
	target := &EncodeError{ContentType: "application/json"}
	if !err.Is(target) {
		t.Error("Should match same content type")
	}
	// Empty content type in target should match any
	emptyTarget := &EncodeError{ContentType: ""}
	if !err.Is(emptyTarget) {
		t.Error("Empty content type target should match any")
	}
	// Different content type should not match
	diffTarget := &EncodeError{ContentType: "application/xml"}
	if err.Is(diffTarget) {
		t.Error("Different content type should not match")
	}
	// Non-EncodeError should not match
	if err.Is(errors.New("other")) {
		t.Error("Should not match non-EncodeError")
	}
}

func TestIsTemporary_NonTemporaryError(t *testing.T) {
	// Regular error without Temporary() method
	err := errors.New("regular error")
	if IsTemporary(err) {
		t.Error("Regular error should not be temporary")
	}
}

func TestIsTemporary_TemporaryError(t *testing.T) {
	// ConnectionError implements Temporary()
	err := &ConnectionError{Op: "Dial", Err: errors.New("refused")}
	if !IsTemporary(err) {
		t.Error("ConnectionError should be temporary")
	}
	// TimeoutError implements Temporary()
	timeoutErr := &TimeoutError{Op: "Read", Err: errors.New("timeout")}
	if !IsTemporary(timeoutErr) {
		t.Error("TimeoutError should be temporary")
	}
}

// Feature: http-client-refactor
// Property 14: ResponseError Contains Full Context
// For any ResponseError, it SHALL contain the status code, response body, and original Response object.
// Validates: Requirements 7.4
func TestProperty14_ResponseErrorContainsFullContext(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("ResponseError contains status code, body, and response", prop.ForAll(
		func(statusCode int, bodyContent string) bool {
			if statusCode < 100 || statusCode > 599 {
				return true
			}
			body := []byte(bodyContent)
			resp := CreateMockResponse(statusCode, body, nil)
			respErr := &ResponseError{
				StatusCode: statusCode,
				Status:     resp.Status,
				Body:       body,
				Response:   resp,
			}
			// Verify status code
			if respErr.StatusCode != statusCode {
				return false
			}
			// Verify body
			if string(respErr.Body) != bodyContent {
				return false
			}
			// Verify response
			if respErr.Response != resp {
				return false
			}
			// Verify error message contains status code
			errMsg := respErr.Error()
			if !strings.Contains(errMsg, fmt.Sprintf("%d", statusCode)) {
				return false
			}
			return true
		},
		gen.IntRange(100, 599),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

func TestResponseError_FullContext(t *testing.T) {
	body := []byte(`{"error": "not found"}`)
	resp := CreateMockResponse(404, body, nil)
	respErr := &ResponseError{
		StatusCode: 404,
		Status:     "Not Found",
		Body:       body,
		Response:   resp,
		Err:        errors.New("resource not found"),
	}
	// Verify all fields are accessible
	if respErr.StatusCode != 404 {
		t.Error("StatusCode not set")
	}
	if respErr.Status != "Not Found" {
		t.Error("Status not set")
	}
	if string(respErr.Body) != `{"error": "not found"}` {
		t.Error("Body not set")
	}
	if respErr.Response == nil {
		t.Error("Response not set")
	}
	if respErr.Response.StatusCode != 404 {
		t.Error("Response status code mismatch")
	}
	// Verify error message
	errMsg := respErr.Error()
	if !strings.Contains(errMsg, "404") {
		t.Error("Error message should contain status code")
	}
	if !strings.Contains(errMsg, "resource not found") {
		t.Error("Error message should contain underlying error")
	}
}
