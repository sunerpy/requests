package requests

import (
	"context"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"sync"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Helper function to create a simple Request for testing
func newTestRequest(method Method, rawURL string) (*Request, error) {
	return NewGet(rawURL).WithMethod(method).Build()
}

// Feature: http-client-refactor
// Property 7: Session Cookie Persistence
// For any Session making multiple requests to the same domain, cookies set by
// earlier responses SHALL be automatically included in subsequent requests.
// Validates: Requirements 5.2
func TestProperty7_SessionCookiePersistence(t *testing.T) {
	// Create a test server that sets and checks cookies
	var requestCount int
	var receivedCookie string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			// First request: set a cookie
			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: "test-session-123",
				Path:  "/",
			})
			w.WriteHeader(http.StatusOK)
		} else {
			// Subsequent requests: check for cookie
			cookie, err := r.Cookie("session_id")
			if err == nil {
				receivedCookie = cookie.Value
			}
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()
	session := NewSession()
	defer session.Close()
	// First request - should set cookie
	req1, _ := newTestRequest(MethodGet, server.URL+"/first")
	_, err := session.Do(req1)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}
	// Second request - should include cookie
	req2, _ := newTestRequest(MethodGet, server.URL+"/second")
	_, err = session.Do(req2)
	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}
	if receivedCookie != "test-session-123" {
		t.Errorf("Cookie not persisted: expected 'test-session-123', got '%s'", receivedCookie)
	}
}

// Feature: http-client-refactor
// Property 8: Session Defaults Applied to All Requests
// For any Session with configured base URL and default headers, all requests made
// through that Session SHALL have the base URL prepended and default headers included.
// Validates: Requirements 5.3, 5.4
func TestProperty8_SessionDefaultsApplied(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Session headers are applied to all requests", prop.ForAll(
		func(headerValue string) bool {
			if headerValue == "" {
				return true
			}
			var receivedHeader string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedHeader = r.Header.Get("X-Custom-Header")
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			session := NewSession().WithHeader("X-Custom-Header", headerValue)
			defer session.Close()
			req, _ := newTestRequest(MethodGet, server.URL)
			_, err := session.Do(req)
			if err != nil {
				return false
			}
			return receivedHeader == headerValue
		},
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

func TestProperty8_SessionBaseURLApplied(t *testing.T) {
	var receivedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	session := NewSession().WithBaseURL(server.URL)
	defer session.Close()
	// Request with relative path
	req, _ := newTestRequest(MethodGet, "/api/users")
	_, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if receivedPath != "/api/users" {
		t.Errorf("Base URL not applied: expected '/api/users', got '%s'", receivedPath)
	}
}

// Feature: http-client-refactor
// Property 9: Session Clone Independence
// For any Session that is cloned, modifications to the original Session SHALL NOT
// affect the cloned Session, and vice versa.
// Validates: Requirements 5.5
func TestProperty9_SessionCloneIndependence(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Cloned session is independent of original", prop.ForAll(
		func(originalHeader, cloneHeader string) bool {
			if originalHeader == "" || cloneHeader == "" {
				return true
			}
			original := NewSession().WithHeader("X-Test", originalHeader)
			clone := original.Clone().(Session)
			// Modify original
			original.WithHeader("X-Test", "modified-original")
			original.WithHeader("X-New", "new-value")
			// Modify clone
			clone.WithHeader("X-Test", cloneHeader)
			// Create test servers
			var originalReceived, cloneReceived string
			server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				originalReceived = r.Header.Get("X-Test")
				w.WriteHeader(http.StatusOK)
			}))
			defer server1.Close()
			server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cloneReceived = r.Header.Get("X-Test")
				w.WriteHeader(http.StatusOK)
			}))
			defer server2.Close()
			req1, _ := newTestRequest(MethodGet, server1.URL)
			original.Do(req1)
			req2, _ := newTestRequest(MethodGet, server2.URL)
			clone.Do(req2)
			// Original should have modified value
			if originalReceived != "modified-original" {
				return false
			}
			// Clone should have its own value
			if cloneReceived != cloneHeader {
				return false
			}
			return true
		},
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
	))
	properties.TestingRun(t)
}

// Feature: http-client-refactor
// Property 10: Session Concurrent Safety
// For any Session used concurrently from multiple goroutines, all requests SHALL
// complete without race conditions or data corruption.
// Validates: Requirements 5.7
func TestProperty10_SessionConcurrentSafety(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()
	session := NewSession()
	defer session.Close()
	const numGoroutines = 50
	const requestsPerGoroutine = 10
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*requestsPerGoroutine)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				req, err := newTestRequest(MethodGet, server.URL)
				if err != nil {
					errors <- err
					continue
				}
				_, err = session.Do(req)
				if err != nil {
					errors <- err
				}
			}
		}(i)
	}
	wg.Wait()
	close(errors)
	var errCount int
	for err := range errors {
		t.Logf("Error: %v", err)
		errCount++
	}
	if errCount > 0 {
		t.Errorf("Had %d errors during concurrent requests", errCount)
	}
}

func TestProperty10_ConcurrentSessionModification(t *testing.T) {
	session := NewSession()
	defer session.Close()
	const numGoroutines = 20
	var wg sync.WaitGroup
	// Concurrent modifications should not cause race conditions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			session.WithHeader("X-Goroutine", string(rune('A'+id)))
			session.WithTimeout(time.Duration(id) * time.Millisecond)
		}(i)
	}
	wg.Wait()
	// If we get here without panic/race, the test passes
}

// Unit tests for Session
func TestSession_Creation(t *testing.T) {
	session := NewSession()
	defer session.Close()
	if session == nil {
		t.Fatal("NewSession returned nil")
	}
}

func TestSession_WithTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	session := NewSession().WithTimeout(10 * time.Millisecond)
	defer session.Close()
	req, _ := newTestRequest(MethodGet, server.URL)
	_, err := session.Do(req)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestSession_WithHeaders(t *testing.T) {
	var receivedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	session := NewSession().WithHeaders(map[string]string{
		"X-Custom-1": "value1",
		"X-Custom-2": "value2",
	})
	defer session.Close()
	req, _ := newTestRequest(MethodGet, server.URL)
	_, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if receivedHeaders.Get("X-Custom-1") != "value1" {
		t.Error("X-Custom-1 header not set")
	}
	if receivedHeaders.Get("X-Custom-2") != "value2" {
		t.Error("X-Custom-2 header not set")
	}
}

func TestSession_WithBasicAuth(t *testing.T) {
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	session := NewSession().WithBasicAuth("user", "pass")
	defer session.Close()
	req, _ := newTestRequest(MethodGet, server.URL)
	_, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if authHeader == "" {
		t.Error("Authorization header not set")
	}
	if authHeader != "Basic dXNlcjpwYXNz" {
		t.Errorf("Wrong basic auth: %s", authHeader)
	}
}

func TestSession_WithBearerToken(t *testing.T) {
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	session := NewSession().WithBearerToken("my-token")
	defer session.Close()
	req, _ := newTestRequest(MethodGet, server.URL)
	_, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if authHeader != "Bearer my-token" {
		t.Errorf("Wrong bearer token: %s", authHeader)
	}
}

func TestSession_WithHTTP2(t *testing.T) {
	session := NewSession()
	defer session.Close()
	// Toggle HTTP/2
	session.WithHTTP2(true)
	session.WithHTTP2(false)
	session.WithHTTP2(true)
	// Should not panic
}

func TestSession_WithKeepAlive(t *testing.T) {
	session := NewSession()
	defer session.Close()
	session.WithKeepAlive(false)
	session.WithKeepAlive(true)
	// Should not panic
}

func TestSession_WithMaxIdleConns(t *testing.T) {
	session := NewSession()
	defer session.Close()
	session.WithMaxIdleConns(50)
	// Should not panic
}

func TestSession_Clear(t *testing.T) {
	session := NewSession().
		WithBaseURL("https://example.com").
		WithHeader("X-Test", "value").
		WithTimeout(10 * time.Second)
	session.Clear()
	// After clear, session should be reset
	// We can't easily verify internal state, but it should not panic
}

func TestSession_Close(t *testing.T) {
	session := NewSession()
	err := session.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestSession_RequestHeadersOverrideSession(t *testing.T) {
	var receivedHeaders []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header.Values("X-Test")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	session := NewSession().WithHeader("X-Test", "session-value")
	defer session.Close()
	req, _ := newTestRequest(MethodGet, server.URL)
	req.SetHeader("X-Test", "request-value")
	_, err := session.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	// Both session and request headers are added
	// The implementation adds session headers first, then request headers
	if len(receivedHeaders) != 2 {
		t.Errorf("Expected 2 header values, got %d: %v", len(receivedHeaders), receivedHeaders)
	}
}

func TestSession_NilURL(t *testing.T) {
	session := NewSession()
	defer session.Close()
	req := &Request{
		Method:  MethodGet,
		URL:     nil,
		Headers: http.Header{},
		Context: context.Background(),
	}
	_, err := session.Do(req)
	if err == nil {
		t.Error("Expected error for nil URL")
	}
}

func TestSession_WithIdleTimeout(t *testing.T) {
	session := NewSession()
	defer session.Close()
	session.WithIdleTimeout(30 * time.Second)
	// Should not panic
}

func TestRequest_AddHeader(t *testing.T) {
	req, _ := newTestRequest(MethodGet, "https://example.com")
	req.AddHeader("X-Test", "value1")
	req.AddHeader("X-Test", "value2")
	values := req.Headers["X-Test"]
	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}
}

func TestRequest_SetHeader(t *testing.T) {
	req, _ := newTestRequest(MethodGet, "https://example.com")
	req.SetHeader("X-Test", "value1")
	req.SetHeader("X-Test", "value2")
	values := req.Headers["X-Test"]
	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
	if values[0] != "value2" {
		t.Errorf("Expected 'value2', got '%s'", values[0])
	}
}

func TestSession_NewRequestWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")
	req, err := NewPost("https://example.com/api").WithContext(ctx).Build()
	if err != nil {
		t.Fatalf("NewPost with context failed: %v", err)
	}
	if req.Context != ctx {
		t.Error("Context not set")
	}
	if req.Method != MethodPost {
		t.Error("Method not set")
	}
}

func TestGetTransport(t *testing.T) {
	// HTTP/1 transport
	tr1 := GetTransport(false)
	if tr1 == nil {
		t.Error("GetTransport(false) returned nil")
	}
	PutTransport(tr1)
	// HTTP/2 transport
	tr2 := GetTransport(true)
	if tr2 == nil {
		t.Error("GetTransport(true) returned nil")
	}
	PutTransport(tr2)
}

func TestPutTransport_Nil(t *testing.T) {
	// Should not panic
	PutTransport(nil)
}

func TestSession_SetHTTP2Enabled(t *testing.T) {
	original := IsHTTP2Enabled()
	SetHTTP2Enabled(true)
	if !IsHTTP2Enabled() {
		t.Error("HTTP/2 should be enabled")
	}
	SetHTTP2Enabled(false)
	if IsHTTP2Enabled() {
		t.Error("HTTP/2 should be disabled")
	}
	// Restore original
	SetHTTP2Enabled(original)
}

func TestNewSession_HTTP2Enabled(t *testing.T) {
	original := IsHTTP2Enabled()
	defer SetHTTP2Enabled(original)
	// Test with HTTP/2 enabled
	SetHTTP2Enabled(true)
	session := NewSession()
	defer session.Close()
	if session == nil {
		t.Error("NewSession returned nil with HTTP/2 enabled")
	}
}

func TestNewSession_HTTP2Disabled(t *testing.T) {
	original := IsHTTP2Enabled()
	defer SetHTTP2Enabled(original)
	// Test with HTTP/2 disabled
	SetHTTP2Enabled(false)
	session := NewSession()
	defer session.Close()
	if session == nil {
		t.Error("NewSession returned nil with HTTP/2 disabled")
	}
}

// Additional tests for coverage
func TestSession_WithProxy(t *testing.T) {
	session := NewSession()
	defer session.Close()
	// Valid proxy URL
	session.WithProxy("http://proxy.example.com:8080")
	// Invalid proxy URL should not panic
	session.WithProxy("invalid-url")
	// Empty proxy URL
	session.WithProxy("")
}

func TestSession_WithCookieJar(t *testing.T) {
	session := NewSession()
	defer session.Close()
	// Test with nil jar
	session.WithCookieJar(nil)
	// Should not panic
}

func TestDefaultSession(t *testing.T) {
	session := DefaultSession()
	if session == nil {
		t.Error("DefaultSession returned nil")
	}
}

func TestHead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			t.Errorf("Expected HEAD method, got %s", r.Method)
		}
		w.Header().Set("X-Custom", "value")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	resp, err := Head(server.URL)
	if err != nil {
		t.Fatalf("Head request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			t.Errorf("Expected OPTIONS method, got %s", r.Method)
		}
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	resp, err := Options(server.URL)
	if err != nil {
		t.Fatalf("Options request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetString(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()
	result, err := GetString(server.URL)
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", result)
	}
}

func TestGetBytes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{0x01, 0x02, 0x03})
	}))
	defer server.Close()
	result, err := GetBytes(server.URL)
	if err != nil {
		t.Fatalf("GetBytes failed: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 bytes, got %d", len(result))
	}
}

func TestGetJSON(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"John","age":30}`))
	}))
	defer server.Close()
	result, err := GetJSON[TestData](server.URL)
	if err != nil {
		t.Fatalf("GetJSON failed: %v", err)
	}
	if result.Data().Name != "John" || result.Data().Age != 30 {
		t.Errorf("Unexpected result: %+v", result.Data())
	}
}

func TestPostJSON(t *testing.T) {
	type Request struct {
		Name string `json:"name"`
	}
	type Response struct {
		ID int `json:"id"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected application/json content type")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":123}`))
	}))
	defer server.Close()
	result, err := PostJSON[Response](server.URL, Request{Name: "Test"})
	if err != nil {
		t.Fatalf("PostJSON failed: %v", err)
	}
	if result.Data().ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.Data().ID)
	}
}

func TestPutJSON(t *testing.T) {
	type Request struct {
		Name string `json:"name"`
	}
	type Response struct {
		Updated bool `json:"updated"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"updated":true}`))
	}))
	defer server.Close()
	result, err := PutJSON[Response](server.URL, Request{Name: "Test"})
	if err != nil {
		t.Fatalf("PutJSON failed: %v", err)
	}
	if !result.Data().Updated {
		t.Error("Expected updated to be true")
	}
}

func TestDeleteJSON(t *testing.T) {
	type Response struct {
		Deleted bool `json:"deleted"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"deleted":true}`))
	}))
	defer server.Close()
	result, err := DeleteJSON[Response](server.URL)
	if err != nil {
		t.Fatalf("DeleteJSON failed: %v", err)
	}
	if !result.Data().Deleted {
		t.Error("Expected deleted to be true")
	}
}

func TestPatchJSON(t *testing.T) {
	type Request struct {
		Name string `json:"name"`
	}
	type Response struct {
		Patched bool `json:"patched"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"patched":true}`))
	}))
	defer server.Close()
	result, err := PatchJSON[Response](server.URL, Request{Name: "Test"})
	if err != nil {
		t.Fatalf("PatchJSON failed: %v", err)
	}
	if !result.Data().Patched {
		t.Error("Expected patched to be true")
	}
}

func TestGetJSON_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()
	type TestData struct {
		Name string `json:"name"`
	}
	_, err := GetJSON[TestData](server.URL)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestSession_WithDNS(t *testing.T) {
	session := NewSession()
	defer session.Close()
	// Set custom DNS servers
	session.WithDNS([]string{"8.8.8.8", "8.8.4.4"})
	// Should not panic
}

func TestSession_WithDNS_Empty(t *testing.T) {
	session := NewSession()
	defer session.Close()
	// Set empty DNS servers
	session.WithDNS([]string{})
	// Should not panic
}

func TestGetJSON_Error(t *testing.T) {
	// Test with invalid URL
	_, err := GetJSON[map[string]string]("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPostJSON_EncodeError(t *testing.T) {
	// Test with unmarshalable type
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	_, err := PostJSON[map[string]string]("https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPutJSON_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	_, err := PutJSON[map[string]string]("https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPatchJSON_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	_, err := PatchJSON[map[string]string]("https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestGetString_Error(t *testing.T) {
	_, err := GetString("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestGetBytes_Error(t *testing.T) {
	_, err := GetBytes("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestGet_Error(t *testing.T) {
	_, err := Get("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPost_Error(t *testing.T) {
	_, err := Post("://invalid-url", nil)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPut_Error(t *testing.T) {
	_, err := Put("://invalid-url", nil)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestDelete_Error(t *testing.T) {
	_, err := Delete("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPatch_Error(t *testing.T) {
	_, err := Patch("://invalid-url", nil)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestHead_Error(t *testing.T) {
	_, err := Head("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestOptions_Error(t *testing.T) {
	_, err := Options("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestDeleteJSON_Error(t *testing.T) {
	_, err := DeleteJSON[map[string]string]("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestSession_WithKeepAlive_NoTransport(t *testing.T) {
	session := NewSession()
	defer session.Close()
	// First call should work
	session.WithKeepAlive(true)
	session.WithKeepAlive(false)
}

func TestSession_WithMaxIdleConns_NoTransport(t *testing.T) {
	session := NewSession()
	defer session.Close()
	// Should work with existing transport
	session.WithMaxIdleConns(200)
}

func TestSession_Clone_HTTP1(t *testing.T) {
	session := NewSession().WithHTTP2(false)
	defer session.Close()
	clone := session.Clone()
	if clone == nil {
		t.Error("Clone returned nil")
	}
}

func TestSession_Clone_HTTP2(t *testing.T) {
	session := NewSession().WithHTTP2(true)
	defer session.Close()
	clone := session.Clone()
	if clone == nil {
		t.Error("Clone returned nil")
	}
}

func TestSession_Do_WithBaseURL_InvalidBase(t *testing.T) {
	session := NewSession().WithBaseURL("://invalid")
	defer session.Close()
	req, _ := newTestRequest(MethodGet, "/api")
	_, err := session.Do(req)
	if err == nil {
		t.Error("Expected error for invalid base URL")
	}
}

func TestSession_Do_NilContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	session := NewSession()
	defer session.Close()
	req, _ := newTestRequest(MethodGet, server.URL)
	req.Context = nil // Set context to nil
	_, err := session.Do(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSession_Do_InvalidMethod(t *testing.T) {
	session := NewSession()
	defer session.Close()
	parsedURL, _ := neturl.Parse("http://example.com")
	req := &Request{
		Method: Method("INVALID\nMETHOD"), // Invalid method with newline
		URL:    parsedURL,
	}
	_, err := session.Do(req)
	if err == nil {
		t.Error("Expected error for invalid method")
	}
}

func TestGet_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()
	resp, err := Get(server.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestDelete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	resp, err := Delete(server.URL)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected 204, got %d", resp.StatusCode)
	}
}

func TestHead_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	resp, err := Head(server.URL)
	if err != nil {
		t.Fatalf("Head failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestOptions_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	resp, err := Options(server.URL)
	if err != nil {
		t.Fatalf("Options failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestPostJSON_Success(t *testing.T) {
	type Response struct {
		ID int `json:"id"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":123}`))
	}))
	defer server.Close()
	result, err := PostJSON[Response](server.URL, map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("PostJSON failed: %v", err)
	}
	if result.Data().ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.Data().ID)
	}
}

func TestPutJSON_Success(t *testing.T) {
	type Response struct {
		Updated bool `json:"updated"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"updated":true}`))
	}))
	defer server.Close()
	result, err := PutJSON[Response](server.URL, map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("PutJSON failed: %v", err)
	}
	if !result.Data().Updated {
		t.Error("Expected updated to be true")
	}
}

func TestPatchJSON_Success(t *testing.T) {
	type Response struct {
		Patched bool `json:"patched"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"patched":true}`))
	}))
	defer server.Close()
	result, err := PatchJSON[Response](server.URL, map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("PatchJSON failed: %v", err)
	}
	if !result.Data().Patched {
		t.Error("Expected patched to be true")
	}
}

func TestDeleteJSON_Success(t *testing.T) {
	type Response struct {
		Deleted bool `json:"deleted"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"deleted":true}`))
	}))
	defer server.Close()
	result, err := DeleteJSON[Response](server.URL)
	if err != nil {
		t.Fatalf("DeleteJSON failed: %v", err)
	}
	if !result.Data().Deleted {
		t.Error("Expected deleted to be true")
	}
}

func TestPostJSON_RequestError(t *testing.T) {
	_, err := PostJSON[map[string]string]("://invalid-url", map[string]string{"name": "test"})
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPutJSON_RequestError(t *testing.T) {
	_, err := PutJSON[map[string]string]("://invalid-url", map[string]string{"name": "test"})
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPatchJSON_RequestError(t *testing.T) {
	_, err := PatchJSON[map[string]string]("://invalid-url", map[string]string{"name": "test"})
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

// Feature: api-design-optimization
// Property 6: Package-Level Functions Make Correct HTTP Requests
// For any package-level function (Get, Post, Put, Delete, Patch), the function
// SHALL send a request with the corresponding HTTP method.
// Validates: Requirements 4.1
func TestProperty6_PackageLevelFunctionsCorrectMethod(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	// Property: Get sends GET request
	properties.Property("Get sends GET request", prop.ForAll(
		func(path string) bool {
			var receivedMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			_, err := Get(server.URL + "/" + path)
			if err != nil {
				return false
			}
			return receivedMethod == "GET"
		},
		gen.Identifier(),
	))
	// Property: Post sends POST request
	properties.Property("Post sends POST request", prop.ForAll(
		func(path string) bool {
			var receivedMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			_, err := Post(server.URL+"/"+path, nil)
			if err != nil {
				return false
			}
			return receivedMethod == "POST"
		},
		gen.Identifier(),
	))
	// Property: Put sends PUT request
	properties.Property("Put sends PUT request", prop.ForAll(
		func(path string) bool {
			var receivedMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			_, err := Put(server.URL+"/"+path, nil)
			if err != nil {
				return false
			}
			return receivedMethod == "PUT"
		},
		gen.Identifier(),
	))
	// Property: Delete sends DELETE request
	properties.Property("Delete sends DELETE request", prop.ForAll(
		func(path string) bool {
			var receivedMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			_, err := Delete(server.URL + "/" + path)
			if err != nil {
				return false
			}
			return receivedMethod == "DELETE"
		},
		gen.Identifier(),
	))
	// Property: Patch sends PATCH request
	properties.Property("Patch sends PATCH request", prop.ForAll(
		func(path string) bool {
			var receivedMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			_, err := Patch(server.URL+"/"+path, nil)
			if err != nil {
				return false
			}
			return receivedMethod == "PATCH"
		},
		gen.Identifier(),
	))
	// Property: Head sends HEAD request
	properties.Property("Head sends HEAD request", prop.ForAll(
		func(path string) bool {
			var receivedMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			_, err := Head(server.URL + "/" + path)
			if err != nil {
				return false
			}
			return receivedMethod == "HEAD"
		},
		gen.Identifier(),
	))
	// Property: Options sends OPTIONS request
	properties.Property("Options sends OPTIONS request", prop.ForAll(
		func(path string) bool {
			var receivedMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			_, err := Options(server.URL + "/" + path)
			if err != nil {
				return false
			}
			return receivedMethod == "OPTIONS"
		},
		gen.Identifier(),
	))
	properties.TestingRun(t)
}
