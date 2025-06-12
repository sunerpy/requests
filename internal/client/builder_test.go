package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: http-client-refactor
// Property 3: RequestBuilder Preserves All Set Values
// For any sequence of RequestBuilder method calls (WithHeader, WithQuery, WithJSON, WithForm),
// the resulting Request from Build() SHALL contain all the values that were set.
// Validates: Requirements 2.1, 2.4, 2.5, 2.7, 2.8
func TestProperty3_BuilderPreservesHeaders(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Builder preserves all headers", prop.ForAll(
		func(headerKeys []string, headerValues []string) bool {
			// Use fixed number of headers to avoid case-sensitivity collisions
			// HTTP headers are case-insensitive, so "pK" and "pk" are the same
			if len(headerKeys) == 0 || len(headerValues) == 0 {
				return true
			}
			// Use only the first few unique keys
			numHeaders := len(headerKeys)
			if numHeaders > len(headerValues) {
				numHeaders = len(headerValues)
			}
			if numHeaders > 5 {
				numHeaders = 5
			}
			builder := NewRequest(MethodGet, "https://example.com")
			expectedHeaders := make(map[string]string)
			for i := 0; i < numHeaders; i++ {
				key := "X-Test-" + headerKeys[i] // Prefix to ensure uniqueness
				value := headerValues[i]
				builder.WithHeader(key, value)
				expectedHeaders[key] = value
			}
			req, err := builder.Build()
			if err != nil {
				return false
			}
			// Verify all headers are present
			for k, v := range expectedHeaders {
				if req.Headers.Get(k) != v {
					return false
				}
			}
			return true
		},
		gen.SliceOf(gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 })),
		gen.SliceOf(gen.AlphaString()),
	))
	properties.TestingRun(t)
}

func TestProperty3_BuilderPreservesQueryParams(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Builder preserves all query parameters", prop.ForAll(
		func(params map[string]string) bool {
			if len(params) == 0 {
				return true
			}
			builder := NewRequest(MethodGet, "https://example.com")
			for k, v := range params {
				builder.WithQuery(k, v)
			}
			req, err := builder.Build()
			if err != nil {
				return false
			}
			// Verify all query params are present in URL
			query := req.URL.Query()
			for k, v := range params {
				if query.Get(k) != v {
					return false
				}
			}
			return true
		},
		gen.MapOf(gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
			gen.AlphaString()),
	))
	properties.TestingRun(t)
}

func TestProperty3_BuilderPreservesJSON(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Builder preserves JSON body", prop.ForAll(
		func(id int, name string) bool {
			if name == "" {
				return true
			}
			data := map[string]any{"id": id, "name": name}
			builder := NewRequest(MethodPost, "https://example.com").WithJSON(data)
			req, err := builder.Build()
			if err != nil {
				return false
			}
			// Verify Content-Type header
			if req.Headers.Get("Content-Type") != "application/json" {
				return false
			}
			// Verify body content
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return false
			}
			var parsed map[string]any
			if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
				return false
			}
			// Check values (JSON numbers are float64)
			if int(parsed["id"].(float64)) != id {
				return false
			}
			if parsed["name"].(string) != name {
				return false
			}
			return true
		},
		gen.Int(),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

func TestProperty3_BuilderPreservesForm(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Builder preserves form data", prop.ForAll(
		func(formData map[string]string) bool {
			if len(formData) == 0 {
				return true
			}
			builder := NewRequest(MethodPost, "https://example.com").WithForm(formData)
			req, err := builder.Build()
			if err != nil {
				return false
			}
			// Verify Content-Type header
			if req.Headers.Get("Content-Type") != "application/x-www-form-urlencoded" {
				return false
			}
			// Verify body content
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return false
			}
			parsed, err := url.ParseQuery(string(bodyBytes))
			if err != nil {
				return false
			}
			for k, v := range formData {
				if parsed.Get(k) != v {
					return false
				}
			}
			return true
		},
		gen.MapOf(gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
			gen.AlphaString()),
	))
	properties.TestingRun(t)
}

func TestProperty3_BuilderPreservesMethod(t *testing.T) {
	methods := []Method{MethodGet, MethodPost, MethodPut, MethodDelete, MethodPatch, MethodHead, MethodOptions}
	for _, method := range methods {
		t.Run(string(method), func(t *testing.T) {
			builder := NewRequest(method, "https://example.com")
			req, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if req.Method != method {
				t.Errorf("Expected method %s, got %s", method, req.Method)
			}
		})
	}
}

func TestProperty3_BuilderPreservesURL(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Builder preserves URL", prop.ForAll(
		func(path string) bool {
			if path == "" || strings.ContainsAny(path, " \t\n") {
				return true
			}
			rawURL := "https://example.com/" + path
			builder := NewRequest(MethodGet, rawURL)
			req, err := builder.Build()
			if err != nil {
				return false
			}
			// URL should contain the path
			return strings.Contains(req.URL.String(), path)
		},
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

// Feature: http-client-refactor
// Property 4: RequestBuilder Validation
// For any RequestBuilder with missing required fields (URL or Method),
// calling Build() SHALL return a non-nil error describing the missing field.
// Validates: Requirements 2.2, 2.3
func TestProperty4_MissingURLReturnsError(t *testing.T) {
	builder := &RequestBuilder{
		method: MethodGet,
	}
	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error for missing URL")
	}
	// Should be a RequestError
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError, got %T", err)
	}
	// Should wrap ErrMissingURL
	if reqErr != nil && reqErr.Err != ErrMissingURL {
		t.Errorf("Expected ErrMissingURL, got %v", reqErr.Err)
	}
}

func TestProperty4_MissingMethodReturnsError(t *testing.T) {
	builder := &RequestBuilder{
		rawURL: "https://example.com",
	}
	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error for missing method")
	}
	// Should be a RequestError
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError, got %T", err)
	}
	// Should wrap ErrMissingMethod
	if reqErr != nil && reqErr.Err != ErrMissingMethod {
		t.Errorf("Expected ErrMissingMethod, got %v", reqErr.Err)
	}
}

func TestProperty4_InvalidMethodReturnsError(t *testing.T) {
	builder := NewRequest(Method("INVALID"), "https://example.com")
	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error for invalid method")
	}
	// Should be a RequestError
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError, got %T", err)
	}
	// Should wrap ErrInvalidMethod
	if reqErr != nil && reqErr.Err != ErrInvalidMethod {
		t.Errorf("Expected ErrInvalidMethod, got %v", reqErr.Err)
	}
}

func TestProperty4_InvalidURLReturnsError(t *testing.T) {
	builder := NewRequest(MethodGet, "://invalid-url")
	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
	// Should be a RequestError
	_, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError, got %T", err)
	}
}

func TestProperty4_ValidationPropertyTest(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Empty URL always returns error", prop.ForAll(
		func(method Method) bool {
			if !method.IsValid() {
				return true
			}
			builder := NewRequest(method, "")
			_, err := builder.Build()
			return err != nil
		},
		gen.OneConstOf(MethodGet, MethodPost, MethodPut, MethodDelete, MethodPatch),
	))
	properties.Property("Empty method always returns error", prop.ForAll(
		func(rawURL string) bool {
			if rawURL == "" {
				return true
			}
			builder := &RequestBuilder{rawURL: "https://example.com/" + rawURL}
			_, err := builder.Build()
			return err != nil
		},
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

// Unit tests for builder methods
func TestNewRequest(t *testing.T) {
	builder := NewRequest(MethodPost, "https://api.example.com/users")
	if builder.method != MethodPost {
		t.Errorf("Expected POST, got %s", builder.method)
	}
	if builder.rawURL != "https://api.example.com/users" {
		t.Errorf("Expected URL, got %s", builder.rawURL)
	}
}

func TestBuilderConvenienceFunctions(t *testing.T) {
	tests := []struct {
		name     string
		builder  *RequestBuilder
		expected Method
	}{
		{"NewGetRequest", NewGetRequest("https://example.com"), MethodGet},
		{"NewPostRequest", NewPostRequest("https://example.com"), MethodPost},
		{"NewPutRequest", NewPutRequest("https://example.com"), MethodPut},
		{"NewDeleteRequest", NewDeleteRequest("https://example.com"), MethodDelete},
		{"NewPatchRequest", NewPatchRequest("https://example.com"), MethodPatch},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.builder.method != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, tc.builder.method)
			}
		})
	}
}

func TestBuilderWithHeaders(t *testing.T) {
	builder := NewRequest(MethodGet, "https://example.com").
		WithHeaders(map[string]string{
			"X-Custom-Header": "value1",
			"Accept":          "application/json",
		})
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if req.Headers.Get("X-Custom-Header") != "value1" {
		t.Error("Missing X-Custom-Header")
	}
	if req.Headers.Get("Accept") != "application/json" {
		t.Error("Missing Accept header")
	}
}

func TestBuilderWithQueryParams(t *testing.T) {
	builder := NewRequest(MethodGet, "https://example.com/search").
		WithQueryParams(map[string]string{
			"q":    "golang",
			"page": "1",
		})
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	query := req.URL.Query()
	if query.Get("q") != "golang" {
		t.Error("Missing q parameter")
	}
	if query.Get("page") != "1" {
		t.Error("Missing page parameter")
	}
}

func TestBuilderWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")
	builder := NewRequest(MethodGet, "https://example.com").WithContext(ctx)
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if req.Context.Value("key") != "value" {
		t.Error("Context not preserved")
	}
}

func TestBuilderWithTimeout(t *testing.T) {
	builder := NewRequest(MethodGet, "https://example.com").WithTimeout(5 * time.Second)
	if builder.timeout != 5*time.Second {
		t.Errorf("Expected 5s timeout, got %v", builder.timeout)
	}
}

func TestBuilderWithBasicAuth(t *testing.T) {
	builder := NewRequest(MethodGet, "https://example.com").WithBasicAuth("user", "pass")
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	auth := req.Headers.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		t.Error("Missing Basic auth header")
	}
}

func TestBuilderWithBearerToken(t *testing.T) {
	builder := NewRequest(MethodGet, "https://example.com").WithBearerToken("my-token")
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	auth := req.Headers.Get("Authorization")
	if auth != "Bearer my-token" {
		t.Errorf("Expected 'Bearer my-token', got '%s'", auth)
	}
}

func TestBuilderWithBodyBytes(t *testing.T) {
	data := []byte("raw body content")
	builder := NewRequest(MethodPost, "https://example.com").WithBodyBytes(data)
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	if !bytes.Equal(bodyBytes, data) {
		t.Error("Body content mismatch")
	}
}

func TestBuilderClone(t *testing.T) {
	original := NewRequest(MethodPost, "https://example.com").
		WithHeader("X-Original", "value").
		WithQuery("key", "value").
		WithJSON(map[string]string{"name": "test"})
	clone := original.Clone()
	// Modify original
	original.WithHeader("X-Modified", "modified")
	original.WithQuery("new", "param")
	// Clone should not be affected
	if clone.headers.Get("X-Modified") != "" {
		t.Error("Clone was affected by original modification")
	}
	cloneQuery := clone.GetQuery()
	if cloneQuery.Get("new") != "" {
		t.Error("Clone query was affected by original modification")
	}
}

func TestBuilderGetters(t *testing.T) {
	builder := NewRequest(MethodPost, "https://example.com").
		WithHeader("X-Test", "value").
		WithQuery("q", "search").
		WithTimeout(10 * time.Second)
	if builder.GetMethod() != MethodPost {
		t.Error("GetMethod failed")
	}
	if builder.GetURL() != "https://example.com" {
		t.Error("GetURL failed")
	}
	if builder.GetHeaders().Get("X-Test") != "value" {
		t.Error("GetHeaders failed")
	}
	if builder.GetQuery().Get("q") != "search" {
		t.Error("GetQuery failed")
	}
	if builder.GetTimeout() != 10*time.Second {
		t.Error("GetTimeout failed")
	}
}

func TestBuilderWithFormValues(t *testing.T) {
	values := url.Values{}
	values.Add("field1", "value1")
	values.Add("field2", "value2")
	builder := NewRequest(MethodPost, "https://example.com").WithFormValues(values)
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if req.Headers.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Error("Wrong Content-Type")
	}
	bodyBytes, _ := io.ReadAll(req.Body)
	parsed, _ := url.ParseQuery(string(bodyBytes))
	if parsed.Get("field1") != "value1" || parsed.Get("field2") != "value2" {
		t.Error("Form values not preserved")
	}
}

func TestBuilderChaining(t *testing.T) {
	req, err := NewRequest(MethodPost, "https://api.example.com/users").
		WithHeader("Content-Type", "application/json").
		WithHeader("Accept", "application/json").
		WithQuery("version", "v1").
		WithBearerToken("secret-token").
		WithTimeout(30 * time.Second).
		WithJSON(map[string]string{"name": "John"}).
		Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if req.Method != MethodPost {
		t.Error("Wrong method")
	}
	if req.Headers.Get("Authorization") != "Bearer secret-token" {
		t.Error("Missing auth header")
	}
	if req.URL.Query().Get("version") != "v1" {
		t.Error("Missing query param")
	}
}

func TestBuilderWithFile(t *testing.T) {
	content := bytes.NewReader([]byte("file content"))
	builder := NewRequest(MethodPost, "https://example.com/upload").
		WithFile("document", "test.txt", content)
	files := builder.GetFiles()
	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}
	if files[0].FieldName != "document" {
		t.Error("Wrong field name")
	}
	if files[0].FileName != "test.txt" {
		t.Error("Wrong file name")
	}
}

func TestBuilderErrorAccumulation(t *testing.T) {
	// Create a type that cannot be marshaled to JSON
	type badType struct {
		Ch chan int `json:"ch"`
	}
	builder := NewRequest(MethodPost, "https://example.com").
		WithJSON(badType{Ch: make(chan int)})
	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
	encErr, ok := err.(*EncodeError)
	if !ok {
		t.Errorf("Expected EncodeError, got %T", err)
	}
	if encErr != nil && encErr.ContentType != "application/json" {
		t.Error("Wrong content type in error")
	}
}

func TestBuilderMergesExistingQueryParams(t *testing.T) {
	builder := NewRequest(MethodGet, "https://example.com/search?existing=param").
		WithQuery("new", "value")
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	query := req.URL.Query()
	if query.Get("existing") != "param" {
		t.Error("Existing query param lost")
	}
	if query.Get("new") != "value" {
		t.Error("New query param not added")
	}
}

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"f", "Zg=="},
		{"fo", "Zm8="},
		{"foo", "Zm9v"},
		{"foob", "Zm9vYg=="},
		{"fooba", "Zm9vYmE="},
		{"foobar", "Zm9vYmFy"},
		{"user:pass", "dXNlcjpwYXNz"},
	}
	for _, tc := range tests {
		result := base64Encode([]byte(tc.input))
		if result != tc.expected {
			t.Errorf("base64Encode(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestBuilderWithXML(t *testing.T) {
	type Person struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
	}
	builder := NewRequest(MethodPost, "https://example.com").
		WithXML(Person{Name: "John", Age: 30})
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if req.Headers.Get("Content-Type") != "application/xml" {
		t.Error("Wrong Content-Type for XML")
	}
	bodyBytes, _ := io.ReadAll(req.Body)
	if !strings.Contains(string(bodyBytes), "<name>John</name>") {
		t.Error("XML body not correct")
	}
}

func TestBuilderPreservesAllValues(t *testing.T) {
	// Comprehensive test that all builder values are preserved
	builder := NewRequest(MethodPost, "https://api.example.com/test").
		WithHeader("X-Custom", "custom-value").
		WithHeader("Accept", "application/json").
		WithQuery("page", "1").
		WithQuery("limit", "10").
		WithBearerToken("token123").
		WithTimeout(15 * time.Second)
	// Get values before build
	method := builder.GetMethod()
	_ = builder.GetURL() // rawURL
	headers := builder.GetHeaders()
	query := builder.GetQuery()
	timeout := builder.GetTimeout()
	// Build and verify
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	// Verify method
	if req.Method != method {
		t.Errorf("Method mismatch: expected %s, got %s", method, req.Method)
	}
	// Verify URL contains original
	if !strings.Contains(req.URL.String(), "api.example.com/test") {
		t.Error("URL not preserved")
	}
	// Verify headers
	for k := range headers {
		if req.Headers.Get(k) != headers.Get(k) {
			t.Errorf("Header %s not preserved", k)
		}
	}
	// Verify query params
	reqQuery := req.URL.Query()
	for k := range query {
		if reqQuery.Get(k) != query.Get(k) {
			t.Errorf("Query param %s not preserved", k)
		}
	}
	// Verify timeout was set
	if timeout != 15*time.Second {
		t.Error("Timeout not preserved in builder")
	}
}

func TestBuilderWithBody(t *testing.T) {
	body := strings.NewReader("custom body content")
	builder := NewRequest(MethodPost, "https://example.com").WithBody(body)
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	content, _ := io.ReadAll(req.Body)
	if string(content) != "custom body content" {
		t.Error("Body not preserved")
	}
}

func TestBuilderWithMethodChange(t *testing.T) {
	builder := NewRequest(MethodGet, "https://example.com").
		WithMethod(MethodPost)
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if req.Method != MethodPost {
		t.Errorf("Expected POST, got %s", req.Method)
	}
}

func TestBuilderWithURLChange(t *testing.T) {
	builder := NewRequest(MethodGet, "https://old.example.com").
		WithURL("https://new.example.com")
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if !strings.Contains(req.URL.String(), "new.example.com") {
		t.Error("URL not updated")
	}
}

func TestBuilderHeadersCloned(t *testing.T) {
	builder := NewRequest(MethodGet, "https://example.com").
		WithHeader("X-Test", "value")
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	// Modify builder headers after build
	builder.WithHeader("X-New", "new-value")
	// Request headers should not be affected
	if req.Headers.Get("X-New") != "" {
		t.Error("Request headers were modified after build")
	}
}

func TestBuilderGetBodyBytes(t *testing.T) {
	data := map[string]string{"key": "value"}
	builder := NewRequest(MethodPost, "https://example.com").WithJSON(data)
	bodyBytes := builder.GetBodyBytes()
	if len(bodyBytes) == 0 {
		t.Error("GetBodyBytes returned empty")
	}
	var parsed map[string]string
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		t.Fatalf("Failed to parse body bytes: %v", err)
	}
	if !reflect.DeepEqual(data, parsed) {
		t.Error("Body bytes content mismatch")
	}
}

// Additional tests for coverage
func TestBuilderWithEncoder(t *testing.T) {
	// Create a custom encoder
	encoder := &testEncoder{
		contentType: "application/custom",
		encodeFunc: func(v any) ([]byte, error) {
			return []byte("custom-encoded"), nil
		},
	}
	builder := NewRequest(MethodPost, "https://example.com").
		WithEncoder(encoder, map[string]string{"key": "value"})
	req, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if req.Headers.Get("Content-Type") != "application/custom" {
		t.Errorf("Expected application/custom, got %s", req.Headers.Get("Content-Type"))
	}
	bodyBytes, _ := io.ReadAll(req.Body)
	if string(bodyBytes) != "custom-encoded" {
		t.Errorf("Expected 'custom-encoded', got '%s'", string(bodyBytes))
	}
}

type testEncoder struct {
	contentType string
	encodeFunc  func(v any) ([]byte, error)
}

func (e *testEncoder) Encode(v any) ([]byte, error) {
	return e.encodeFunc(v)
}

func (e *testEncoder) ContentType() string {
	return e.contentType
}

func TestBuilderWithEncoderError(t *testing.T) {
	encoder := &testEncoder{
		contentType: "application/custom",
		encodeFunc: func(v any) ([]byte, error) {
			return nil, io.ErrUnexpectedEOF
		},
	}
	builder := NewRequest(MethodPost, "https://example.com").
		WithEncoder(encoder, map[string]string{"key": "value"})
	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error from encoder")
	}
	encErr, ok := err.(*EncodeError)
	if !ok {
		t.Errorf("Expected EncodeError, got %T", err)
	}
	if encErr != nil && encErr.ContentType != "application/custom" {
		t.Errorf("Expected content type 'application/custom', got '%s'", encErr.ContentType)
	}
}

func TestBuilderGetForm(t *testing.T) {
	builder := NewRequest(MethodPost, "https://example.com").
		WithForm(map[string]string{
			"field1": "value1",
			"field2": "value2",
		})
	form := builder.GetForm()
	if form.Get("field1") != "value1" {
		t.Error("field1 not found in form")
	}
	if form.Get("field2") != "value2" {
		t.Error("field2 not found in form")
	}
}

func TestBuilderCloneWithBodyBytes(t *testing.T) {
	original := NewRequest(MethodPost, "https://example.com").
		WithJSON(map[string]string{"key": "value"})
	clone := original.Clone()
	// Both should have body bytes
	if len(original.GetBodyBytes()) == 0 {
		t.Error("Original should have body bytes")
	}
	if len(clone.GetBodyBytes()) == 0 {
		t.Error("Clone should have body bytes")
	}
	// Clone body should be readable
	cloneReq, err := clone.Build()
	if err != nil {
		t.Fatalf("Clone build failed: %v", err)
	}
	bodyBytes, _ := io.ReadAll(cloneReq.Body)
	if len(bodyBytes) == 0 {
		t.Error("Clone body should be readable")
	}
}

func TestBuilderCloneWithFiles(t *testing.T) {
	original := NewRequest(MethodPost, "https://example.com").
		WithFile("file1", "test1.txt", strings.NewReader("content1")).
		WithFile("file2", "test2.txt", strings.NewReader("content2"))
	clone := original.Clone()
	originalFiles := original.GetFiles()
	cloneFiles := clone.GetFiles()
	if len(originalFiles) != 2 {
		t.Errorf("Original should have 2 files, got %d", len(originalFiles))
	}
	if len(cloneFiles) != 2 {
		t.Errorf("Clone should have 2 files, got %d", len(cloneFiles))
	}
}

func TestBuilderCloneWithForm(t *testing.T) {
	original := NewRequest(MethodPost, "https://example.com").
		WithForm(map[string]string{"key": "value"})
	clone := original.Clone()
	// Modify original form
	original.WithForm(map[string]string{"new": "data"})
	// Clone should not be affected
	cloneForm := clone.GetForm()
	if cloneForm.Get("new") != "" {
		t.Error("Clone form was affected by original modification")
	}
}

func TestBuilderWithXMLError(t *testing.T) {
	// Create a type that cannot be marshaled to XML
	type badType struct {
		Ch chan int `xml:"ch"`
	}
	builder := NewRequest(MethodPost, "https://example.com").
		WithXML(badType{Ch: make(chan int)})
	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error for unmarshalable XML type")
	}
	encErr, ok := err.(*EncodeError)
	if !ok {
		t.Errorf("Expected EncodeError, got %T", err)
	}
	if encErr != nil && encErr.ContentType != "application/xml" {
		t.Error("Wrong content type in error")
	}
}

// Tests for short convenience constructors
func TestNewGet(t *testing.T) {
	builder := NewGet("https://example.com/test")
	if builder.GetMethod() != MethodGet {
		t.Errorf("Expected GET, got %s", builder.GetMethod())
	}
	if builder.GetURL() != "https://example.com/test" {
		t.Errorf("Expected URL https://example.com/test, got %s", builder.GetURL())
	}
}

func TestNewPost(t *testing.T) {
	builder := NewPost("https://example.com/test")
	if builder.GetMethod() != MethodPost {
		t.Errorf("Expected POST, got %s", builder.GetMethod())
	}
}

func TestNewPut(t *testing.T) {
	builder := NewPut("https://example.com/test")
	if builder.GetMethod() != MethodPut {
		t.Errorf("Expected PUT, got %s", builder.GetMethod())
	}
}

func TestNewDelete(t *testing.T) {
	builder := NewDelete("https://example.com/test")
	if builder.GetMethod() != MethodDelete {
		t.Errorf("Expected DELETE, got %s", builder.GetMethod())
	}
}

func TestNewPatch(t *testing.T) {
	builder := NewPatch("https://example.com/test")
	if builder.GetMethod() != MethodPatch {
		t.Errorf("Expected PATCH, got %s", builder.GetMethod())
	}
}

// Property test for convenience constructors
func TestProperty_ConvenienceConstructorsSetCorrectMethod(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)
	properties.Property("NewGet sets GET method", prop.ForAll(
		func(path string) bool {
			url := "https://example.com/" + path
			builder := NewGet(url)
			return builder.GetMethod() == MethodGet && builder.GetURL() == url
		},
		gen.AlphaString(),
	))
	properties.Property("NewPost sets POST method", prop.ForAll(
		func(path string) bool {
			url := "https://example.com/" + path
			builder := NewPost(url)
			return builder.GetMethod() == MethodPost && builder.GetURL() == url
		},
		gen.AlphaString(),
	))
	properties.Property("NewPut sets PUT method", prop.ForAll(
		func(path string) bool {
			url := "https://example.com/" + path
			builder := NewPut(url)
			return builder.GetMethod() == MethodPut && builder.GetURL() == url
		},
		gen.AlphaString(),
	))
	properties.Property("NewDelete sets DELETE method", prop.ForAll(
		func(path string) bool {
			url := "https://example.com/" + path
			builder := NewDelete(url)
			return builder.GetMethod() == MethodDelete && builder.GetURL() == url
		},
		gen.AlphaString(),
	))
	properties.Property("NewPatch sets PATCH method", prop.ForAll(
		func(path string) bool {
			url := "https://example.com/" + path
			builder := NewPatch(url)
			return builder.GetMethod() == MethodPatch && builder.GetURL() == url
		},
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

// Tests for Do() and DoJSON() methods
func TestBuilderDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("key") != "value" {
			t.Error("Query param not received")
		}
		if r.Header.Get("X-Custom") != "header" {
			t.Error("Custom header not received")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()
	resp, err := NewGet(server.URL).
		WithQuery("key", "value").
		WithHeader("X-Custom", "header").
		Do()
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	if resp.Text() != "OK" {
		t.Errorf("Expected 'OK', got '%s'", resp.Text())
	}
}

func TestBuilderDo_Error(t *testing.T) {
	// Test with invalid URL
	_, err := NewGet("://invalid").Do()
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestBuilderDo_MissingURL(t *testing.T) {
	builder := NewRequest(MethodGet, "")
	_, err := builder.Do()
	if err == nil {
		t.Error("Expected error for missing URL")
	}
}

type BuilderTestUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestDoJSON(t *testing.T) {
	expected := BuilderTestUser{ID: 1, Name: "Test", Email: "test@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := DoJSON[BuilderTestUser](NewGet(server.URL))
	if err != nil {
		t.Fatalf("DoJSON failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
	if result.StatusCode() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", result.StatusCode())
	}
}

func TestDoJSON_WithPost(t *testing.T) {
	expected := BuilderTestUser{ID: 1, Name: "Created", Email: "created@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := DoJSON[BuilderTestUser](
		NewPost(server.URL).
			WithJSON(map[string]string{"name": "Test"}),
	)
	if err != nil {
		t.Fatalf("DoJSON failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestDoJSON_Error(t *testing.T) {
	_, err := DoJSON[BuilderTestUser](NewGet("://invalid"))
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestDoJSON_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	_, err := DoJSON[BuilderTestUser](NewGet(server.URL))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestDoXML(t *testing.T) {
	type XMLUser struct {
		ID    int    `xml:"id"`
		Name  string `xml:"name"`
		Email string `xml:"email"`
	}
	expected := XMLUser{ID: 1, Name: "Test", Email: "test@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<XMLUser><id>1</id><name>Test</name><email>test@example.com</email></XMLUser>`))
	}))
	defer server.Close()
	result, err := DoXML[XMLUser](NewGet(server.URL))
	if err != nil {
		t.Fatalf("DoXML failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestDoXML_Error(t *testing.T) {
	_, err := DoXML[BuilderTestUser](NewGet("://invalid"))
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestDoXML_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid xml"))
	}))
	defer server.Close()
	_, err := DoXML[BuilderTestUser](NewGet(server.URL))
	if err == nil {
		t.Error("Expected error for invalid XML")
	}
}

// Property test for Do() method
func TestProperty_DoExecutesRequest(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)
	properties.Property("Do executes request with correct method and headers", prop.ForAll(
		func(headerKey, headerValue string) bool {
			if headerKey == "" {
				return true
			}
			var receivedMethod string
			var receivedHeader string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				receivedHeader = r.Header.Get("X-" + headerKey)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			_, err := NewGet(server.URL).
				WithHeader("X-"+headerKey, headerValue).
				Do()
			if err != nil {
				return false
			}
			return receivedMethod == "GET" && receivedHeader == headerValue
		},
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}
