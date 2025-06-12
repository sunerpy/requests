package client

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: http-client-refactor
// Property 6: Request Options Merge with Session Defaults
// For any Session with default headers/timeout and a request with additional options,
// the final request SHALL contain both session defaults and request-specific options,
// with request options taking precedence on conflicts.
// Validates: Requirements 4.1, 4.3
func TestProperty6_OptionsMergeWithDefaults(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Request options merge with session defaults", prop.ForAll(
		func(sessionTimeout, requestTimeout int) bool {
			// Create session config with defaults
			sessionConfig := NewRequestConfig()
			sessionConfig.Timeout = time.Duration(sessionTimeout) * time.Millisecond
			sessionConfig.Headers.Set("X-Session-Header", "session-value")
			sessionConfig.Headers.Set("X-Shared-Header", "session-shared")
			// Create request config with overrides
			requestConfig := NewRequestConfig()
			if requestTimeout > 0 {
				requestConfig.Timeout = time.Duration(requestTimeout) * time.Millisecond
			}
			requestConfig.Headers.Set("X-Request-Header", "request-value")
			requestConfig.Headers.Set("X-Shared-Header", "request-shared")
			// Merge request into session (session is base, request overrides)
			merged := sessionConfig.Clone()
			merged.Merge(requestConfig)
			// Session-only header should be preserved
			if merged.Headers.Get("X-Session-Header") != "session-value" {
				return false
			}
			// Request-only header should be present
			if merged.Headers.Get("X-Request-Header") != "request-value" {
				return false
			}
			// Shared header should have request value (request takes precedence)
			if merged.Headers.Get("X-Shared-Header") != "request-shared" {
				return false
			}
			// Timeout: request takes precedence if set
			if requestTimeout > 0 {
				if merged.Timeout != time.Duration(requestTimeout)*time.Millisecond {
					return false
				}
			} else {
				if merged.Timeout != time.Duration(sessionTimeout)*time.Millisecond {
					return false
				}
			}
			return true
		},
		gen.IntRange(100, 5000),
		gen.IntRange(0, 5000),
	))
	properties.TestingRun(t)
}

func TestProperty6_OptionsPreserveAllValues(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("All option values are preserved after Apply", prop.ForAll(
		func(timeout int, headerKey, headerValue string) bool {
			if headerKey == "" || headerValue == "" {
				return true
			}
			config := NewRequestConfig()
			config.Apply(
				WithTimeout(time.Duration(timeout)*time.Millisecond),
				WithHeader(headerKey, headerValue),
			)
			if config.Timeout != time.Duration(timeout)*time.Millisecond {
				return false
			}
			if config.Headers.Get(headerKey) != headerValue {
				return false
			}
			return true
		},
		gen.IntRange(100, 10000),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

func TestProperty6_CloneIndependence(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Cloned config is independent of original", prop.ForAll(
		func(timeout int, headerValue string) bool {
			original := NewRequestConfig()
			original.Timeout = time.Duration(timeout) * time.Millisecond
			original.Headers.Set("X-Test", headerValue)
			clone := original.Clone()
			// Modify original
			original.Timeout = time.Duration(timeout+1000) * time.Millisecond
			original.Headers.Set("X-Test", "modified")
			original.Headers.Set("X-New", "new-value")
			// Clone should not be affected
			if clone.Timeout != time.Duration(timeout)*time.Millisecond {
				return false
			}
			if clone.Headers.Get("X-Test") != headerValue {
				return false
			}
			if clone.Headers.Get("X-New") != "" {
				return false
			}
			return true
		},
		gen.IntRange(100, 5000),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

// Unit tests for RequestConfig
func TestNewRequestConfig(t *testing.T) {
	config := NewRequestConfig()
	if config.Headers == nil {
		t.Error("Headers should be initialized")
	}
	if config.Query == nil {
		t.Error("Query should be initialized")
	}
	if config.Timeout != 0 {
		t.Error("Timeout should be zero")
	}
}

func TestRequestConfig_Apply(t *testing.T) {
	config := NewRequestConfig()
	config.Apply(
		WithTimeout(5*time.Second),
		WithHeader("X-Custom", "value"),
		WithQuery("page", "1"),
		WithBasicAuth("user", "pass"),
	)
	if config.Timeout != 5*time.Second {
		t.Error("Timeout not applied")
	}
	if config.Headers.Get("X-Custom") != "value" {
		t.Error("Header not applied")
	}
	if config.Query.Get("page") != "1" {
		t.Error("Query not applied")
	}
	if config.BasicAuth == nil || config.BasicAuth.Username != "user" {
		t.Error("BasicAuth not applied")
	}
}

func TestRequestConfig_Clone(t *testing.T) {
	original := NewRequestConfig()
	original.Timeout = 10 * time.Second
	original.Headers.Set("X-Test", "value")
	original.Query.Set("key", "value")
	original.BasicAuth = &BasicAuth{Username: "user", Password: "pass"}
	original.BearerToken = "token"
	original.Retry = &RetryPolicy{MaxAttempts: 3}
	clone := original.Clone()
	// Verify all values are copied
	if clone.Timeout != original.Timeout {
		t.Error("Timeout not cloned")
	}
	if clone.Headers.Get("X-Test") != "value" {
		t.Error("Headers not cloned")
	}
	if clone.Query.Get("key") != "value" {
		t.Error("Query not cloned")
	}
	if clone.BasicAuth == nil || clone.BasicAuth.Username != "user" {
		t.Error("BasicAuth not cloned")
	}
	if clone.BearerToken != "token" {
		t.Error("BearerToken not cloned")
	}
	if clone.Retry == nil || clone.Retry.MaxAttempts != 3 {
		t.Error("Retry not cloned")
	}
	// Verify independence
	original.Headers.Set("X-Test", "modified")
	if clone.Headers.Get("X-Test") == "modified" {
		t.Error("Clone headers affected by original")
	}
	original.Query.Set("key", "modified")
	if clone.Query.Get("key") == "modified" {
		t.Error("Clone query affected by original")
	}
}

func TestRequestConfig_Merge(t *testing.T) {
	base := NewRequestConfig()
	base.Timeout = 5 * time.Second
	base.Headers.Set("X-Base", "base-value")
	base.Headers.Set("X-Shared", "base-shared")
	override := NewRequestConfig()
	override.Timeout = 10 * time.Second
	override.Headers.Set("X-Override", "override-value")
	override.Headers.Set("X-Shared", "override-shared")
	base.Merge(override)
	// Override values should take precedence
	if base.Timeout != 10*time.Second {
		t.Error("Timeout not overridden")
	}
	// Base-only values should be preserved
	if base.Headers.Get("X-Base") != "base-value" {
		t.Error("Base header lost")
	}
	// Override-only values should be added
	if base.Headers.Get("X-Override") != "override-value" {
		t.Error("Override header not added")
	}
	// Shared values should use override
	if base.Headers.Get("X-Shared") != "override-shared" {
		t.Error("Shared header not overridden")
	}
}

func TestRequestConfig_MergeNil(t *testing.T) {
	config := NewRequestConfig()
	config.Timeout = 5 * time.Second
	// Should not panic
	config.Merge(nil)
	if config.Timeout != 5*time.Second {
		t.Error("Config modified by nil merge")
	}
}

func TestRequestConfig_ApplyToRequest(t *testing.T) {
	config := NewRequestConfig()
	config.Headers.Set("X-Custom", "value")
	config.Query.Set("page", "1")
	config.BasicAuth = &BasicAuth{Username: "user", Password: "pass"}
	req, _ := http.NewRequest("GET", "https://example.com/api", nil)
	config.ApplyToRequest(req)
	if req.Header.Get("X-Custom") != "value" {
		t.Error("Header not applied to request")
	}
	if req.URL.Query().Get("page") != "1" {
		t.Error("Query not applied to request")
	}
	auth := req.Header.Get("Authorization")
	if auth == "" {
		t.Error("Authorization header not set")
	}
}

func TestRequestConfig_ApplyToRequest_BearerToken(t *testing.T) {
	config := NewRequestConfig()
	config.BearerToken = "my-token"
	req, _ := http.NewRequest("GET", "https://example.com/api", nil)
	config.ApplyToRequest(req)
	if req.Header.Get("Authorization") != "Bearer my-token" {
		t.Error("Bearer token not applied")
	}
}

// Unit tests for individual options
func TestWithTimeout(t *testing.T) {
	config := NewRequestConfig()
	WithTimeout(30 * time.Second)(config)
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected 30s, got %v", config.Timeout)
	}
}

func TestWithHeader(t *testing.T) {
	config := NewRequestConfig()
	WithHeader("X-Test", "value")(config)
	if config.Headers.Get("X-Test") != "value" {
		t.Error("Header not set")
	}
}

func TestWithHeaders(t *testing.T) {
	config := NewRequestConfig()
	WithHeaders(map[string]string{
		"X-One": "1",
		"X-Two": "2",
	})(config)
	if config.Headers.Get("X-One") != "1" || config.Headers.Get("X-Two") != "2" {
		t.Error("Headers not set")
	}
}

func TestWithQuery(t *testing.T) {
	config := NewRequestConfig()
	WithQuery("key", "value")(config)
	if config.Query.Get("key") != "value" {
		t.Error("Query not set")
	}
}

func TestWithQueryParams(t *testing.T) {
	config := NewRequestConfig()
	WithQueryParams(map[string]string{
		"page":  "1",
		"limit": "10",
	})(config)
	if config.Query.Get("page") != "1" || config.Query.Get("limit") != "10" {
		t.Error("Query params not set")
	}
}

func TestWithBasicAuth(t *testing.T) {
	config := NewRequestConfig()
	WithBasicAuth("user", "pass")(config)
	if config.BasicAuth == nil {
		t.Fatal("BasicAuth not set")
	}
	if config.BasicAuth.Username != "user" || config.BasicAuth.Password != "pass" {
		t.Error("BasicAuth values incorrect")
	}
}

func TestWithBearerToken(t *testing.T) {
	config := NewRequestConfig()
	WithBearerToken("token123")(config)
	if config.BearerToken != "token123" {
		t.Error("BearerToken not set")
	}
}

func TestWithRetry(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 5}
	config := NewRequestConfig()
	WithRetry(policy)(config)
	if config.Retry == nil || config.Retry.MaxAttempts != 5 {
		t.Error("Retry policy not set")
	}
}

func TestWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")
	config := NewRequestConfig()
	WithContext(ctx)(config)
	if config.Context != ctx {
		t.Error("Context not set")
	}
}

func TestWithContentType(t *testing.T) {
	config := NewRequestConfig()
	WithContentType("application/json")(config)
	if config.Headers.Get("Content-Type") != "application/json" {
		t.Error("Content-Type not set")
	}
}

func TestWithAccept(t *testing.T) {
	config := NewRequestConfig()
	WithAccept("application/json")(config)
	if config.Headers.Get("Accept") != "application/json" {
		t.Error("Accept not set")
	}
}

func TestWithUserAgent(t *testing.T) {
	config := NewRequestConfig()
	WithUserAgent("MyApp/1.0")(config)
	if config.Headers.Get("User-Agent") != "MyApp/1.0" {
		t.Error("User-Agent not set")
	}
}

func TestOptionsChaining(t *testing.T) {
	config := NewRequestConfig()
	config.Apply(
		WithTimeout(10*time.Second),
		WithHeader("X-Custom", "value"),
		WithQuery("page", "1"),
		WithBearerToken("token"),
		WithContentType("application/json"),
		WithAccept("application/json"),
		WithUserAgent("TestApp/1.0"),
	)
	if config.Timeout != 10*time.Second {
		t.Error("Timeout not set")
	}
	if config.Headers.Get("X-Custom") != "value" {
		t.Error("Custom header not set")
	}
	if config.Query.Get("page") != "1" {
		t.Error("Query not set")
	}
	if config.BearerToken != "token" {
		t.Error("Bearer token not set")
	}
	if config.Headers.Get("Content-Type") != "application/json" {
		t.Error("Content-Type not set")
	}
	if config.Headers.Get("Accept") != "application/json" {
		t.Error("Accept not set")
	}
	if config.Headers.Get("User-Agent") != "TestApp/1.0" {
		t.Error("User-Agent not set")
	}
}

func TestRequestConfig_MergePreservesNonConflicting(t *testing.T) {
	base := NewRequestConfig()
	base.Headers.Set("X-Base-Only", "base")
	base.Query.Set("base-param", "value")
	base.BasicAuth = &BasicAuth{Username: "base-user", Password: "base-pass"}
	override := NewRequestConfig()
	override.Headers.Set("X-Override-Only", "override")
	override.Query.Set("override-param", "value")
	// No BasicAuth in override
	base.Merge(override)
	// Base-only values preserved
	if base.Headers.Get("X-Base-Only") != "base" {
		t.Error("Base header lost")
	}
	if base.Query.Get("base-param") != "value" {
		t.Error("Base query lost")
	}
	if base.BasicAuth == nil || base.BasicAuth.Username != "base-user" {
		t.Error("Base BasicAuth lost")
	}
	// Override values added
	if base.Headers.Get("X-Override-Only") != "override" {
		t.Error("Override header not added")
	}
	if base.Query.Get("override-param") != "value" {
		t.Error("Override query not added")
	}
}

func TestRequestConfig_ApplyToRequest_MergesQuery(t *testing.T) {
	config := NewRequestConfig()
	config.Query.Set("config-param", "config-value")
	req, _ := http.NewRequest("GET", "https://example.com/api?existing=value", nil)
	config.ApplyToRequest(req)
	query := req.URL.Query()
	if query.Get("existing") != "value" {
		t.Error("Existing query param lost")
	}
	if query.Get("config-param") != "config-value" {
		t.Error("Config query param not added")
	}
}

func TestRequestConfig_CloneWithNilValues(t *testing.T) {
	original := NewRequestConfig()
	// Leave BasicAuth and Retry as nil
	clone := original.Clone()
	if clone.BasicAuth != nil {
		t.Error("Clone should have nil BasicAuth")
	}
	if clone.Retry != nil {
		t.Error("Clone should have nil Retry")
	}
}

func TestWithHeader_NilHeaders(t *testing.T) {
	config := &RequestConfig{} // Headers is nil
	WithHeader("X-Test", "value")(config)
	if config.Headers == nil {
		t.Fatal("Headers should be initialized")
	}
	if config.Headers.Get("X-Test") != "value" {
		t.Error("Header not set")
	}
}

func TestWithQuery_NilQuery(t *testing.T) {
	config := &RequestConfig{} // Query is nil
	WithQuery("key", "value")(config)
	if config.Query == nil {
		t.Fatal("Query should be initialized")
	}
	if config.Query.Get("key") != "value" {
		t.Error("Query not set")
	}
}

func TestRequestConfig_DeepEquality(t *testing.T) {
	original := NewRequestConfig()
	original.Timeout = 5 * time.Second
	original.Headers.Set("X-Test", "value")
	original.Query.Set("key", "value")
	original.BasicAuth = &BasicAuth{Username: "user", Password: "pass"}
	original.BearerToken = "token"
	original.Retry = &RetryPolicy{MaxAttempts: 3}
	clone := original.Clone()
	// Compare values (not pointers)
	if clone.Timeout != original.Timeout {
		t.Error("Timeout mismatch")
	}
	if !reflect.DeepEqual(clone.Headers, original.Headers) {
		t.Error("Headers mismatch")
	}
	if !reflect.DeepEqual(clone.Query, original.Query) {
		t.Error("Query mismatch")
	}
	if clone.BasicAuth.Username != original.BasicAuth.Username {
		t.Error("BasicAuth username mismatch")
	}
	if clone.BearerToken != original.BearerToken {
		t.Error("BearerToken mismatch")
	}
	if clone.Retry.MaxAttempts != original.Retry.MaxAttempts {
		t.Error("Retry mismatch")
	}
}

// Tests for file upload options
func TestWithFile(t *testing.T) {
	config := NewRequestConfig()
	WithFile("document", "/path/to/file.txt")(config)
	if len(config.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(config.Files))
	}
	file := config.Files[0]
	if file.FieldName != "document" {
		t.Errorf("Expected field name 'document', got '%s'", file.FieldName)
	}
	if file.FilePath != "/path/to/file.txt" {
		t.Errorf("Expected file path '/path/to/file.txt', got '%s'", file.FilePath)
	}
}

func TestWithFileReader(t *testing.T) {
	config := NewRequestConfig()
	reader := struct{}{}
	WithFileReader("avatar", "photo.jpg", reader, 1024)(config)
	if len(config.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(config.Files))
	}
	file := config.Files[0]
	if file.FieldName != "avatar" {
		t.Errorf("Expected field name 'avatar', got '%s'", file.FieldName)
	}
	if file.FileName != "photo.jpg" {
		t.Errorf("Expected file name 'photo.jpg', got '%s'", file.FileName)
	}
	if file.Size != 1024 {
		t.Errorf("Expected size 1024, got %d", file.Size)
	}
	if file.Reader != reader {
		t.Error("Reader not set correctly")
	}
}

func TestWithProgress(t *testing.T) {
	var called bool
	callback := func(uploaded, total int64) {
		called = true
	}
	config := NewRequestConfig()
	WithProgress(callback)(config)
	if config.ProgressCallback == nil {
		t.Fatal("ProgressCallback should be set")
	}
	// Call the callback to verify it works
	config.ProgressCallback(100, 1000)
	if !called {
		t.Error("Callback was not called")
	}
}

func TestWithMultipleFiles(t *testing.T) {
	config := NewRequestConfig()
	WithFile("file1", "/path/to/file1.txt")(config)
	WithFile("file2", "/path/to/file2.txt")(config)
	WithFileReader("file3", "file3.txt", struct{}{}, 512)(config)
	if len(config.Files) != 3 {
		t.Fatalf("Expected 3 files, got %d", len(config.Files))
	}
}

func TestRequestConfig_CloneWithFiles(t *testing.T) {
	original := NewRequestConfig()
	WithFile("doc", "/path/to/doc.pdf")(original)
	WithProgress(func(uploaded, total int64) {})(original)
	clone := original.Clone()
	if len(clone.Files) != 1 {
		t.Errorf("Expected 1 file in clone, got %d", len(clone.Files))
	}
	if clone.ProgressCallback == nil {
		t.Error("ProgressCallback should be cloned")
	}
	// Verify independence
	WithFile("new", "/path/to/new.txt")(original)
	if len(clone.Files) != 1 {
		t.Error("Clone files should not be affected by original modification")
	}
}

func TestRequestConfig_MergeWithFiles(t *testing.T) {
	base := NewRequestConfig()
	WithFile("base", "/path/to/base.txt")(base)
	override := NewRequestConfig()
	WithFile("override", "/path/to/override.txt")(override)
	WithProgress(func(uploaded, total int64) {})(override)
	base.Merge(override)
	// Files should be appended
	if len(base.Files) != 2 {
		t.Errorf("Expected 2 files after merge, got %d", len(base.Files))
	}
	// ProgressCallback should be from override
	if base.ProgressCallback == nil {
		t.Error("ProgressCallback should be merged from override")
	}
}

func TestWithHeaders_NilHeaders(t *testing.T) {
	config := &RequestConfig{} // Headers is nil
	WithHeaders(map[string]string{
		"X-One": "1",
		"X-Two": "2",
	})(config)
	if config.Headers == nil {
		t.Fatal("Headers should be initialized")
	}
	if config.Headers.Get("X-One") != "1" {
		t.Error("X-One not set")
	}
	if config.Headers.Get("X-Two") != "2" {
		t.Error("X-Two not set")
	}
}

func TestWithQueryParams_NilQuery(t *testing.T) {
	config := &RequestConfig{} // Query is nil
	WithQueryParams(map[string]string{
		"page":  "1",
		"limit": "10",
	})(config)
	if config.Query == nil {
		t.Fatal("Query should be initialized")
	}
	if config.Query.Get("page") != "1" {
		t.Error("page not set")
	}
	if config.Query.Get("limit") != "10" {
		t.Error("limit not set")
	}
}

func TestWithContentType_NilHeaders(t *testing.T) {
	config := &RequestConfig{} // Headers is nil
	WithContentType("application/json")(config)
	if config.Headers == nil {
		t.Fatal("Headers should be initialized")
	}
	if config.Headers.Get("Content-Type") != "application/json" {
		t.Error("Content-Type not set")
	}
}

func TestWithAccept_NilHeaders(t *testing.T) {
	config := &RequestConfig{} // Headers is nil
	WithAccept("application/xml")(config)
	if config.Headers == nil {
		t.Fatal("Headers should be initialized")
	}
	if config.Headers.Get("Accept") != "application/xml" {
		t.Error("Accept not set")
	}
}

func TestWithUserAgent_NilHeaders(t *testing.T) {
	config := &RequestConfig{} // Headers is nil
	WithUserAgent("CustomAgent/1.0")(config)
	if config.Headers == nil {
		t.Fatal("Headers should be initialized")
	}
	if config.Headers.Get("User-Agent") != "CustomAgent/1.0" {
		t.Error("User-Agent not set")
	}
}

func TestRequestConfig_Merge_AllFields(t *testing.T) {
	base := NewRequestConfig()
	base.Timeout = 5 * time.Second
	base.Headers.Set("X-Base", "base")
	base.Query.Set("base-param", "value")
	base.BasicAuth = &BasicAuth{Username: "base-user", Password: "base-pass"}
	base.BearerToken = "base-token"
	base.Retry = &RetryPolicy{MaxAttempts: 3}
	base.Context = context.Background()
	base.Files = []FileUploadConfig{{FieldName: "base-file"}}
	base.ProgressCallback = func(uploaded, total int64) {}
	override := NewRequestConfig()
	override.Timeout = 10 * time.Second
	override.Headers.Set("X-Override", "override")
	override.Query.Set("override-param", "value")
	override.BasicAuth = &BasicAuth{Username: "override-user", Password: "override-pass"}
	override.BearerToken = "override-token"
	override.Retry = &RetryPolicy{MaxAttempts: 5}
	override.Context = context.WithValue(context.Background(), "key", "value")
	override.Files = []FileUploadConfig{{FieldName: "override-file"}}
	override.ProgressCallback = func(uploaded, total int64) {}
	base.Merge(override)
	// All override values should take precedence
	if base.Timeout != 10*time.Second {
		t.Error("Timeout not overridden")
	}
	if base.Headers.Get("X-Override") != "override" {
		t.Error("Override header not added")
	}
	if base.Query.Get("override-param") != "value" {
		t.Error("Override query not added")
	}
	if base.BasicAuth.Username != "override-user" {
		t.Error("BasicAuth not overridden")
	}
	if base.BearerToken != "override-token" {
		t.Error("BearerToken not overridden")
	}
	if base.Retry.MaxAttempts != 5 {
		t.Error("Retry not overridden")
	}
	if base.Context.Value("key") != "value" {
		t.Error("Context not overridden")
	}
	if len(base.Files) != 2 {
		t.Errorf("Files should be appended, got %d", len(base.Files))
	}
	if base.ProgressCallback == nil {
		t.Error("ProgressCallback not overridden")
	}
}

func TestRequestConfig_Merge_PartialOverride(t *testing.T) {
	base := NewRequestConfig()
	base.Timeout = 5 * time.Second
	base.BearerToken = "base-token"
	override := NewRequestConfig()
	// Only set timeout, leave other fields empty
	override.Timeout = 10 * time.Second
	base.Merge(override)
	// Timeout should be overridden
	if base.Timeout != 10*time.Second {
		t.Error("Timeout not overridden")
	}
	// BearerToken should remain unchanged
	if base.BearerToken != "base-token" {
		t.Error("BearerToken should not be changed")
	}
}
