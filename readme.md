# Go Requests

[![Language](https://img.shields.io/badge/language-golang-brightgreen)](https://github.com/sunerpy/requests) [![Last Commit](https://img.shields.io/github/last-commit/sunerpy/requests)](https://github.com/sunerpy/requests) [![CI](https://github.com/sunerpy/requests/workflows/Go/badge.svg)](https://github.com/sunerpy/requests/actions) [![codecov](https://codecov.io/gh/sunerpy/requests/branch/main/graph/badge.svg)](https://codecov.io/gh/sunerpy/requests) [![Benchmark](https://github.com/sunerpy/requests/actions/workflows/benchmark.yml/badge.svg)](https://sunerpy.github.io/requests/dev/bench) [![Go Reference](https://pkg.go.dev/badge/github.com/sunerpy/requests.svg)](https://pkg.go.dev/github.com/sunerpy/requests) [![Go Report Card](https://goreportcard.com/badge/github.com/sunerpy/requests)](https://goreportcard.com/report/github.com/sunerpy/requests) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

English | [简体中文](readme-cn.md)

Go Requests is a modern, type-safe HTTP client library for Go, inspired by Python's Requests library. It leverages Go 1.18+ generics to provide compile-time type safety and a fluent API for building HTTP requests.

## Features

* **Type-Safe Generic Methods**: `GetJSON[T]`, `PostJSON[T]`, etc. return `Result[T]` for automatic response parsing
* **Fluent Request Builder**: Chain methods to construct complex requests with `NewGet`, `NewPost`, etc.
* **Unified Result Type**: `Result[T]` wraps both parsed data and response metadata
* **Middleware Support**: Add cross-cutting concerns like logging, authentication, retry
* **Session Management**: Connection pooling, cookie persistence, default headers
* **Pluggable Codecs**: JSON, XML built-in with custom codec support
* **Retry Mechanism**: Configurable retry policies with exponential backoff
* **Request/Response Hooks**: Observe lifecycle events for logging and metrics
* **File Upload**: Multipart uploads with progress tracking
* **HTTP/2 Support**: Automatic HTTP/2 with fallback to HTTP/1.1
* **Connection Pool Optimization**: Configurable pool sizes and idle timeouts
* **High Performance**: Object pooling, zero-copy optimizations, close to net/http performance

## Install

```bash
go get github.com/sunerpy/requests
```

## Quick Start

### Simple GET Request

```go
import "github.com/sunerpy/requests"

// Basic GET request
resp, err := requests.Get("https://api.github.com")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Status:", resp.StatusCode)
fmt.Println("Body:", resp.Text())

// GET with query parameters using options
resp, err := requests.Get("https://api.github.com/search/repos",
    requests.WithQuery("q", "golang"),
    requests.WithQuery("sort", "stars"),
)
```

### Generic JSON Response with Result[T]

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// GetJSON returns Result[T] which wraps both data and response metadata
result, err := requests.GetJSON[User]("https://api.example.com/user/1")
if err != nil {
    log.Fatal(err)
}

// Access parsed data
fmt.Printf("User: %+v\n", result.Data())

// Access response metadata
fmt.Printf("Status: %d\n", result.StatusCode())
fmt.Printf("Success: %v\n", result.IsSuccess())
fmt.Printf("Headers: %v\n", result.Headers())
```

### POST with JSON Body

```go
type CreateUser struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponse struct {
    ID int `json:"id"`
}

// PostJSON returns Result[T]
result, err := requests.PostJSON[UserResponse](
    "https://api.example.com/users",
    CreateUser{Name: "John", Email: "john@example.com"},
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created user ID: %d\n", result.Data().ID)
fmt.Printf("Status: %d\n", result.StatusCode())
```

### Using Request Builder

```go
// Use convenience constructors for common methods
req, err := requests.NewPost("https://api.example.com/users").
    WithHeader("X-Custom-Header", "value").
    WithQuery("version", "v2").
    WithJSON(map[string]string{"name": "John"}).
    WithTimeout(10 * time.Second).
    Build()

// Or use NewRequestBuilder for any method
req, err := requests.NewRequestBuilder(requests.MethodPatch, "https://api.example.com/users/1").
    WithJSON(map[string]string{"name": "Jane"}).
    Build()

// Execute with DoJSON for automatic JSON parsing
result, err := requests.DoJSON[UserResponse](
    requests.NewPost("https://api.example.com/users").
        WithJSON(CreateUser{Name: "John", Email: "john@example.com"}),
)
```

### Using Request Options

```go
// Add headers, query params, auth, etc. using variadic options
resp, err := requests.Get("https://api.example.com/data",
    requests.WithHeader("X-Custom-Header", "value"),
    requests.WithQuery("page", "1"),
    requests.WithBearerToken("my-token"),
    requests.WithTimeout(10 * time.Second),
)
```

## Advanced Usage

### Session with Defaults

```go
session := requests.NewSession().
    WithBaseURL("https://api.example.com").
    WithTimeout(30 * time.Second).
    WithHeader("Authorization", "Bearer token").
    WithHTTP2(true).
    WithMaxIdleConns(100)

defer session.Close()

// Build request using the new Builder API
req, _ := requests.NewGet("/users").Build()

// Execute with session
resp, err := session.Do(req)
```

### Middleware

```go
// Logging middleware
loggingMiddleware := requests.MiddlewareFunc(func(req *requests.Request, next requests.Handler) (*requests.Response, error) {
    log.Printf("Request: %s %s", req.Method, req.URL)
    resp, err := next(req)
    if resp != nil {
        log.Printf("Response: %d", resp.StatusCode)
    }
    return resp, err
})

chain := requests.NewMiddlewareChain(loggingMiddleware)
```

### Retry Policy

```go
policy := requests.ExponentialRetryPolicy(3, 100*time.Millisecond, 10*time.Second)

executor := requests.NewRetryExecutor(policy)

resp, err := executor.Execute(nil, func() (*requests.Response, error) {
    return requests.Get("https://api.example.com/data")
})
```

### Request/Response Hooks

```go
hooks := requests.NewHooks()

// Log all requests
hooks.OnRequest(func(req *requests.Request) {
    log.Printf("Sending: %s %s", req.Method, req.URL)
})

// Log all responses with duration
hooks.OnResponse(func(req *requests.Request, resp *requests.Response, duration time.Duration) {
    log.Printf("Received: %d in %v", resp.StatusCode, duration)
})

// Log errors
hooks.OnError(func(req *requests.Request, err error, duration time.Duration) {
    log.Printf("Error: %v", err)
})
```

### File Upload

```go
import "github.com/sunerpy/requests/url"

upload, err := url.NewFileUpload("file", "/path/to/file.pdf")
if err != nil {
    log.Fatal(err)
}

upload.WithProgress(func(uploaded, total int64) {
    fmt.Printf("Progress: %d/%d (%.1f%%)\n", uploaded, total, float64(uploaded)/float64(total)*100)
})
```

### Custom Codec

```go
import "github.com/sunerpy/requests/codec"

// Register custom codec
type YAMLCodec struct{}

func (c *YAMLCodec) Encode(v any) ([]byte, error) { /* ... */ }
func (c *YAMLCodec) Decode(data []byte, v any) error { /* ... */ }
func (c *YAMLCodec) ContentType() string { return "application/yaml" }

codec.Register("application/yaml", &YAMLCodec{})
```

## API Reference

### Package `requests`

| Function | Description |
| -------- | ----------- |
| `Get(url, opts...)` | Send GET request |
| `Post(url, body, opts...)` | Send POST request with body |
| `Put(url, body, opts...)` | Send PUT request with body |
| `Delete(url, opts...)` | Send DELETE request |
| `Patch(url, body, opts...)` | Send PATCH request with body |
| `Head(url, opts...)` | Send HEAD request |
| `Options(url, opts...)` | Send OPTIONS request |
| `GetJSON[T](url, opts...)` | GET with JSON parsing, returns `(Result[T], error)` |
| `PostJSON[T](url, data, opts...)` | POST JSON, returns `(Result[T], error)` |
| `PutJSON[T](url, data, opts...)` | PUT JSON, returns `(Result[T], error)` |
| `DeleteJSON[T](url, opts...)` | DELETE with JSON parsing, returns `(Result[T], error)` |
| `PatchJSON[T](url, data, opts...)` | PATCH JSON, returns `(Result[T], error)` |
| `GetString(url, opts...)` | GET and return body as string |
| `GetBytes(url, opts...)` | GET and return body as bytes |
| `DoJSON[T](builder)` | Execute builder and parse JSON response |
| `NewSession()` | Create new session |

### Request Builder

| Function | Description |
| -------- | ----------- |
| `NewGet(url)` | Create GET request builder |
| `NewPost(url)` | Create POST request builder |
| `NewPut(url)` | Create PUT request builder |
| `NewDeleteBuilder(url)` | Create DELETE request builder |
| `NewPatch(url)` | Create PATCH request builder |
| `NewRequestBuilder(method, url)` | Create request builder with any method |

### Result[T] Methods

| Method | Description |
| ------ | ----------- |
| `Data()` | Get parsed response data of type T |
| `Response()` | Get underlying Response object |
| `StatusCode()` | Get HTTP status code |
| `Headers()` | Get response headers |
| `IsSuccess()` | Check if status is 2xx |
| `IsError()` | Check if status is 4xx or 5xx |
| `Cookies()` | Get response cookies |

### Request Options

| Option | Description |
| ------ | ----------- |
| `WithHeader(key, value)` | Add a header |
| `WithHeaders(map)` | Add multiple headers |
| `WithQuery(key, value)` | Add a query parameter |
| `WithQueryParams(map)` | Add multiple query parameters |
| `WithBasicAuth(user, pass)` | Set basic authentication |
| `WithBearerToken(token)` | Set bearer token |
| `WithTimeout(duration)` | Set request timeout |
| `WithContentType(type)` | Set Content-Type header |
| `WithContext(ctx)` | Set request context |
| `WithAccept(type)` | Set Accept header |

## Migration Guide

### From v1.x to v2.x

1. **Result[T] Type**: Generic methods now return `(Result[T], error)` instead of `(T, *Response, error)`

   ```go
   // Old
   user, resp, err := requests.GetJSONWithResponse[User](url)
   
   // New
   result, err := requests.GetJSON[User](url)
   user := result.Data()
   statusCode := result.StatusCode()
   ```

2. **Request Builder**: Use `NewGet`, `NewPost`, etc. instead of `NewRequest(method, url, params, body)`

   ```go
   // Old
   req, err := requests.NewRequest(requests.MethodGet, url, params, nil)
   
   // New
   req, err := requests.NewGet(url).WithQuery("key", "value").Build()
   ```

3. **Removed Methods**: `*WithResponse` methods are removed, use `Result[T]` instead

4. **Single Import**: Import only `github.com/sunerpy/requests` for most use cases

## Contributing

Contributions are welcome! Please:

1. Fork this repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## Performance

The library is optimized for high performance with:

* **Object Pooling**: `RequestBuilder`, `RequestConfig`, `FastValues` pools to reduce allocations
* **Header Ownership Transfer**: Avoids unnecessary cloning operations
* **Response Buffer Pool**: Tiered buffer pools (4KB/32KB) for response reading
* **Close to net/http**: Only ~6 more allocations than raw net/http

```bash
BenchmarkBasicRequests/NetHTTP_Get     40μs, 69 allocs, 6.3KB
BenchmarkBasicRequests/Requests_Get    44μs, 75 allocs, 6.3KB  (only +6 allocs)
```

For performance-critical code, use `FastValues` for URL parameter building:

```go
import "github.com/sunerpy/requests/url"

// Lock-free FastValues for single-threaded scenarios
fv := url.AcquireFastValues()
defer url.ReleaseFastValues(fv)
fv.Add("page", "1")
fv.Add("limit", "10")
```

See [Performance Optimization Guide](docs/guides/performance-optimization.md) for details.

## Documentation

For detailed documentation, see:

* [Getting Started](docs/guides/getting-started.md)
* [API Reference](docs/guides/api-reference.md)
* [Performance Optimization](docs/guides/performance-optimization.md)
* [Migration Guide](docs/guides/migration-guide.md)
* [Examples](docs/examples/)

## Project Structure

```bash
requests/
├── internal/           # Internal packages (not for external use)
│   ├── client/         # HTTP client core implementation
│   ├── models/         # Internal data models
│   └── utils/          # Internal utilities
├── codec/              # Encoder/decoder implementations
├── url/                # URL utilities and file upload
├── docs/               # Documentation
│   ├── guides/         # User guides
│   └── examples/       # Code examples
├── test/               # Benchmark tests
├── example/            # Example application
├── session.go          # Session management
├── methods.go          # HTTP method functions
├── types.go            # Type exports
└── Makefile            # Build automation
```

## License

MIT License - see [LICENSE](LICENSE) for details.
