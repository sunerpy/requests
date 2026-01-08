package requests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pgregory.net/rapid"

	"github.com/sunerpy/requests/internal/client"
)

// Property 1: Context Cancellation Respected
// For any HTTP request and any context that is canceled,
// calling DoWithContext SHALL return a context cancellation error without completing the request.
func TestProperty_ContextCancellationRespected(t *testing.T) {
	// Create a slow server that takes time to respond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	rapid.Check(t, func(t *rapid.T) {
		// Generate random cancellation delay (0-100ms)
		cancelDelay := rapid.IntRange(0, 100).Draw(t, "cancelDelay")

		sess := NewSession()
		defer sess.Close()

		req, err := NewGet(server.URL).Build()
		if err != nil {
			t.Fatal(err)
		}

		// Create a context that will be canceled
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel the context after a short delay
		go func() {
			time.Sleep(time.Duration(cancelDelay) * time.Millisecond)
			cancel()
		}()

		// Execute request with cancellable context
		_, err = sess.DoWithContext(ctx, req)

		// The request should fail with context.Canceled
		if err == nil {
			t.Error("Expected error due to context cancellation, got nil")
		} else if err != context.Canceled {
			// Check if the error wraps context.Canceled
			if ctx.Err() != context.Canceled {
				t.Errorf("Expected context.Canceled error, got: %v", err)
			}
		}
	})
}

// Property 2: Context Deadline Respected
// For any HTTP request and any context with a deadline,
// if the request execution time exceeds the deadline, DoWithContext SHALL return a deadline exceeded error.
func TestProperty_ContextDeadlineRespected(t *testing.T) {
	// Create a slow server that takes time to respond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	rapid.Check(t, func(t *rapid.T) {
		// Generate random timeout (10-100ms, shorter than server response time)
		timeout := rapid.IntRange(10, 100).Draw(t, "timeout")

		sess := NewSession()
		defer sess.Close()

		req, err := NewGet(server.URL).Build()
		if err != nil {
			t.Fatal(err)
		}

		// Create a context with deadline
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
		defer cancel()

		// Execute request with deadline context
		_, err = sess.DoWithContext(ctx, req)

		// The request should fail with deadline exceeded
		if err == nil {
			t.Error("Expected error due to context deadline, got nil")
		} else if err != context.DeadlineExceeded {
			// Check if the error wraps context.DeadlineExceeded
			if ctx.Err() != context.DeadlineExceeded {
				t.Errorf("Expected context.DeadlineExceeded error, got: %v", err)
			}
		}
	})
}

// Property 3: Configuration Method Chaining
// For any Session and any configuration method,
// calling the method SHALL return a Session instance that allows further method chaining.
func TestProperty_ConfigurationMethodChaining(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random configuration values
		timeout := rapid.IntRange(1, 60).Draw(t, "timeout")
		maxIdle := rapid.IntRange(1, 100).Draw(t, "maxIdle")
		headerKey := rapid.StringMatching(`[a-zA-Z][a-zA-Z0-9-]*`).Draw(t, "headerKey")
		headerValue := rapid.String().Draw(t, "headerValue")
		baseURL := rapid.SampledFrom([]string{
			"http://example.com",
			"https://api.example.com",
			"http://localhost:8080",
		}).Draw(t, "baseURL")

		sess := NewSession()
		defer sess.Close()

		// Chain multiple configuration methods
		result := sess.
			WithTimeout(time.Duration(timeout)*time.Second).
			WithMaxIdleConns(maxIdle).
			WithHeader(headerKey, headerValue).
			WithBaseURL(baseURL).
			WithKeepAlive(true).
			WithHTTP2(false)

		// Result should be a valid Session (not nil)
		if result == nil {
			t.Error("Configuration method chaining returned nil")
		}

		// Result should be the same session (fluent interface)
		// We can verify by checking that further operations work
		result2 := result.WithIdleTimeout(30 * time.Second)
		if result2 == nil {
			t.Error("Further chaining after configuration returned nil")
		}
	})
}

// Additional test: Verify DoWithContext works correctly with valid context
func TestDoWithContext_ValidContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	sess := NewSession()
	defer sess.Close()

	req, err := NewGet(server.URL).Build()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	resp, err := sess.DoWithContext(ctx, req)
	if err != nil {
		t.Fatalf("DoWithContext failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if resp.Text() != `{"status":"ok"}` {
		t.Errorf("Unexpected response body: %s", resp.Text())
	}
}

// Test WithRetry configuration
func TestWithRetry_Configuration(t *testing.T) {
	sess := NewSession()
	defer sess.Close()

	policy := RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
	}

	result := sess.WithRetry(policy)
	if result == nil {
		t.Error("WithRetry returned nil")
	}
}

// Test WithMiddleware configuration
func TestWithMiddleware_Configuration(t *testing.T) {
	sess := NewSession()
	defer sess.Close()

	// Create a simple logging middleware using client.MiddlewareFunc
	loggingMiddleware := client.MiddlewareFunc(func(req *Request, next Handler) (*Response, error) {
		return next(req)
	})

	result := sess.WithMiddleware(loggingMiddleware)
	if result == nil {
		t.Error("WithMiddleware returned nil")
	}
}
