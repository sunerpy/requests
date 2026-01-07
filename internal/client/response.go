package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/sunerpy/requests/codec"
	"github.com/sunerpy/requests/internal/models"
)

// Response represents an HTTP response with generic parsing capabilities.
type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Cookies    []*http.Cookie
	Proto      string
	body       []byte
	finalURL   string
	rawResp    *http.Response
}

// NewResponse creates a new Response from an http.Response.
func NewResponse(resp *http.Response, finalURL string) (*Response, error) {
	if resp == nil {
		return nil, &RequestError{Op: "NewResponse", Err: ErrNilResponse}
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RequestError{Op: "ReadBody", URL: finalURL, Err: err}
	}
	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Cookies:    resp.Cookies(),
		Proto:      resp.Proto,
		body:       data,
		finalURL:   finalURL,
		rawResp:    resp,
	}, nil
}

// Text returns the response body as a string.
func (r *Response) Text() string {
	return string(r.body)
}

// Bytes returns the response body as bytes.
func (r *Response) Bytes() []byte {
	return r.body
}

// Raw returns the underlying http.Response.
func (r *Response) Raw() *http.Response {
	return r.rawResp
}

// GetURL returns the final URL after redirects.
func (r *Response) GetURL() string {
	return r.finalURL
}

// IsSuccess returns true if status code is 2xx.
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError returns true if status code is 4xx or 5xx.
func (r *Response) IsError() bool {
	return r.StatusCode >= 400
}

// IsRedirect returns true if status code is 3xx.
func (r *Response) IsRedirect() bool {
	return r.StatusCode >= 300 && r.StatusCode < 400
}

// IsClientError returns true if status code is 4xx.
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if status code is 5xx.
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500
}

// ContentType returns the Content-Type header value.
func (r *Response) ContentType() string {
	return r.Headers.Get("Content-Type")
}

// DecodeJSON decodes the response body as JSON into the destination.
func (r *Response) DecodeJSON(dest any) error {
	if err := codec.Unmarshal(r.body, dest); err != nil {
		return &DecodeError{ContentType: "application/json", Err: err}
	}
	return nil
}

// DecodeXML decodes the response body as XML into the destination.
func (r *Response) DecodeXML(dest any) error {
	if err := xml.Unmarshal(r.body, dest); err != nil {
		return &DecodeError{ContentType: "application/xml", Err: err}
	}
	return nil
}

// Decode decodes the response body using the specified decoder.
func (r *Response) Decode(dest any, decoder codec.Decoder) error {
	return decoder.Decode(r.body, dest)
}

// DecodeAuto decodes the response body based on Content-Type header.
func (r *Response) DecodeAuto(dest any) error {
	contentType := r.ContentType()
	if contentType == "" {
		// Default to JSON if no content type
		return r.DecodeJSON(dest)
	}
	decoder, ok := codec.GetDecoder(contentType)
	if !ok {
		// Fall back to JSON for unknown content types
		return r.DecodeJSON(dest)
	}
	if err := decoder.Decode(r.body, dest); err != nil {
		return &DecodeError{ContentType: contentType, Err: err}
	}
	return nil
}

// JSON is a generic function that parses the response body as JSON.
func JSON[T any](r *Response) (T, error) {
	var result T
	if err := codec.Unmarshal(r.body, &result); err != nil {
		return result, &DecodeError{ContentType: "application/json", Err: err}
	}
	return result, nil
}

// XML is a generic function that parses the response body as XML.
func XML[T any](r *Response) (T, error) {
	var result T
	if err := xml.Unmarshal(r.body, &result); err != nil {
		return result, &DecodeError{ContentType: "application/xml", Err: err}
	}
	return result, nil
}

// Decode is a generic function that parses the response body using a decoder.
func Decode[T any](r *Response, decoder codec.Decoder) (T, error) {
	var result T
	if err := decoder.Decode(r.body, &result); err != nil {
		return result, err
	}
	return result, nil
}

// DecodeAuto is a generic function that parses based on Content-Type.
func DecodeAuto[T any](r *Response) (T, error) {
	contentType := r.ContentType()
	// Normalize content type
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = strings.TrimSpace(contentType[:idx])
	}
	// Try to get decoder from registry
	decoder, ok := codec.GetDecoder(contentType)
	if ok {
		var result T
		if err := decoder.Decode(r.body, &result); err != nil {
			return result, &DecodeError{ContentType: contentType, Err: err}
		}
		return result, nil
	}
	// Fall back to built-in decoders
	switch strings.ToLower(contentType) {
	case "application/xml", "text/xml":
		return XML[T](r)
	default:
		// Default to JSON
		return JSON[T](r)
	}
}

// MustJSON parses JSON and panics on error.
func MustJSON[T any](r *Response) T {
	result, err := JSON[T](r)
	if err != nil {
		panic(err)
	}
	return result
}

// MustXML parses XML and panics on error.
func MustXML[T any](r *Response) T {
	result, err := XML[T](r)
	if err != nil {
		panic(err)
	}
	return result
}

// CreateMockResponse creates a models.Response for testing purposes.
func CreateMockResponse(statusCode int, body []byte, headers http.Header) *models.Response {
	return models.CreateMockResponse(statusCode, body, headers)
}
