package requests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/sunerpy/requests/internal/client"
	"github.com/sunerpy/requests/internal/models"
)

// Tests for session pool functions
func TestAcquireSession(t *testing.T) {
	session := AcquireSession()
	if session == nil {
		t.Fatal("AcquireSession returned nil")
	}
	defer ReleaseSession(session)
}

func TestAcquireSession_HTTP2Enabled(t *testing.T) {
	original := IsHTTP2Enabled()
	defer SetHTTP2Enabled(original)

	SetHTTP2Enabled(true)
	session := AcquireSession()
	if session == nil {
		t.Fatal("AcquireSession returned nil with HTTP/2 enabled")
	}
	ReleaseSession(session)
}

func TestReleaseSession_Nil(t *testing.T) {
	// Should not panic
	ReleaseSession(nil)
}

func TestReleaseSession_NonDefaultSession(t *testing.T) {
	// Create a mock session that is not *defaultSession
	type mockSession struct {
		client.Session
	}
	// Should not panic
	ReleaseSession(&mockSession{})
}

func TestReleaseSession_HTTP2Transport(t *testing.T) {
	original := IsHTTP2Enabled()
	defer SetHTTP2Enabled(original)

	SetHTTP2Enabled(true)
	session := AcquireSession()
	ReleaseSession(session)
	// Should not panic
}

func TestAcquireReleaseSession_Concurrent(t *testing.T) {
	const numGoroutines = 20
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			session := AcquireSession()
			// Simulate some work
			time.Sleep(time.Millisecond)
			ReleaseSession(session)
		}()
	}
	wg.Wait()
}

// Tests for DoFast
func TestSession_DoFast(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	session := NewSession()
	defer session.Close()

	req, err := NewGet(server.URL).Build()
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := session.(*defaultSession).DoFast(req)
	if err != nil {
		t.Fatalf("DoFast failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestSession_DoFast_WithHeaders(t *testing.T) {
	var receivedHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("X-Custom")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	session := NewSession()
	defer session.Close()

	req, err := NewGet(server.URL).WithHeader("X-Custom", "test-value").Build()
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	_, err = session.(*defaultSession).DoFast(req)
	if err != nil {
		t.Fatalf("DoFast failed: %v", err)
	}
	if receivedHeader != "test-value" {
		t.Errorf("Expected header 'test-value', got '%s'", receivedHeader)
	}
}

func TestSession_DoFast_Error(t *testing.T) {
	session := NewSession()
	defer session.Close()

	req, err := NewGet("http://invalid-host-that-does-not-exist.local").Build()
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	_, err = session.(*defaultSession).DoFast(req)
	if err == nil {
		t.Error("Expected error for invalid host")
	}
}

// Tests for middleware
func TestSession_WithMiddleware(t *testing.T) {
	var middlewareCalled bool
	middleware := client.MiddlewareFunc(func(req *client.Request, next client.Handler) (*models.Response, error) {
		middlewareCalled = true
		return next(req)
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	session := NewSession().WithMiddleware(middleware)
	defer session.Close()

	req, _ := NewGet(server.URL).Build()
	_, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if !middlewareCalled {
		t.Error("Middleware was not called")
	}
}

func TestSession_WithMultipleMiddlewares(t *testing.T) {
	var order []int
	middleware1 := client.MiddlewareFunc(func(req *client.Request, next client.Handler) (*models.Response, error) {
		order = append(order, 1)
		resp, err := next(req)
		order = append(order, 4)
		return resp, err
	})
	middleware2 := client.MiddlewareFunc(func(req *client.Request, next client.Handler) (*models.Response, error) {
		order = append(order, 2)
		resp, err := next(req)
		order = append(order, 3)
		return resp, err
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	session := NewSession().WithMiddleware(middleware1).WithMiddleware(middleware2)
	defer session.Close()

	req, _ := NewGet(server.URL).Build()
	_, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// Middleware should be called in order: 1, 2, (request), 3, 4
	expected := []int{1, 2, 3, 4}
	if len(order) != len(expected) {
		t.Errorf("Expected %d middleware calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if i < len(order) && order[i] != v {
			t.Errorf("Expected order[%d]=%d, got %d", i, v, order[i])
		}
	}
}

// Tests for retry
func TestSession_WithRetry(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	policy := client.RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		RetryIf: func(resp *models.Response, err error) bool {
			return err != nil || (resp != nil && resp.StatusCode >= 500)
		},
	}

	session := NewSession().WithRetry(policy)
	defer session.Close()

	req, _ := NewGet(server.URL).Build()
	resp, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestSession_WithRetry_MaxAttemptsExceeded(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	policy := client.RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		RetryIf: func(resp *models.Response, err error) bool {
			return err != nil || (resp != nil && resp.StatusCode >= 500)
		},
	}

	session := NewSession().WithRetry(policy)
	defer session.Close()

	req, _ := NewGet(server.URL).Build()
	resp, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed with error: %v", err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestSession_WithRetry_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	policy := client.RetryPolicy{
		MaxAttempts:     5,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     time.Second,
		Multiplier:      2.0,
		RetryIf: func(resp *models.Response, err error) bool {
			return true
		},
	}

	session := NewSession().WithRetry(policy)
	defer session.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req, _ := NewGet(server.URL).WithContext(ctx).Build()
	_, err := session.DoWithContext(ctx, req)
	if err == nil {
		t.Error("Expected context canceled error")
	}
}

func TestSession_WithRetry_NoRetryNeeded(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	policy := client.RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		RetryIf: func(resp *models.Response, err error) bool {
			return err != nil || (resp != nil && resp.StatusCode >= 500)
		},
	}

	session := NewSession().WithRetry(policy)
	defer session.Close()

	req, _ := NewGet(server.URL).Build()
	resp, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestSession_WithRetryAndMiddleware(t *testing.T) {
	var middlewareCalls int
	var serverCalls int

	middleware := client.MiddlewareFunc(func(req *client.Request, next client.Handler) (*models.Response, error) {
		middlewareCalls++
		return next(req)
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalls++
		if serverCalls < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	policy := client.RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		RetryIf: func(resp *models.Response, err error) bool {
			return err != nil || (resp != nil && resp.StatusCode >= 500)
		},
	}

	session := NewSession().WithMiddleware(middleware).WithRetry(policy)
	defer session.Close()

	req, _ := NewGet(server.URL).Build()
	resp, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	// Middleware is called once, retry happens inside the final handler
	if middlewareCalls != 1 {
		t.Errorf("Expected 1 middleware call, got %d", middlewareCalls)
	}
	if serverCalls != 2 {
		t.Errorf("Expected 2 server calls, got %d", serverCalls)
	}
}

// Tests for Clone with retry and middleware
func TestSession_Clone_WithRetryAndMiddleware(t *testing.T) {
	middleware := client.MiddlewareFunc(func(req *client.Request, next client.Handler) (*models.Response, error) {
		return next(req)
	})

	policy := client.RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
	}

	original := NewSession().WithMiddleware(middleware).WithRetry(policy)
	defer original.Close()

	clone := original.Clone()
	if clone == nil {
		t.Fatal("Clone returned nil")
	}
}

// Test DoWithContext with timeout
func TestSession_DoWithContext_WithSessionTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	session := NewSession().WithTimeout(50 * time.Millisecond)
	defer session.Close()

	req, _ := NewGet(server.URL).Build()
	_, err := session.DoWithContext(context.Background(), req)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

// Test retry with error (not response)
func TestSession_WithRetry_NetworkError(t *testing.T) {
	policy := client.RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		RetryIf: func(resp *models.Response, err error) bool {
			return err != nil
		},
	}

	session := NewSession().WithRetry(policy)
	defer session.Close()

	// Use a server that we immediately close to cause connection errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	serverURL := server.URL
	server.Close() // Close immediately to cause connection error

	req, _ := NewGet(serverURL).Build()
	_, err := session.Do(req)
	if err == nil {
		t.Error("Expected network error")
	}
}

// Test retry with RetryIf returning false
func TestSession_WithRetry_RetryIfFalse(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest) // 400 error
	}))
	defer server.Close()

	policy := client.RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		RetryIf: func(resp *models.Response, err error) bool {
			// Only retry on 5xx errors
			return err != nil || (resp != nil && resp.StatusCode >= 500)
		},
	}

	session := NewSession().WithRetry(policy)
	defer session.Close()

	req, _ := NewGet(server.URL).Build()
	resp, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
	// Should only attempt once since RetryIf returns false for 400
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}
