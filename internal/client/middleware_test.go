package client

import (
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"github.com/sunerpy/requests/internal/models"
)

// Feature: http-client-refactor
// Property 5: Middleware Execution Order
// For any Client with N registered middleware, request middleware SHALL be executed
// in registration order (1, 2, ..., N), and response middleware SHALL be executed
// in reverse order (N, ..., 2, 1).
// Validates: Requirements 3.1, 3.2, 3.3, 3.4
func TestProperty5_MiddlewareExecutionOrder(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Middleware executes in registration order", prop.ForAll(
		func(n int) bool {
			if n < 1 || n > 10 {
				return true
			}
			var mu sync.Mutex
			var requestOrder []int
			var responseOrder []int
			chain := NewMiddlewareChain()
			// Add n middlewares
			for i := 1; i <= n; i++ {
				idx := i // capture
				chain.UseFunc(func(req *Request, next Handler) (*models.Response, error) {
					mu.Lock()
					requestOrder = append(requestOrder, idx)
					mu.Unlock()
					resp, err := next(req)
					mu.Lock()
					responseOrder = append(responseOrder, idx)
					mu.Unlock()
					return resp, err
				})
			}
			// Create a mock request
			parsedURL, _ := url.Parse("https://example.com")
			req := &Request{
				Method:  MethodGet,
				URL:     parsedURL,
				Headers: make(http.Header),
			}
			// Final handler
			finalHandler := func(r *Request) (*models.Response, error) {
				return CreateMockResponse(200, nil, nil), nil
			}
			_, _ = chain.Execute(req, finalHandler)
			// Verify request order is 1, 2, ..., n
			for i := 0; i < n; i++ {
				if requestOrder[i] != i+1 {
					return false
				}
			}
			// Verify response order is n, n-1, ..., 1
			for i := 0; i < n; i++ {
				if responseOrder[i] != n-i {
					return false
				}
			}
			return true
		},
		gen.IntRange(1, 10),
	))
	properties.TestingRun(t)
}

func TestProperty5_MiddlewareCanShortCircuit(t *testing.T) {
	var executed []string
	chain := NewMiddlewareChain()
	// First middleware - executes normally
	chain.UseFunc(func(req *Request, next Handler) (*models.Response, error) {
		executed = append(executed, "m1-before")
		resp, err := next(req)
		executed = append(executed, "m1-after")
		return resp, err
	})
	// Second middleware - short circuits
	chain.UseFunc(func(req *Request, next Handler) (*models.Response, error) {
		executed = append(executed, "m2-shortcircuit")
		return CreateMockResponse(401, []byte("Unauthorized"), nil), nil
	})
	// Third middleware - should not be called
	chain.UseFunc(func(req *Request, next Handler) (*models.Response, error) {
		executed = append(executed, "m3-should-not-run")
		return next(req)
	})
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	finalHandler := func(r *Request) (*models.Response, error) {
		executed = append(executed, "final-should-not-run")
		return CreateMockResponse(200, nil, nil), nil
	}
	resp, err := chain.Execute(req, finalHandler)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("Expected 401, got %d", resp.StatusCode)
	}
	// Verify execution order
	expected := []string{"m1-before", "m2-shortcircuit", "m1-after"}
	if len(executed) != len(expected) {
		t.Errorf("Expected %d executions, got %d: %v", len(expected), len(executed), executed)
	}
	for i, e := range expected {
		if i >= len(executed) || executed[i] != e {
			t.Errorf("Execution order mismatch at %d: expected %s, got %v", i, e, executed)
		}
	}
}

func TestMiddlewareChain_Empty(t *testing.T) {
	chain := NewMiddlewareChain()
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	called := false
	finalHandler := func(r *Request) (*models.Response, error) {
		called = true
		return CreateMockResponse(200, nil, nil), nil
	}
	resp, err := chain.Execute(req, finalHandler)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !called {
		t.Error("Final handler was not called")
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestMiddlewareChain_Clone(t *testing.T) {
	original := NewMiddlewareChain()
	original.UseFunc(func(req *Request, next Handler) (*models.Response, error) {
		req.Headers.Set("X-Original", "true")
		return next(req)
	})
	clone := original.Clone()
	// Add middleware to clone
	clone.UseFunc(func(req *Request, next Handler) (*models.Response, error) {
		req.Headers.Set("X-Clone", "true")
		return next(req)
	})
	// Original should have 1 middleware
	if original.Len() != 1 {
		t.Errorf("Original should have 1 middleware, got %d", original.Len())
	}
	// Clone should have 2 middlewares
	if clone.Len() != 2 {
		t.Errorf("Clone should have 2 middlewares, got %d", clone.Len())
	}
}

func TestMiddlewareChain_Use(t *testing.T) {
	chain := NewMiddlewareChain()
	if chain.Len() != 0 {
		t.Error("New chain should be empty")
	}
	chain.Use(MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		return next(req)
	}))
	if chain.Len() != 1 {
		t.Error("Chain should have 1 middleware")
	}
	chain.UseFunc(func(req *Request, next Handler) (*models.Response, error) {
		return next(req)
	})
	if chain.Len() != 2 {
		t.Error("Chain should have 2 middlewares")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	var logs []string
	logger := func(format string, args ...any) {
		logs = append(logs, format)
	}
	m := LoggingMiddleware(logger)
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	finalHandler := func(r *Request) (*models.Response, error) {
		return CreateMockResponse(200, nil, nil), nil
	}
	_, _ = m.Process(req, finalHandler)
	if len(logs) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(logs))
	}
}

func TestLoggingMiddleware_WithError(t *testing.T) {
	var logs []string
	logger := func(format string, args ...any) {
		logs = append(logs, format)
	}
	m := LoggingMiddleware(logger)
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	finalHandler := func(r *Request) (*models.Response, error) {
		return nil, &RequestError{Op: "Test", Err: ErrNilRequest}
	}
	_, err := m.Process(req, finalHandler)
	if err == nil {
		t.Error("Expected error")
	}
	// Should still log the request
	if len(logs) < 1 {
		t.Error("Expected at least 1 log entry")
	}
}

func TestHeaderMiddleware(t *testing.T) {
	m := HeaderMiddleware(map[string]string{
		"X-Custom":   "value",
		"X-Existing": "new-value",
	})
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	req.Headers.Set("X-Existing", "old-value")
	finalHandler := func(r *Request) (*models.Response, error) {
		// X-Custom should be added
		if r.Headers.Get("X-Custom") != "value" {
			t.Error("X-Custom header not added")
		}
		// X-Existing should NOT be overwritten
		if r.Headers.Get("X-Existing") != "old-value" {
			t.Error("X-Existing header was overwritten")
		}
		return CreateMockResponse(200, nil, nil), nil
	}
	_, _ = m.Process(req, finalHandler)
}

func TestUserAgentMiddleware(t *testing.T) {
	m := UserAgentMiddleware("MyApp/1.0")
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	finalHandler := func(r *Request) (*models.Response, error) {
		if r.Headers.Get("User-Agent") != "MyApp/1.0" {
			t.Error("User-Agent not set")
		}
		return CreateMockResponse(200, nil, nil), nil
	}
	_, _ = m.Process(req, finalHandler)
}

func TestBearerTokenMiddleware(t *testing.T) {
	m := BearerTokenMiddleware("secret-token")
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	finalHandler := func(r *Request) (*models.Response, error) {
		if r.Headers.Get("Authorization") != "Bearer secret-token" {
			t.Error("Bearer token not set")
		}
		return CreateMockResponse(200, nil, nil), nil
	}
	_, _ = m.Process(req, finalHandler)
}

func TestBasicAuthMiddleware(t *testing.T) {
	m := BasicAuthMiddleware("user", "pass")
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	finalHandler := func(r *Request) (*models.Response, error) {
		auth := r.Headers.Get("Authorization")
		if auth == "" {
			t.Error("Authorization header not set")
		}
		if auth != "Basic dXNlcjpwYXNz" {
			t.Errorf("Wrong basic auth: %s", auth)
		}
		return CreateMockResponse(200, nil, nil), nil
	}
	_, _ = m.Process(req, finalHandler)
}

func TestRecoveryMiddleware(t *testing.T) {
	var recovered any
	m := RecoveryMiddleware(func(req *Request, r any) {
		recovered = r
	})
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	finalHandler := func(r *Request) (*models.Response, error) {
		panic("test panic")
	}
	resp, err := m.Process(req, finalHandler)
	if recovered != "test panic" {
		t.Error("Panic was not recovered")
	}
	if err == nil {
		t.Error("Expected error after panic")
	}
	if resp != nil {
		t.Error("Response should be nil after panic")
	}
}

func TestConditionalMiddleware(t *testing.T) {
	executed := false
	inner := MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		executed = true
		return next(req)
	})
	// Condition that checks for a specific header
	condition := func(req *Request) bool {
		return req.Headers.Get("X-Execute") == "true"
	}
	m := ConditionalMiddleware(condition, inner)
	parsedURL, _ := url.Parse("https://example.com")
	finalHandler := func(r *Request) (*models.Response, error) {
		return CreateMockResponse(200, nil, nil), nil
	}
	// Test without header - should not execute
	req1 := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	_, _ = m.Process(req1, finalHandler)
	if executed {
		t.Error("Middleware should not have executed")
	}
	// Test with header - should execute
	req2 := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	req2.Headers.Set("X-Execute", "true")
	_, _ = m.Process(req2, finalHandler)
	if !executed {
		t.Error("Middleware should have executed")
	}
}

func TestChainMiddleware(t *testing.T) {
	var order []int
	m1 := MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		order = append(order, 1)
		return next(req)
	})
	m2 := MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		order = append(order, 2)
		return next(req)
	})
	combined := ChainMiddleware(m1, m2)
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	finalHandler := func(r *Request) (*models.Response, error) {
		order = append(order, 3)
		return CreateMockResponse(200, nil, nil), nil
	}
	_, _ = combined.Process(req, finalHandler)
	if len(order) != 3 {
		t.Errorf("Expected 3 executions, got %d", len(order))
	}
	for i, expected := range []int{1, 2, 3} {
		if order[i] != expected {
			t.Errorf("Order[%d] = %d, want %d", i, order[i], expected)
		}
	}
}

func TestAuthMiddleware_DoesNotOverwrite(t *testing.T) {
	m := AuthMiddleware("Bearer new-token")
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	req.Headers.Set("Authorization", "Bearer existing-token")
	finalHandler := func(r *Request) (*models.Response, error) {
		// Should keep existing token
		if r.Headers.Get("Authorization") != "Bearer existing-token" {
			t.Error("Existing auth was overwritten")
		}
		return CreateMockResponse(200, nil, nil), nil
	}
	_, _ = m.Process(req, finalHandler)
}

func TestMiddlewareChain_ErrorPropagation(t *testing.T) {
	chain := NewMiddlewareChain()
	chain.UseFunc(func(req *Request, next Handler) (*models.Response, error) {
		resp, err := next(req)
		if err != nil {
			// Wrap the error
			return nil, &RequestError{Op: "Middleware1", Err: err}
		}
		return resp, nil
	})
	chain.UseFunc(func(req *Request, next Handler) (*models.Response, error) {
		return nil, &RequestError{Op: "Middleware2", Err: ErrNilRequest}
	})
	parsedURL, _ := url.Parse("https://example.com")
	req := &Request{
		Method:  MethodGet,
		URL:     parsedURL,
		Headers: make(http.Header),
	}
	finalHandler := func(r *Request) (*models.Response, error) {
		return CreateMockResponse(200, nil, nil), nil
	}
	_, err := chain.Execute(req, finalHandler)
	if err == nil {
		t.Error("Expected error to propagate")
	}
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError, got %T", err)
	}
	if reqErr.Op != "Middleware1" {
		t.Errorf("Expected Middleware1, got %s", reqErr.Op)
	}
}
