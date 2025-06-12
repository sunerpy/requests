package client

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// ResultTestUser is a test struct for JSON parsing in result tests
type ResultTestUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// TestResultAccessorConsistency tests Property 3: Result[T] Accessor Consistency
// Feature: api-design-optimization, Property 3: Result[T] Accessor Consistency
// Validates: Requirements 2.3, 2.4, 2.5, 2.6, 2.7
func TestResultAccessorConsistency(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	// Generate random status codes
	statusCodeGen := gen.IntRange(100, 599)
	// Generate random headers
	headerGen := gen.MapOf(gen.AlphaString(), gen.AlphaString())
	properties.Property("StatusCode() returns same value as Response().StatusCode", prop.ForAll(
		func(statusCode int) bool {
			resp := CreateMockResponse(statusCode, []byte("test"), nil)
			result := NewResult("test", resp)
			return result.StatusCode() == result.Response().StatusCode
		},
		statusCodeGen,
	))
	properties.Property("Headers() returns same value as Response().Headers", prop.ForAll(
		func(headers map[string]string) bool {
			httpHeaders := make(http.Header)
			for k, v := range headers {
				if k != "" {
					httpHeaders.Set(k, v)
				}
			}
			resp := CreateMockResponse(200, []byte("test"), httpHeaders)
			result := NewResult("test", resp)
			resultHeaders := result.Headers()
			respHeaders := result.Response().Headers
			if len(resultHeaders) != len(respHeaders) {
				return false
			}
			for k, v := range resultHeaders {
				if respHeaders.Get(k) != v[0] {
					return false
				}
			}
			return true
		},
		headerGen,
	))
	properties.Property("IsSuccess() returns true iff StatusCode is in [200, 300)", prop.ForAll(
		func(statusCode int) bool {
			resp := CreateMockResponse(statusCode, []byte("test"), nil)
			result := NewResult("test", resp)
			expected := statusCode >= 200 && statusCode < 300
			return result.IsSuccess() == expected
		},
		statusCodeGen,
	))
	properties.Property("IsError() returns true iff StatusCode >= 400", prop.ForAll(
		func(statusCode int) bool {
			resp := CreateMockResponse(statusCode, []byte("test"), nil)
			result := NewResult("test", resp)
			expected := statusCode >= 400
			return result.IsError() == expected
		},
		statusCodeGen,
	))
	properties.Property("IsClientError() returns true iff StatusCode is in [400, 500)", prop.ForAll(
		func(statusCode int) bool {
			resp := CreateMockResponse(statusCode, []byte("test"), nil)
			result := NewResult("test", resp)
			expected := statusCode >= 400 && statusCode < 500
			return result.IsClientError() == expected
		},
		statusCodeGen,
	))
	properties.Property("IsServerError() returns true iff StatusCode >= 500", prop.ForAll(
		func(statusCode int) bool {
			resp := CreateMockResponse(statusCode, []byte("test"), nil)
			result := NewResult("test", resp)
			expected := statusCode >= 500
			return result.IsServerError() == expected
		},
		statusCodeGen,
	))
	properties.TestingRun(t)
}

// TestResultDataRoundTrip tests Property 5: Result Data Round-Trip
// Feature: api-design-optimization, Property 5: Result Data Round-Trip
// Validates: Requirements 2.9
func TestResultDataRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	// Generate random user data
	idGen := gen.IntRange(1, 10000)
	nameGen := gen.AlphaString()
	emailGen := gen.AlphaString()
	properties.Property("JSON round-trip preserves data", prop.ForAll(
		func(id int, name, email string) bool {
			user := ResultTestUser{ID: id, Name: name, Email: email}
			// Serialize to JSON
			jsonData, err := json.Marshal(user)
			if err != nil {
				return false
			}
			// Create mock response with JSON data
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json")
			resp := CreateMockResponse(200, jsonData, headers)
			// Parse via Result
			parsed, err := JSON[ResultTestUser](resp)
			if err != nil {
				return false
			}
			result := NewResult(parsed, resp)
			// Verify round-trip
			return result.Data().ID == user.ID &&
				result.Data().Name == user.Name &&
				result.Data().Email == user.Email
		},
		idGen, nameGen, emailGen,
	))
	properties.TestingRun(t)
}

// TestResultNilResponse tests Result behavior with nil response
func TestResultNilResponse(t *testing.T) {
	result := NewResult("test", nil)
	if result.StatusCode() != 0 {
		t.Errorf("Expected StatusCode 0 for nil response, got %d", result.StatusCode())
	}
	if result.Headers() != nil {
		t.Errorf("Expected nil Headers for nil response")
	}
	if result.IsSuccess() {
		t.Errorf("Expected IsSuccess false for nil response")
	}
	if result.Cookies() != nil {
		t.Errorf("Expected nil Cookies for nil response")
	}
	if result.ContentType() != "" {
		t.Errorf("Expected empty ContentType for nil response")
	}
	if result.Text() != "" {
		t.Errorf("Expected empty Text for nil response")
	}
	if result.Bytes() != nil {
		t.Errorf("Expected nil Bytes for nil response")
	}
	if result.URL() != "" {
		t.Errorf("Expected empty URL for nil response")
	}
	// Data should still be accessible
	if result.Data() != "test" {
		t.Errorf("Expected Data to be 'test', got '%s'", result.Data())
	}
}

// TestResultDataAccess tests basic data access
func TestResultDataAccess(t *testing.T) {
	user := ResultTestUser{ID: 1, Name: "John", Email: "john@example.com"}
	resp := CreateMockResponse(200, []byte(`{"id":1,"name":"John","email":"john@example.com"}`), nil)
	result := NewResult(user, resp)
	if result.Data().ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.Data().ID)
	}
	if result.Data().Name != "John" {
		t.Errorf("Expected Name 'John', got '%s'", result.Data().Name)
	}
	if result.Data().Email != "john@example.com" {
		t.Errorf("Expected Email 'john@example.com', got '%s'", result.Data().Email)
	}
}
