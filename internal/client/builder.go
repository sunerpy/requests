package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sunerpy/requests/codec"
	"github.com/sunerpy/requests/internal/models"
)

// queryParam represents a single query parameter
type queryParam struct {
	key   string
	value string
}

// builderPool is a pool of RequestBuilder objects for reuse
var builderPool = sync.Pool{
	New: func() any {
		return &RequestBuilder{
			ctx: context.Background(),
		}
	},
}

// AcquireBuilder gets a RequestBuilder from the pool
func AcquireBuilder(method Method, rawURL string) *RequestBuilder {
	b := builderPool.Get().(*RequestBuilder)
	b.method = method
	b.rawURL = rawURL
	b.ctx = context.Background()
	return b
}

// ReleaseBuilder returns a RequestBuilder to the pool
func ReleaseBuilder(b *RequestBuilder) {
	if b == nil {
		return
	}
	b.Reset()
	builderPool.Put(b)
}

// Reset clears the builder for reuse
func (b *RequestBuilder) Reset() {
	b.method = ""
	b.rawURL = ""
	b.headers = nil
	b.query = nil
	b.queryParams = b.queryParams[:0] // Reset slice but keep capacity
	b.body = nil
	b.bodyBytes = nil
	b.bodyEncoder = nil
	b.bodyData = nil
	b.ctx = nil
	b.timeout = 0
	if b.files != nil {
		b.files = b.files[:0]
	}
	b.form = nil
	b.err = nil
}

// RequestBuilder provides a fluent interface for building HTTP requests.
type RequestBuilder struct {
	method      Method
	rawURL      string
	headers     http.Header
	query       url.Values
	queryParams []queryParam // Fast path for simple queries
	body        io.Reader
	bodyBytes   []byte
	bodyEncoder codec.Encoder
	bodyData    any
	ctx         context.Context
	timeout     time.Duration
	files       []FileUpload
	form        url.Values
	err         error
}

// NewRequest creates a new RequestBuilder with the specified method and URL.
// Uses lazy initialization for headers, query, and form to reduce allocations.
func NewRequest(method Method, rawURL string) *RequestBuilder {
	return &RequestBuilder{
		method: method,
		rawURL: rawURL,
		ctx:    context.Background(),
	}
}

// NewGetRequest creates a new GET request builder.
func NewGetRequest(rawURL string) *RequestBuilder {
	return NewRequest(MethodGet, rawURL)
}

// NewPostRequest creates a new POST request builder.
func NewPostRequest(rawURL string) *RequestBuilder {
	return NewRequest(MethodPost, rawURL)
}

// NewPutRequest creates a new PUT request builder.
func NewPutRequest(rawURL string) *RequestBuilder {
	return NewRequest(MethodPut, rawURL)
}

// NewDeleteRequest creates a new DELETE request builder.
func NewDeleteRequest(rawURL string) *RequestBuilder {
	return NewRequest(MethodDelete, rawURL)
}

// NewPatchRequest creates a new PATCH request builder.
func NewPatchRequest(rawURL string) *RequestBuilder {
	return NewRequest(MethodPatch, rawURL)
}

// NewGet creates a new GET request builder (short alias for NewGetRequest).
func NewGet(rawURL string) *RequestBuilder {
	return NewRequest(MethodGet, rawURL)
}

// NewPost creates a new POST request builder (short alias for NewPostRequest).
func NewPost(rawURL string) *RequestBuilder {
	return NewRequest(MethodPost, rawURL)
}

// NewPut creates a new PUT request builder (short alias for NewPutRequest).
func NewPut(rawURL string) *RequestBuilder {
	return NewRequest(MethodPut, rawURL)
}

// NewDelete creates a new DELETE request builder (short alias for NewDeleteRequest).
func NewDelete(rawURL string) *RequestBuilder {
	return NewRequest(MethodDelete, rawURL)
}

// NewPatch creates a new PATCH request builder (short alias for NewPatchRequest).
func NewPatch(rawURL string) *RequestBuilder {
	return NewRequest(MethodPatch, rawURL)
}

// WithURL sets the request URL.
func (b *RequestBuilder) WithURL(rawURL string) *RequestBuilder {
	b.rawURL = rawURL
	return b
}

// WithMethod sets the HTTP method.
func (b *RequestBuilder) WithMethod(method Method) *RequestBuilder {
	b.method = method
	return b
}

// WithHeader adds a single header to the request.
func (b *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(http.Header, 4)
	}
	b.headers.Set(key, value)
	return b
}

// WithHeaders adds multiple headers to the request.
func (b *RequestBuilder) WithHeaders(headers map[string]string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(http.Header, len(headers))
	}
	for k, v := range headers {
		b.headers.Set(k, v)
	}
	return b
}

// WithQuery adds a single query parameter.
// Uses fast path with slice storage to avoid map allocations.
func (b *RequestBuilder) WithQuery(key, value string) *RequestBuilder {
	// Use fast path with slice for simple cases
	b.queryParams = append(b.queryParams, queryParam{key: key, value: value})
	return b
}

// WithQueryParams adds multiple query parameters.
func (b *RequestBuilder) WithQueryParams(params map[string]string) *RequestBuilder {
	for k, v := range params {
		b.queryParams = append(b.queryParams, queryParam{key: k, value: v})
	}
	return b
}

// WithJSON sets the request body as JSON.
func (b *RequestBuilder) WithJSON(data any) *RequestBuilder {
	jsonData, err := codec.Marshal(data)
	if err != nil {
		b.err = &EncodeError{ContentType: "application/json", Err: err}
		return b
	}
	b.bodyBytes = jsonData
	b.body = bytes.NewReader(jsonData)
	b.bodyData = data
	if b.headers == nil {
		b.headers = make(http.Header, 2)
	}
	b.headers.Set("Content-Type", "application/json")
	return b
}

// WithXML sets the request body as XML.
func (b *RequestBuilder) WithXML(data any) *RequestBuilder {
	xmlCodec := codec.NewXMLCodec()
	xmlData, err := xmlCodec.Encode(data)
	if err != nil {
		b.err = &EncodeError{ContentType: "application/xml", Err: err}
		return b
	}
	b.bodyBytes = xmlData
	b.body = bytes.NewReader(xmlData)
	b.bodyData = data
	if b.headers == nil {
		b.headers = make(http.Header, 2)
	}
	b.headers.Set("Content-Type", "application/xml")
	return b
}

// WithForm sets the request body as form data.
func (b *RequestBuilder) WithForm(data map[string]string) *RequestBuilder {
	if b.form == nil {
		b.form = make(url.Values, len(data))
	}
	for k, v := range data {
		b.form.Set(k, v)
	}
	formData := b.form.Encode()
	b.bodyBytes = []byte(formData)
	b.body = strings.NewReader(formData)
	if b.headers == nil {
		b.headers = make(http.Header, 2)
	}
	b.headers.Set("Content-Type", "application/x-www-form-urlencoded")
	return b
}

// WithFormValues sets the request body from url.Values.
func (b *RequestBuilder) WithFormValues(values url.Values) *RequestBuilder {
	if b.form == nil {
		b.form = make(url.Values, len(values))
	}
	for k, vs := range values {
		for _, v := range vs {
			b.form.Add(k, v)
		}
	}
	formData := b.form.Encode()
	b.bodyBytes = []byte(formData)
	b.body = strings.NewReader(formData)
	if b.headers == nil {
		b.headers = make(http.Header, 2)
	}
	b.headers.Set("Content-Type", "application/x-www-form-urlencoded")
	return b
}

// WithBody sets a raw body reader.
func (b *RequestBuilder) WithBody(body io.Reader) *RequestBuilder {
	b.body = body
	return b
}

// WithBodyBytes sets the body from bytes.
func (b *RequestBuilder) WithBodyBytes(data []byte) *RequestBuilder {
	b.bodyBytes = data
	b.body = bytes.NewReader(data)
	return b
}

// WithFile adds a file for multipart upload.
func (b *RequestBuilder) WithFile(fieldName, fileName string, reader io.Reader) *RequestBuilder {
	b.files = append(b.files, FileUpload{
		FieldName: fieldName,
		FileName:  fileName,
		Reader:    reader,
	})
	return b
}

// WithContext sets the request context.
func (b *RequestBuilder) WithContext(ctx context.Context) *RequestBuilder {
	b.ctx = ctx
	return b
}

// WithTimeout sets the request timeout.
func (b *RequestBuilder) WithTimeout(d time.Duration) *RequestBuilder {
	b.timeout = d
	return b
}

// WithBasicAuth adds basic authentication header.
func (b *RequestBuilder) WithBasicAuth(username, password string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(http.Header, 2)
	}
	b.headers.Set("Authorization", "Basic "+basicAuth(username, password))
	return b
}

// WithBearerToken adds bearer token authentication header.
func (b *RequestBuilder) WithBearerToken(token string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(http.Header, 2)
	}
	b.headers.Set("Authorization", "Bearer "+token)
	return b
}

// WithEncoder sets a custom encoder for the body data.
func (b *RequestBuilder) WithEncoder(encoder codec.Encoder, data any) *RequestBuilder {
	encoded, err := encoder.Encode(data)
	if err != nil {
		b.err = &EncodeError{ContentType: encoder.ContentType(), Err: err}
		return b
	}
	b.bodyBytes = encoded
	b.body = bytes.NewReader(encoded)
	b.bodyData = data
	b.bodyEncoder = encoder
	if b.headers == nil {
		b.headers = make(http.Header, 2)
	}
	b.headers.Set("Content-Type", encoder.ContentType())
	return b
}

// Build constructs the final Request from the builder.
// Note: After calling Build(), the builder should not be reused unless Reset() is called.
// The headers are transferred by ownership to avoid allocation.
func (b *RequestBuilder) Build() (*Request, error) {
	// Validate builder state
	if err := b.validate(); err != nil {
		return nil, err
	}
	// Parse URL and merge query parameters
	parsedURL, err := b.parseAndMergeURL()
	if err != nil {
		return nil, err
	}
	// Apply timeout to context if set
	ctx := b.applyTimeout()
	// Transfer ownership of headers instead of cloning (optimization)
	// After this, b.headers will be nil and builder should not be reused
	var headers http.Header
	if b.headers != nil {
		headers = b.headers
		b.headers = nil // Transfer ownership
	} else {
		headers = make(http.Header)
	}
	// Transfer ownership of files slice
	files := b.files
	b.files = nil
	// Transfer ownership of form
	form := b.form
	b.form = nil
	return &Request{
		Method:    b.method,
		URL:       parsedURL,
		Headers:   headers,
		Body:      b.body,
		Context:   ctx,
		bodyBytes: b.bodyBytes,
		files:     files,
		form:      form,
	}, nil
}

// BuildCopy constructs a Request while preserving the builder state.
// Use this if you need to reuse the builder after building.
func (b *RequestBuilder) BuildCopy() (*Request, error) {
	// Validate builder state
	if err := b.validate(); err != nil {
		return nil, err
	}
	// Parse URL and merge query parameters
	parsedURL, err := b.parseAndMergeURL()
	if err != nil {
		return nil, err
	}
	// Apply timeout to context if set
	ctx := b.applyTimeout()
	// Clone headers (preserves builder state)
	headers := b.cloneHeaders()
	// Copy files slice
	files := b.cloneFiles()
	// Copy form
	form := cloneForm(b.form)
	return &Request{
		Method:    b.method,
		URL:       parsedURL,
		Headers:   headers,
		Body:      b.body,
		Context:   ctx,
		bodyBytes: b.bodyBytes,
		files:     files,
		form:      form,
	}, nil
}

// Do builds and executes the request, returning the response.
func (b *RequestBuilder) Do() (*models.Response, error) {
	req, err := b.Build()
	if err != nil {
		return nil, err
	}
	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(req.Context, req.Method.String(), req.URL.String(), req.Body)
	if err != nil {
		return nil, &RequestError{Op: "NewRequest", URL: req.URL.String(), Err: err}
	}
	// Apply headers
	for k, vs := range req.Headers {
		for _, v := range vs {
			httpReq.Header.Add(k, v)
		}
	}
	// Execute request
	httpResp, err := DefaultHTTPClient.Do(httpReq)
	if err != nil {
		return nil, &RequestError{Op: "Do", URL: req.URL.String(), Err: err}
	}
	return models.NewResponse(httpResp, req.URL.String())
}

// DoJSON builds and executes the request, parsing the JSON response into Result[T].
func DoJSON[T any](b *RequestBuilder) (Result[T], error) {
	var zero Result[T]
	resp, err := b.Do()
	if err != nil {
		return zero, err
	}
	data, err := models.JSON[T](resp)
	if err != nil {
		return zero, err
	}
	return NewResult(data, resp), nil
}

// DoXML builds and executes the request, parsing the XML response into Result[T].
func DoXML[T any](b *RequestBuilder) (Result[T], error) {
	var zero Result[T]
	resp, err := b.Do()
	if err != nil {
		return zero, err
	}
	data, err := models.XML[T](resp)
	if err != nil {
		return zero, err
	}
	return NewResult(data, resp), nil
}

// Clone creates a copy of the builder.
func (b *RequestBuilder) Clone() *RequestBuilder {
	clone := &RequestBuilder{
		method:      b.method,
		rawURL:      b.rawURL,
		bodyBytes:   b.bodyBytes,
		bodyData:    b.bodyData,
		bodyEncoder: b.bodyEncoder,
		ctx:         b.ctx,
		timeout:     b.timeout,
		err:         b.err,
	}
	// Clone headers if present
	if b.headers != nil {
		clone.headers = b.headers.Clone()
	}
	// Copy queryParams if present
	if len(b.queryParams) > 0 {
		clone.queryParams = make([]queryParam, len(b.queryParams))
		copy(clone.queryParams, b.queryParams)
	}
	// Copy query params if present (legacy)
	if b.query != nil {
		clone.query = make(url.Values, len(b.query))
		for k, vs := range b.query {
			clone.query[k] = append([]string(nil), vs...)
		}
	}
	// Copy form values if present
	if b.form != nil {
		clone.form = make(url.Values, len(b.form))
		for k, vs := range b.form {
			clone.form[k] = append([]string(nil), vs...)
		}
	}
	// Copy files
	if len(b.files) > 0 {
		clone.files = make([]FileUpload, len(b.files))
		copy(clone.files, b.files)
	}
	// Reset body reader if we have bytes
	if len(b.bodyBytes) > 0 {
		clone.body = bytes.NewReader(b.bodyBytes)
	}
	return clone
}

// GetMethod returns the current method.
func (b *RequestBuilder) GetMethod() Method {
	return b.method
}

// GetURL returns the current URL.
func (b *RequestBuilder) GetURL() string {
	return b.rawURL
}

// GetHeaders returns a copy of the current headers.
func (b *RequestBuilder) GetHeaders() http.Header {
	if b.headers == nil {
		return make(http.Header)
	}
	return b.headers.Clone()
}

// GetQuery returns a copy of the current query parameters.
func (b *RequestBuilder) GetQuery() url.Values {
	result := make(url.Values, len(b.queryParams)+len(b.query))
	// Add fast path params
	for _, p := range b.queryParams {
		result.Add(p.key, p.value)
	}
	// Add legacy params
	for k, vs := range b.query {
		for _, v := range vs {
			result.Add(k, v)
		}
	}
	return result
}

// GetBodyBytes returns the body bytes if available.
func (b *RequestBuilder) GetBodyBytes() []byte {
	return b.bodyBytes
}

// GetFiles returns the files to be uploaded.
func (b *RequestBuilder) GetFiles() []FileUpload {
	return b.files
}

// GetForm returns a copy of the form values.
func (b *RequestBuilder) GetForm() url.Values {
	if b.form == nil {
		return make(url.Values)
	}
	result := make(url.Values, len(b.form))
	for k, vs := range b.form {
		result[k] = append([]string(nil), vs...)
	}
	return result
}

// GetTimeout returns the timeout duration.
func (b *RequestBuilder) GetTimeout() time.Duration {
	return b.timeout
}

// validate checks builder state for required fields and returns an error if invalid.
func (b *RequestBuilder) validate() error {
	if b.err != nil {
		return b.err
	}
	if b.rawURL == "" {
		return &RequestError{Op: "Build", Err: ErrMissingURL}
	}
	if b.method == "" {
		return &RequestError{Op: "Build", Err: ErrMissingMethod}
	}
	if !b.method.IsValid() {
		return &RequestError{Op: "Build", Err: ErrInvalidMethod}
	}
	return nil
}

// parseAndMergeURL parses the raw URL and merges query parameters.
func (b *RequestBuilder) parseAndMergeURL() (*url.URL, error) {
	parsedURL, err := url.Parse(b.rawURL)
	if err != nil {
		return nil, &RequestError{Op: "ParseURL", URL: b.rawURL, Err: err}
	}
	// Fast path: use queryParams slice directly
	if len(b.queryParams) > 0 {
		var sb strings.Builder
		existingQuery := parsedURL.RawQuery
		if existingQuery != "" {
			sb.WriteString(existingQuery)
		}
		for _, p := range b.queryParams {
			if sb.Len() > 0 {
				sb.WriteByte('&')
			}
			sb.WriteString(url.QueryEscape(p.key))
			sb.WriteByte('=')
			sb.WriteString(url.QueryEscape(p.value))
		}
		parsedURL.RawQuery = sb.String()
	}
	// Legacy path: use url.Values if set
	if len(b.query) > 0 {
		existingQuery := parsedURL.Query()
		for k, vs := range b.query {
			for _, v := range vs {
				existingQuery.Add(k, v)
			}
		}
		parsedURL.RawQuery = existingQuery.Encode()
	}
	return parsedURL, nil
}

// applyTimeout applies timeout to context if set.
func (b *RequestBuilder) applyTimeout() context.Context {
	ctx := b.ctx
	if b.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(b.ctx, b.timeout)
		_ = cancel
	}
	return ctx
}

// cloneHeaders creates a copy of headers or returns an empty header.
func (b *RequestBuilder) cloneHeaders() http.Header {
	if b.headers != nil {
		return b.headers.Clone()
	}
	return make(http.Header)
}

// cloneFiles creates a copy of the files slice.
func (b *RequestBuilder) cloneFiles() []FileUpload {
	if len(b.files) == 0 {
		return nil
	}
	files := make([]FileUpload, len(b.files))
	copy(files, b.files)
	return files
}

// cloneForm creates a copy of form values.
func cloneForm(form url.Values) url.Values {
	if form == nil {
		return nil
	}
	result := make(url.Values, len(form))
	for k, vs := range form {
		result[k] = append([]string(nil), vs...)
	}
	return result
}

// basicAuth encodes username and password for basic auth.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64Encode([]byte(auth))
}

// base64Encode encodes bytes to base64 string.
func base64Encode(data []byte) string {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder
	result.Grow((len(data) + 2) / 3 * 4)
	for i := 0; i < len(data); i += 3 {
		var n uint32
		remaining := len(data) - i
		switch remaining {
		case 1:
			n = uint32(data[i]) << 16
			result.WriteByte(base64Chars[n>>18&0x3F])
			result.WriteByte(base64Chars[n>>12&0x3F])
			result.WriteString("==")
		case 2:
			n = uint32(data[i])<<16 | uint32(data[i+1])<<8
			result.WriteByte(base64Chars[n>>18&0x3F])
			result.WriteByte(base64Chars[n>>12&0x3F])
			result.WriteByte(base64Chars[n>>6&0x3F])
			result.WriteByte('=')
		default:
			n = uint32(data[i])<<16 | uint32(data[i+1])<<8 | uint32(data[i+2])
			result.WriteByte(base64Chars[n>>18&0x3F])
			result.WriteByte(base64Chars[n>>12&0x3F])
			result.WriteByte(base64Chars[n>>6&0x3F])
			result.WriteByte(base64Chars[n&0x3F])
		}
	}
	return result.String()
}
