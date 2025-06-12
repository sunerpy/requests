package client

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// TestUser is a test struct for JSON/XML parsing tests.
type TestUser struct {
	ID    int    `json:"id" xml:"id"`
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
	Age   int    `json:"age" xml:"age"`
}

// TestNestedData is a test struct with nested fields.
type TestNestedData struct {
	Title  string   `json:"title" xml:"title"`
	Tags   []string `json:"tags" xml:"tags>tag"`
	Author TestUser `json:"author" xml:"author"`
}

// Feature: http-client-refactor
// Property 1: Response JSON/XML Parsing Round-Trip
// For any valid Go struct that can be serialized to JSON or XML, creating a Response
// with that serialized content and parsing it back using JSON[T]() or XML[T]()
// SHALL produce an equivalent struct.
// Validates: Requirements 1.1, 1.3, 1.4, 1.5
func TestProperty1_JSONRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("JSON round-trip preserves struct values", prop.ForAll(
		func(id int, name string, email string, age int) bool {
			// Skip empty strings that might cause issues
			if name == "" || email == "" {
				return true
			}
			original := TestUser{
				ID:    id,
				Name:  name,
				Email: email,
				Age:   age,
			}
			// Serialize to JSON
			data, err := json.Marshal(original)
			if err != nil {
				return false
			}
			// Create mock response
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json")
			resp := CreateMockResponse(200, data, headers)
			// Parse back using generic function
			parsed, err := JSON[TestUser](resp)
			if err != nil {
				return false
			}
			// Verify equality
			return reflect.DeepEqual(original, parsed)
		},
		gen.Int(),
		gen.AlphaString(),
		gen.AlphaString(),
		gen.IntRange(0, 150),
	))
	properties.TestingRun(t)
}

func TestProperty1_XMLRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("XML round-trip preserves struct values", prop.ForAll(
		func(id int, name string, email string, age int) bool {
			// Skip empty strings and strings with special XML characters
			if name == "" || email == "" {
				return true
			}
			original := TestUser{
				ID:    id,
				Name:  name,
				Email: email,
				Age:   age,
			}
			// Serialize to XML
			data, err := xml.Marshal(original)
			if err != nil {
				return false
			}
			// Create mock response
			headers := make(http.Header)
			headers.Set("Content-Type", "application/xml")
			resp := CreateMockResponse(200, data, headers)
			// Parse back using generic function
			parsed, err := XML[TestUser](resp)
			if err != nil {
				return false
			}
			// Verify equality
			return reflect.DeepEqual(original, parsed)
		},
		gen.Int(),
		gen.AlphaString(),
		gen.AlphaString(),
		gen.IntRange(0, 150),
	))
	properties.TestingRun(t)
}

func TestProperty1_DecodeAutoJSON(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("DecodeAuto correctly parses JSON based on Content-Type", prop.ForAll(
		func(id int, name string) bool {
			if name == "" {
				return true
			}
			original := TestUser{ID: id, Name: name, Email: "test@example.com", Age: 25}
			data, err := json.Marshal(original)
			if err != nil {
				return false
			}
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json; charset=utf-8")
			resp := CreateMockResponse(200, data, headers)
			parsed, err := DecodeAuto[TestUser](resp)
			if err != nil {
				return false
			}
			return reflect.DeepEqual(original, parsed)
		},
		gen.Int(),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

func TestProperty1_DecodeAutoXML(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("DecodeAuto correctly parses XML based on Content-Type", prop.ForAll(
		func(id int, name string) bool {
			if name == "" {
				return true
			}
			original := TestUser{ID: id, Name: name, Email: "test@example.com", Age: 25}
			data, err := xml.Marshal(original)
			if err != nil {
				return false
			}
			headers := make(http.Header)
			headers.Set("Content-Type", "application/xml")
			resp := CreateMockResponse(200, data, headers)
			parsed, err := DecodeAuto[TestUser](resp)
			if err != nil {
				return false
			}
			return reflect.DeepEqual(original, parsed)
		},
		gen.Int(),
		gen.AlphaString(),
	))
	properties.TestingRun(t)
}

// Feature: http-client-refactor
// Property 2: Invalid Response Parsing Returns Error
// For any Response containing invalid JSON/XML content, calling JSON[T]() or XML[T]()
// SHALL return a non-nil error with descriptive context.
// Validates: Requirements 1.2
func TestProperty2_InvalidJSONReturnsError(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Invalid JSON returns DecodeError", prop.ForAll(
		func(invalidData string) bool {
			// Ensure the data is actually invalid JSON
			var test TestUser
			if json.Unmarshal([]byte(invalidData), &test) == nil {
				// Skip if it happens to be valid JSON
				return true
			}
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json")
			resp := CreateMockResponse(200, []byte(invalidData), headers)
			_, err := JSON[TestUser](resp)
			// Must return an error
			if err == nil {
				return false
			}
			// Error should be a DecodeError
			decodeErr, ok := err.(*DecodeError)
			if !ok {
				return false
			}
			// DecodeError should have content type set
			return decodeErr.ContentType == "application/json"
		},
		gen.AnyString().SuchThat(func(s string) bool {
			// Generate strings that are likely invalid JSON
			return len(s) > 0 && s[0] != '{' && s[0] != '['
		}),
	))
	properties.TestingRun(t)
}

func TestProperty2_InvalidXMLReturnsError(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Invalid XML returns DecodeError", prop.ForAll(
		func(invalidData string) bool {
			// Ensure the data is actually invalid XML
			var test TestUser
			if xml.Unmarshal([]byte(invalidData), &test) == nil {
				// Skip if it happens to be valid XML
				return true
			}
			headers := make(http.Header)
			headers.Set("Content-Type", "application/xml")
			resp := CreateMockResponse(200, []byte(invalidData), headers)
			_, err := XML[TestUser](resp)
			// Must return an error
			if err == nil {
				return false
			}
			// Error should be a DecodeError
			decodeErr, ok := err.(*DecodeError)
			if !ok {
				return false
			}
			// DecodeError should have content type set
			return decodeErr.ContentType == "application/xml"
		},
		gen.AnyString().SuchThat(func(s string) bool {
			// Generate strings that are likely invalid XML
			return len(s) > 0 && s[0] != '<'
		}),
	))
	properties.TestingRun(t)
}

func TestProperty2_MalformedJSONStructure(t *testing.T) {
	testCases := []struct {
		name string
		data string
	}{
		{"empty string", ""},
		{"plain text", "hello world"},
		{"incomplete object", `{"id": 1, "name":`},
		{"missing quotes", `{id: 1}`},
		{"trailing comma", `{"id": 1,}`},
		{"wrong type", `{"id": "not a number"}`},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json")
			resp := CreateMockResponse(200, []byte(tc.data), headers)
			_, err := JSON[TestUser](resp)
			if err == nil {
				t.Errorf("Expected error for invalid JSON: %s", tc.data)
			}
			decodeErr, ok := err.(*DecodeError)
			if !ok {
				t.Errorf("Expected DecodeError, got %T", err)
			}
			if decodeErr != nil && decodeErr.ContentType != "application/json" {
				t.Errorf("Expected content type 'application/json', got '%s'", decodeErr.ContentType)
			}
		})
	}
}

func TestProperty2_MalformedXMLStructure(t *testing.T) {
	testCases := []struct {
		name string
		data string
	}{
		{"empty string", ""},
		{"plain text", "hello world"},
		{"unclosed tag", "<TestUser><id>1</id>"},
		{"mismatched tags", "<TestUser><id>1</name></TestUser>"},
		{"invalid characters", "<TestUser><id>&invalid;</id></TestUser>"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := make(http.Header)
			headers.Set("Content-Type", "application/xml")
			resp := CreateMockResponse(200, []byte(tc.data), headers)
			_, err := XML[TestUser](resp)
			if err == nil {
				t.Errorf("Expected error for invalid XML: %s", tc.data)
			}
			decodeErr, ok := err.(*DecodeError)
			if !ok {
				t.Errorf("Expected DecodeError, got %T", err)
			}
			if decodeErr != nil && decodeErr.ContentType != "application/xml" {
				t.Errorf("Expected content type 'application/xml', got '%s'", decodeErr.ContentType)
			}
		})
	}
}

// Unit tests for Response helper methods
func TestResponse_IsSuccess(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   bool
	}{
		{200, true},
		{201, true},
		{204, true},
		{299, true},
		{300, false},
		{400, false},
		{500, false},
	}
	for _, tc := range testCases {
		resp := CreateMockResponse(tc.statusCode, nil, nil)
		if resp.IsSuccess() != tc.expected {
			t.Errorf("IsSuccess() for status %d: expected %v, got %v", tc.statusCode, tc.expected, resp.IsSuccess())
		}
	}
}

func TestResponse_IsError(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{300, false},
		{399, false},
		{400, true},
		{404, true},
		{500, true},
		{503, true},
	}
	for _, tc := range testCases {
		resp := CreateMockResponse(tc.statusCode, nil, nil)
		if resp.IsError() != tc.expected {
			t.Errorf("IsError() for status %d: expected %v, got %v", tc.statusCode, tc.expected, resp.IsError())
		}
	}
}

func TestResponse_IsRedirect(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{300, true},
		{301, true},
		{302, true},
		{307, true},
		{399, true},
		{400, false},
	}
	for _, tc := range testCases {
		resp := CreateMockResponse(tc.statusCode, nil, nil)
		if resp.IsRedirect() != tc.expected {
			t.Errorf("IsRedirect() for status %d: expected %v, got %v", tc.statusCode, tc.expected, resp.IsRedirect())
		}
	}
}

func TestResponse_IsClientError(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   bool
	}{
		{399, false},
		{400, true},
		{404, true},
		{499, true},
		{500, false},
	}
	for _, tc := range testCases {
		resp := CreateMockResponse(tc.statusCode, nil, nil)
		if resp.IsClientError() != tc.expected {
			t.Errorf("IsClientError() for status %d: expected %v, got %v", tc.statusCode, tc.expected, resp.IsClientError())
		}
	}
}

func TestResponse_IsServerError(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   bool
	}{
		{499, false},
		{500, true},
		{502, true},
		{503, true},
		{599, true},
	}
	for _, tc := range testCases {
		resp := CreateMockResponse(tc.statusCode, nil, nil)
		if resp.IsServerError() != tc.expected {
			t.Errorf("IsServerError() for status %d: expected %v, got %v", tc.statusCode, tc.expected, resp.IsServerError())
		}
	}
}

func TestResponse_ContentType(t *testing.T) {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json; charset=utf-8")
	resp := CreateMockResponse(200, nil, headers)
	if resp.ContentType() != "application/json; charset=utf-8" {
		t.Errorf("Expected 'application/json; charset=utf-8', got '%s'", resp.ContentType())
	}
}

func TestResponse_TextAndBytes(t *testing.T) {
	body := []byte("Hello, World!")
	resp := CreateMockResponse(200, body, nil)
	if resp.Text() != "Hello, World!" {
		t.Errorf("Text() expected 'Hello, World!', got '%s'", resp.Text())
	}
	if string(resp.Bytes()) != "Hello, World!" {
		t.Errorf("Bytes() expected 'Hello, World!', got '%s'", string(resp.Bytes()))
	}
}

func TestResponse_DecodeJSON(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)
	resp := CreateMockResponse(200, data, nil)
	var parsed TestUser
	err := resp.DecodeJSON(&parsed)
	if err != nil {
		t.Fatalf("DecodeJSON failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("DecodeJSON result mismatch: expected %+v, got %+v", user, parsed)
	}
}

func TestResponse_DecodeXML(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := xml.Marshal(user)
	resp := CreateMockResponse(200, data, nil)
	var parsed TestUser
	err := resp.DecodeXML(&parsed)
	if err != nil {
		t.Fatalf("DecodeXML failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("DecodeXML result mismatch: expected %+v, got %+v", user, parsed)
	}
}

func TestMustJSON_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustJSON should panic on invalid JSON")
		}
	}()
	resp := CreateMockResponse(200, []byte("invalid json"), nil)
	_ = MustJSON[TestUser](resp)
}

func TestMustXML_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustXML should panic on invalid XML")
		}
	}()
	resp := CreateMockResponse(200, []byte("invalid xml"), nil)
	_ = MustXML[TestUser](resp)
}

// Additional tests for coverage
func TestResponse_Raw(t *testing.T) {
	resp := CreateMockResponse(200, []byte("test"), nil)
	// CreateMockResponse sets rawResp to nil
	if resp.Raw() != nil {
		t.Error("Expected nil Raw() for mock response")
	}
}

func TestResponse_GetURL(t *testing.T) {
	resp := CreateMockResponse(200, []byte("test"), nil)
	// CreateMockResponse sets finalURL to empty string
	if resp.GetURL() != "" {
		t.Errorf("Expected empty URL, got '%s'", resp.GetURL())
	}
}

func TestResponse_DecodeAuto_EmptyContentType(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)
	// No Content-Type header - should default to JSON
	resp := CreateMockResponse(200, data, nil)
	parsed, err := DecodeAuto[TestUser](resp)
	if err != nil {
		t.Fatalf("DecodeAuto failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("DecodeAuto result mismatch")
	}
}

func TestResponse_DecodeAuto_TextXML(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := xml.Marshal(user)
	headers := make(http.Header)
	headers.Set("Content-Type", "text/xml")
	resp := CreateMockResponse(200, data, headers)
	parsed, err := DecodeAuto[TestUser](resp)
	if err != nil {
		t.Fatalf("DecodeAuto failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("DecodeAuto result mismatch")
	}
}

func TestMustJSON_Success(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)
	resp := CreateMockResponse(200, data, nil)
	parsed := MustJSON[TestUser](resp)
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("MustJSON result mismatch")
	}
}

func TestMustXML_Success(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := xml.Marshal(user)
	resp := CreateMockResponse(200, data, nil)
	parsed := MustXML[TestUser](resp)
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("MustXML result mismatch")
	}
}

func TestResponse_DecodeJSON_Error(t *testing.T) {
	resp := CreateMockResponse(200, []byte("invalid json"), nil)
	var user TestUser
	err := resp.DecodeJSON(&user)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	decodeErr, ok := err.(*DecodeError)
	if !ok {
		t.Errorf("Expected DecodeError, got %T", err)
	}
	if decodeErr != nil && decodeErr.ContentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got '%s'", decodeErr.ContentType)
	}
}

func TestResponse_DecodeXML_Error(t *testing.T) {
	resp := CreateMockResponse(200, []byte("invalid xml"), nil)
	var user TestUser
	err := resp.DecodeXML(&user)
	if err == nil {
		t.Error("Expected error for invalid XML")
	}
	decodeErr, ok := err.(*DecodeError)
	if !ok {
		t.Errorf("Expected DecodeError, got %T", err)
	}
	if decodeErr != nil && decodeErr.ContentType != "application/xml" {
		t.Errorf("Expected content type 'application/xml', got '%s'", decodeErr.ContentType)
	}
}

func TestNewResponse_NilResponse(t *testing.T) {
	resp, err := NewResponse(nil, "https://example.com")
	if err == nil {
		t.Error("Expected error for nil response")
	}
	if resp != nil {
		t.Error("Expected nil response")
	}
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError, got %T", err)
	}
	if reqErr != nil && reqErr.Op != "NewResponse" {
		t.Errorf("Expected Op 'NewResponse', got '%s'", reqErr.Op)
	}
}

func TestNewResponse_Success(t *testing.T) {
	// Create a real http.Response
	httpResp := &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("test body")),
		Proto:      "HTTP/1.1",
	}
	httpResp.Header.Set("Content-Type", "text/plain")
	resp, err := NewResponse(httpResp, "https://example.com/test")
	if err != nil {
		t.Fatalf("NewResponse failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if resp.Text() != "test body" {
		t.Errorf("Expected 'test body', got '%s'", resp.Text())
	}
	if resp.GetURL() != "https://example.com/test" {
		t.Errorf("Expected URL 'https://example.com/test', got '%s'", resp.GetURL())
	}
	if resp.Raw() != httpResp {
		t.Error("Raw() should return original response")
	}
}

func TestResponse_DecodeAuto_UnknownContentType(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)
	headers := make(http.Header)
	headers.Set("Content-Type", "application/unknown")
	resp := CreateMockResponse(200, data, headers)
	// Should fall back to JSON
	var parsed TestUser
	err := resp.DecodeAuto(&parsed)
	if err != nil {
		t.Fatalf("DecodeAuto failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("DecodeAuto result mismatch")
	}
}

func TestResponse_DecodeAuto_DecodeError(t *testing.T) {
	// Invalid JSON data with JSON content type
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	resp := CreateMockResponse(200, []byte("invalid json"), headers)
	var parsed TestUser
	err := resp.DecodeAuto(&parsed)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	decodeErr, ok := err.(*DecodeError)
	if !ok {
		t.Errorf("Expected DecodeError, got %T", err)
	}
	if decodeErr != nil && decodeErr.ContentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got '%s'", decodeErr.ContentType)
	}
}

func TestDecodeAuto_Generic_Error(t *testing.T) {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	resp := CreateMockResponse(200, []byte("invalid json"), headers)
	_, err := DecodeAuto[TestUser](resp)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	decodeErr, ok := err.(*DecodeError)
	if !ok {
		t.Errorf("Expected DecodeError, got %T", err)
	}
	if decodeErr != nil && decodeErr.ContentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got '%s'", decodeErr.ContentType)
	}
}

func TestDecodeAuto_Generic_XMLContentType(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := xml.Marshal(user)
	headers := make(http.Header)
	headers.Set("Content-Type", "text/xml; charset=utf-8")
	resp := CreateMockResponse(200, data, headers)
	parsed, err := DecodeAuto[TestUser](resp)
	if err != nil {
		t.Fatalf("DecodeAuto failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("DecodeAuto result mismatch")
	}
}

func TestDecodeAuto_Generic_UnknownContentType(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)
	headers := make(http.Header)
	headers.Set("Content-Type", "application/unknown")
	resp := CreateMockResponse(200, data, headers)
	// Should fall back to JSON
	parsed, err := DecodeAuto[TestUser](resp)
	if err != nil {
		t.Fatalf("DecodeAuto failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("DecodeAuto result mismatch")
	}
}

func TestDecodeAuto_Generic_DecoderError(t *testing.T) {
	// Invalid JSON data with JSON content type (uses registered decoder)
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	resp := CreateMockResponse(200, []byte("invalid json"), headers)
	_, err := DecodeAuto[TestUser](resp)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	decodeErr, ok := err.(*DecodeError)
	if !ok {
		t.Errorf("Expected DecodeError, got %T", err)
	}
	if decodeErr != nil && decodeErr.ContentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got '%s'", decodeErr.ContentType)
	}
}

func TestResponse_Decode_WithDecoder(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)
	resp := CreateMockResponse(200, data, nil)
	// Use a mock decoder
	decoder := &mockDecoder{
		decodeFunc: func(data []byte, v any) error {
			return json.Unmarshal(data, v)
		},
	}
	var parsed TestUser
	err := resp.Decode(&parsed, decoder)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("Decode result mismatch")
	}
}

func TestResponse_DecodeAuto_Method(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	resp := CreateMockResponse(200, data, headers)
	var parsed TestUser
	err := resp.DecodeAuto(&parsed)
	if err != nil {
		t.Fatalf("DecodeAuto failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("DecodeAuto result mismatch")
	}
}

func TestResponse_DecodeAuto_Method_EmptyContentType(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)
	// No Content-Type header
	resp := CreateMockResponse(200, data, nil)
	var parsed TestUser
	err := resp.DecodeAuto(&parsed)
	if err != nil {
		t.Fatalf("DecodeAuto failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("DecodeAuto result mismatch")
	}
}

func TestDecode_Generic_WithDecoder(t *testing.T) {
	user := TestUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)
	resp := CreateMockResponse(200, data, nil)
	decoder := &mockDecoder{
		decodeFunc: func(data []byte, v any) error {
			return json.Unmarshal(data, v)
		},
	}
	parsed, err := Decode[TestUser](resp, decoder)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if !reflect.DeepEqual(user, parsed) {
		t.Errorf("Decode result mismatch")
	}
}

func TestDecode_Generic_Error(t *testing.T) {
	resp := CreateMockResponse(200, []byte("invalid"), nil)
	decoder := &mockDecoder{
		decodeFunc: func(data []byte, v any) error {
			return json.Unmarshal(data, v)
		},
	}
	_, err := Decode[TestUser](resp, decoder)
	if err == nil {
		t.Error("Expected error for invalid data")
	}
}

type mockDecoder struct {
	decodeFunc func(data []byte, v any) error
}

func (d *mockDecoder) Decode(data []byte, v any) error {
	return d.decodeFunc(data, v)
}

func (d *mockDecoder) ContentType() string {
	return "application/mock"
}

// errorReader is a reader that always returns an error
type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (r *errorReader) Close() error {
	return nil
}

func TestNewResponse_ReadBodyError(t *testing.T) {
	httpResp := &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Header:     make(http.Header),
		Body:       &errorReader{},
		Proto:      "HTTP/1.1",
	}
	resp, err := NewResponse(httpResp, "https://example.com")
	if err == nil {
		t.Error("Expected error for read body failure")
	}
	if resp != nil {
		t.Error("Expected nil response")
	}
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError, got %T", err)
	}
	if reqErr != nil && reqErr.Op != "ReadBody" {
		t.Errorf("Expected Op 'ReadBody', got '%s'", reqErr.Op)
	}
}
