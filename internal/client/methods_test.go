package client

import (
	"context"
	"encoding/json"
	"encoding/xml"
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

// TestAPIUser is a test struct for API tests.
type TestAPIUser struct {
	ID    int    `json:"id" xml:"id"`
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

// Feature: http-client-refactor
// Property 11: Generic HTTP Methods Parse Response Correctly
// For any valid JSON response from a server, calling GetJSON[T]() or PostJSON[T]()
// SHALL return the correctly parsed type T wrapped in Result[T].
// Validates: Requirements 6.2, 6.4
func TestProperty11_GenericMethodsParseJSON(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("GetJSON parses response correctly", prop.ForAll(
		func(id int, name, email string) bool {
			if name == "" || email == "" {
				return true
			}
			expected := TestAPIUser{ID: id, Name: name, Email: email}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(expected)
			}))
			defer server.Close()
			result, err := GetJSON[TestAPIUser](server.URL)
			if err != nil {
				return false
			}
			return reflect.DeepEqual(expected, result.Data())
		},
		gen.Int(),
		gen.AlphaString(),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

func TestProperty11_PostJSONParsesResponse(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("PostJSON parses response correctly", prop.ForAll(
		func(id int, name string) bool {
			if name == "" {
				return true
			}
			input := map[string]any{"name": name}
			expected := TestAPIUser{ID: id, Name: name, Email: "created@example.com"}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(expected)
			}))
			defer server.Close()
			result, err := PostJSON[TestAPIUser](server.URL, input)
			if err != nil {
				return false
			}
			return reflect.DeepEqual(expected, result.Data())
		},
		gen.Int(),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

func TestProperty11_GenericMethodsReturnResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TestAPIUser{ID: 1, Name: "Test", Email: "test@example.com"})
	}))
	defer server.Close()
	// Use Result[T] to get response
	result, err := GetJSON[TestAPIUser](server.URL)
	if err != nil {
		t.Fatalf("GetJSON failed: %v", err)
	}
	// Response should be accessible via Result
	if result.Response() == nil {
		t.Fatal("Response is nil")
	}
	// Response should have correct status
	if result.StatusCode() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", result.StatusCode())
	}
	// Response should have headers
	if result.Headers().Get("X-Custom-Header") != "test-value" {
		t.Error("Custom header not in response")
	}
}

// Feature: http-client-refactor
// Property 12: Generic HTTP Methods Error Handling
// For any response that cannot be parsed as type T, the generic HTTP methods
// SHALL return a zero value for T and a non-nil error.
// Validates: Requirements 6.3
func TestProperty12_GenericMethodsErrorHandling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Invalid JSON returns error with zero value", prop.ForAll(
		func(invalidJSON string) bool {
			// Ensure it's actually invalid JSON
			var test TestAPIUser
			if json.Unmarshal([]byte(invalidJSON), &test) == nil {
				return true
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(invalidJSON))
			}))
			defer server.Close()
			result, err := GetJSON[TestAPIUser](server.URL)
			// Should return error
			if err == nil {
				return false
			}
			// Should return zero value in Result
			var zero TestAPIUser
			if !reflect.DeepEqual(result.Data(), zero) {
				return false
			}
			return true
		},
		gen.AnyString().SuchThat(func(s string) bool {
			return len(s) > 0 && s[0] != '{' && s[0] != '['
		}),
	))
	properties.TestingRun(t)
}

func TestProperty12_NetworkErrorReturnsError(t *testing.T) {
	// Use an invalid URL to trigger network error
	_, err := GetJSON[TestAPIUser]("http://invalid.invalid.invalid:12345")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestProperty12_ServerErrorReturnsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()
	result, err := GetJSON[TestAPIUser](server.URL)
	// Should return error (invalid JSON)
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
	// Should return zero value
	var zero TestAPIUser
	if !reflect.DeepEqual(result.Data(), zero) {
		t.Error("Expected zero value")
	}
}

// Unit tests for HTTP methods
func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
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
	if resp.Text() != "OK" {
		t.Errorf("Expected 'OK', got '%s'", resp.Text())
	}
}

func TestPost(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()
	resp, err := Post(server.URL, map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201, got %d", resp.StatusCode)
	}
	if receivedBody == "" {
		t.Error("Body not sent")
	}
}

func TestPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	resp, err := Put(server.URL, map[string]string{"name": "updated"})
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
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

func TestPatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	resp, err := Patch(server.URL, map[string]string{"name": "patched"})
	if err != nil {
		t.Fatalf("Patch failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestHead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			t.Errorf("Expected HEAD, got %s", r.Method)
		}
		w.Header().Set("X-Custom", "value")
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
	if resp.Headers.Get("X-Custom") != "value" {
		t.Error("Custom header not received")
	}
}

func TestOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			t.Errorf("Expected OPTIONS, got %s", r.Method)
		}
		w.Header().Set("Allow", "GET, POST, OPTIONS")
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

func TestGetWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	_, err := GetWithContext(ctx, server.URL)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestGetWithOptions(t *testing.T) {
	var receivedHeader string
	var receivedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("X-Custom")
		receivedQuery = r.URL.Query().Get("page")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	_, err := Get(server.URL,
		WithHeader("X-Custom", "value"),
		WithQuery("page", "1"),
	)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if receivedHeader != "value" {
		t.Error("Custom header not sent")
	}
	if receivedQuery != "1" {
		t.Error("Query param not sent")
	}
}

func TestPostWithFormData(t *testing.T) {
	var receivedContentType string
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	formData := url.Values{}
	formData.Set("field1", "value1")
	formData.Set("field2", "value2")
	_, err := Post(server.URL, formData)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if receivedContentType != "application/x-www-form-urlencoded" {
		t.Errorf("Expected form content type, got %s", receivedContentType)
	}
	if receivedBody == "" {
		t.Error("Form data not sent")
	}
}

func TestGetXML(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Test", Email: "test@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := GetXML[TestAPIUser](server.URL)
	if err != nil {
		t.Fatalf("GetXML failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestPrepareBody(t *testing.T) {
	tests := []struct {
		name        string
		body        any
		expectCT    string
		expectError bool
	}{
		{"nil", nil, "", false},
		{"string", "hello", "", false},
		{"bytes", []byte("hello"), "", false},
		{"url.Values", url.Values{"key": {"value"}}, "application/x-www-form-urlencoded", false},
		{"struct", TestAPIUser{ID: 1}, "application/json", false},
		{"map", map[string]string{"key": "value"}, "application/json", false},
		{"io.Reader", strings.NewReader("hello"), "", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader, ct, err := PrepareBody(tc.body)
			if tc.expectError && err == nil {
				t.Error("Expected error")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if ct != tc.expectCT {
				t.Errorf("Expected content type '%s', got '%s'", tc.expectCT, ct)
			}
			if tc.body != nil && reader == nil && !tc.expectError {
				t.Error("Expected reader for non-nil body")
			}
		})
	}
}

func TestSetDefaultClient(t *testing.T) {
	original := DefaultHTTPClient
	mockClient := &http.Client{Timeout: 1 * time.Second}
	SetDefaultClient(mockClient)
	if DefaultHTTPClient != mockClient {
		t.Error("Default client not set")
	}
	// Restore
	SetDefaultClient(original)
}

func TestGenericMethodsWithAuth(t *testing.T) {
	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TestAPIUser{ID: 1, Name: "Test", Email: "test@example.com"})
	}))
	defer server.Close()
	_, err := GetJSON[TestAPIUser](server.URL, WithBearerToken("my-token"))
	if err != nil {
		t.Fatalf("GetJSON failed: %v", err)
	}
	if receivedAuth != "Bearer my-token" {
		t.Errorf("Expected 'Bearer my-token', got '%s'", receivedAuth)
	}
}

func TestPutJSON(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Updated", Email: "updated@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := PutJSON[TestAPIUser](server.URL, map[string]string{"name": "Updated"})
	if err != nil {
		t.Fatalf("PutJSON failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestDeleteJSON(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Deleted", Email: "deleted@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := DeleteJSON[TestAPIUser](server.URL)
	if err != nil {
		t.Fatalf("DeleteJSON failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestPatchJSON(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Patched", Email: "patched@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := PatchJSON[TestAPIUser](server.URL, map[string]string{"name": "Patched"})
	if err != nil {
		t.Fatalf("PatchJSON failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

// Additional tests for coverage
func TestDeleteWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	ctx := context.Background()
	resp, err := DeleteWithContext(ctx, server.URL)
	if err != nil {
		t.Fatalf("DeleteWithContext failed: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected 204, got %d", resp.StatusCode)
	}
}

func TestHeadWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			t.Errorf("Expected HEAD, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	ctx := context.Background()
	resp, err := HeadWithContext(ctx, server.URL)
	if err != nil {
		t.Fatalf("HeadWithContext failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestOptionsWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			t.Errorf("Expected OPTIONS, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	ctx := context.Background()
	resp, err := OptionsWithContext(ctx, server.URL)
	if err != nil {
		t.Fatalf("OptionsWithContext failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestPostXML(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Created", Email: "created@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := PostXML[TestAPIUser](server.URL, TestAPIUser{Name: "Test"})
	if err != nil {
		t.Fatalf("PostXML failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestPostXMLWithContext(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Created", Email: "created@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	ctx := context.Background()
	result, err := PostXMLWithContext[TestAPIUser](ctx, server.URL, TestAPIUser{Name: "Test"})
	if err != nil {
		t.Fatalf("PostXMLWithContext failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestGetXMLWithContext(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Test", Email: "test@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	ctx := context.Background()
	result, err := GetXMLWithContext[TestAPIUser](ctx, server.URL)
	if err != nil {
		t.Fatalf("GetXMLWithContext failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestPostWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()
	ctx := context.Background()
	resp, err := PostWithContext(ctx, server.URL, map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("PostWithContext failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201, got %d", resp.StatusCode)
	}
}

func TestPutWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	ctx := context.Background()
	resp, err := PutWithContext(ctx, server.URL, map[string]string{"name": "updated"})
	if err != nil {
		t.Fatalf("PutWithContext failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestPatchWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	ctx := context.Background()
	resp, err := PatchWithContext(ctx, server.URL, map[string]string{"name": "patched"})
	if err != nil {
		t.Fatalf("PatchWithContext failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestPrepareBody_UnmarshalableType(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	_, _, err := PrepareBody(BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPostJSONWithContext(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Created", Email: "created@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	ctx := context.Background()
	result, err := PostJSONWithContext[TestAPIUser](ctx, server.URL, map[string]string{"name": "Test"})
	if err != nil {
		t.Fatalf("PostJSONWithContext failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestPutJSONWithContext(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Updated", Email: "updated@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	ctx := context.Background()
	result, err := PutJSONWithContext[TestAPIUser](ctx, server.URL, map[string]string{"name": "Updated"})
	if err != nil {
		t.Fatalf("PutJSONWithContext failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestDeleteJSONWithContext(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Deleted", Email: "deleted@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	ctx := context.Background()
	result, err := DeleteJSONWithContext[TestAPIUser](ctx, server.URL)
	if err != nil {
		t.Fatalf("DeleteJSONWithContext failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestPatchJSONWithContext(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Patched", Email: "patched@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	ctx := context.Background()
	result, err := PatchJSONWithContext[TestAPIUser](ctx, server.URL, map[string]string{"name": "Patched"})
	if err != nil {
		t.Fatalf("PatchJSONWithContext failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestPostJSONWithContext_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	ctx := context.Background()
	_, err := PostJSONWithContext[TestAPIUser](ctx, "https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPutJSONWithContext_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	ctx := context.Background()
	_, err := PutJSONWithContext[TestAPIUser](ctx, "https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPatchJSONWithContext_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	ctx := context.Background()
	_, err := PatchJSONWithContext[TestAPIUser](ctx, "https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPostXMLWithContext_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `xml:"ch"`
	}
	ctx := context.Background()
	_, err := PostXMLWithContext[TestAPIUser](ctx, "https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestGetXMLWithContext_Error(t *testing.T) {
	ctx := context.Background()
	_, err := GetXMLWithContext[TestAPIUser](ctx, "://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestDeleteJSONWithContext_Error(t *testing.T) {
	ctx := context.Background()
	_, err := DeleteJSONWithContext[TestAPIUser](ctx, "://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPostWithContext_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	ctx := context.Background()
	_, err := PostWithContext(ctx, "https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPutWithContext_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	ctx := context.Background()
	_, err := PutWithContext(ctx, "https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPatchWithContext_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	ctx := context.Background()
	_, err := PatchWithContext(ctx, "https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPostJSONWithContext_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	ctx := context.Background()
	_, err := PostJSONWithContext[TestAPIUser](ctx, server.URL, map[string]string{"name": "test"})
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

func TestPutJSONWithContext_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	ctx := context.Background()
	_, err := PutJSONWithContext[TestAPIUser](ctx, server.URL, map[string]string{"name": "test"})
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

func TestPatchJSONWithContext_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	ctx := context.Background()
	_, err := PatchJSONWithContext[TestAPIUser](ctx, server.URL, map[string]string{"name": "test"})
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

func TestDeleteJSONWithContext_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	ctx := context.Background()
	_, err := DeleteJSONWithContext[TestAPIUser](ctx, server.URL)
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

func TestGetXMLWithContext_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid xml"))
	}))
	defer server.Close()
	ctx := context.Background()
	_, err := GetXMLWithContext[TestAPIUser](ctx, server.URL)
	if err == nil {
		t.Error("Expected error for invalid XML response")
	}
}

func TestPostXMLWithContext_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid xml"))
	}))
	defer server.Close()
	ctx := context.Background()
	_, err := PostXMLWithContext[TestAPIUser](ctx, server.URL, TestAPIUser{Name: "test"})
	if err == nil {
		t.Error("Expected error for invalid XML response")
	}
}

// Tests for Result[T] API - verifying response metadata access
func TestGetJSON_ResultMetadata(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Test", Email: "test@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "custom-value")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := GetJSON[TestAPIUser](server.URL)
	if err != nil {
		t.Fatalf("GetJSON failed: %v", err)
	}
	// Verify data
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
	// Verify response metadata
	if result.Response() == nil {
		t.Fatal("Response is nil")
	}
	if result.StatusCode() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", result.StatusCode())
	}
	if result.Headers().Get("X-Custom-Header") != "custom-value" {
		t.Error("Custom header not in response")
	}
	if !result.IsSuccess() {
		t.Error("Expected IsSuccess to be true")
	}
	if result.IsError() {
		t.Error("Expected IsError to be false")
	}
}

func TestPostJSON_ResultMetadata(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Created", Email: "created@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := PostJSON[TestAPIUser](server.URL, map[string]string{"name": "Test"})
	if err != nil {
		t.Fatalf("PostJSON failed: %v", err)
	}
	if result.StatusCode() != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", result.StatusCode())
	}
	if !result.IsSuccess() {
		t.Error("Expected IsSuccess to be true")
	}
}

func TestPutJSON_ResultMetadata(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Updated", Email: "updated@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := PutJSON[TestAPIUser](server.URL, map[string]string{"name": "Updated"})
	if err != nil {
		t.Fatalf("PutJSON failed: %v", err)
	}
	if result.Response() == nil {
		t.Fatal("Response is nil")
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestDeleteJSON_ResultMetadata(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Deleted", Email: "deleted@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := DeleteJSON[TestAPIUser](server.URL)
	if err != nil {
		t.Fatalf("DeleteJSON failed: %v", err)
	}
	if result.Response() == nil {
		t.Fatal("Response is nil")
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestPatchJSON_ResultMetadata(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Patched", Email: "patched@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := PatchJSON[TestAPIUser](server.URL, map[string]string{"name": "Patched"})
	if err != nil {
		t.Fatalf("PatchJSON failed: %v", err)
	}
	if result.Response() == nil {
		t.Fatal("Response is nil")
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestGetXML_ResultMetadata(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Test", Email: "test@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := GetXML[TestAPIUser](server.URL)
	if err != nil {
		t.Fatalf("GetXML failed: %v", err)
	}
	if result.Response() == nil {
		t.Fatal("Response is nil")
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestPostXML_ResultMetadata(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Created", Email: "created@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	result, err := PostXML[TestAPIUser](server.URL, TestAPIUser{Name: "Test"})
	if err != nil {
		t.Fatalf("PostXML failed: %v", err)
	}
	if result.Response() == nil {
		t.Fatal("Response is nil")
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestGetJSONWithContext(t *testing.T) {
	expected := TestAPIUser{ID: 1, Name: "Test", Email: "test@example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()
	ctx := context.Background()
	result, err := GetJSONWithContext[TestAPIUser](ctx, server.URL)
	if err != nil {
		t.Fatalf("GetJSONWithContext failed: %v", err)
	}
	if !reflect.DeepEqual(expected, result.Data()) {
		t.Errorf("Expected %+v, got %+v", expected, result.Data())
	}
}

func TestGetJSONWithContext_Error(t *testing.T) {
	ctx := context.Background()
	_, err := GetJSONWithContext[TestAPIUser](ctx, "://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestGetJSONWithContext_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	ctx := context.Background()
	_, err := GetJSONWithContext[TestAPIUser](ctx, server.URL)
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

// Error path tests for client package methods
func TestPostJSON_Error(t *testing.T) {
	_, err := PostJSON[TestAPIUser]("://invalid-url", nil)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPostJSON_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	_, err := PostJSON[TestAPIUser]("https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPostJSON_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	_, err := PostJSON[TestAPIUser](server.URL, map[string]string{"name": "test"})
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

func TestPutJSON_Error(t *testing.T) {
	_, err := PutJSON[TestAPIUser]("://invalid-url", nil)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPutJSON_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	_, err := PutJSON[TestAPIUser]("https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPutJSON_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	_, err := PutJSON[TestAPIUser](server.URL, map[string]string{"name": "test"})
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

func TestDeleteJSON_Error(t *testing.T) {
	_, err := DeleteJSON[TestAPIUser]("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestDeleteJSON_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	_, err := DeleteJSON[TestAPIUser](server.URL)
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

func TestPatchJSON_Error(t *testing.T) {
	_, err := PatchJSON[TestAPIUser]("://invalid-url", nil)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPatchJSON_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `json:"ch"`
	}
	_, err := PatchJSON[TestAPIUser]("https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPatchJSON_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	_, err := PatchJSON[TestAPIUser](server.URL, map[string]string{"name": "test"})
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

func TestGetXML_Error(t *testing.T) {
	_, err := GetXML[TestAPIUser]("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestGetXML_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid xml"))
	}))
	defer server.Close()
	_, err := GetXML[TestAPIUser](server.URL)
	if err == nil {
		t.Error("Expected error for invalid XML response")
	}
}

func TestPostXML_Error(t *testing.T) {
	_, err := PostXML[TestAPIUser]("://invalid-url", nil)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestPostXML_EncodeError(t *testing.T) {
	type BadType struct {
		Ch chan int `xml:"ch"`
	}
	_, err := PostXML[TestAPIUser]("https://example.com", BadType{Ch: make(chan int)})
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestPostXML_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid xml"))
	}))
	defer server.Close()
	_, err := PostXML[TestAPIUser](server.URL, TestAPIUser{Name: "test"})
	if err == nil {
		t.Error("Expected error for invalid XML response")
	}
}

// Tests for GetString and GetBytes convenience methods
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

func TestGetString_Error(t *testing.T) {
	_, err := GetString("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestGetBytes(t *testing.T) {
	expected := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(expected)
	}))
	defer server.Close()
	result, err := GetBytes(server.URL)
	if err != nil {
		t.Fatalf("GetBytes failed: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetBytes_Error(t *testing.T) {
	_, err := GetBytes("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

// Property test for Result[T] generic methods
func TestProperty_GenericMethodsReturnResult(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)
	properties.Property("All generic methods return Result with accessible response", prop.ForAll(
		func(id int, name string) bool {
			if name == "" {
				return true
			}
			expected := TestAPIUser{ID: id, Name: name, Email: "test@example.com"}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Request-ID", "test-123")
				json.NewEncoder(w).Encode(expected)
			}))
			defer server.Close()
			result, err := GetJSON[TestAPIUser](server.URL)
			if err != nil {
				return false
			}
			// Verify data is accessible
			if result.Data().ID != id || result.Data().Name != name {
				return false
			}
			// Verify response metadata is accessible
			if result.Response() == nil {
				return false
			}
			if result.StatusCode() != http.StatusOK {
				return false
			}
			if result.Headers().Get("X-Request-ID") != "test-123" {
				return false
			}
			if !result.IsSuccess() {
				return false
			}
			return true
		},
		gen.Int(),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}
