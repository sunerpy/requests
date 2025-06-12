package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sunerpy/requests"
)

// User represents a user from the API
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// HTTPBinResponse represents a response from httpbin.org
type HTTPBinResponse struct {
	Args    map[string]string `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	URL     string            `json:"url"`
}

func main() {
	fmt.Println("=== HTTP Client Library Examples ===")
	fmt.Println()
	// Example 1: Basic GET request using package-level functions
	basicGetExample()
	// Example 2: Using Session with configuration
	sessionExample()
	// Example 3: Generic HTTP methods with automatic JSON parsing
	genericMethodsExample()
	// Example 4: Using RequestBuilder
	builderExample()
	// Example 5: Middleware usage
	middlewareExample()
	// Example 6: Retry mechanism
	retryExample()
	// Example 7: Hooks for observability
	hooksExample()
}

func basicGetExample() {
	fmt.Println("--- Example 1: Basic GET Request ---")
	resp, err := requests.Get("https://httpbin.org/get", requests.WithQuery("query", "golang"))
	if err != nil {
		log.Printf("GET Error: %v\n", err)
		return
	}
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Protocol: %s\n", resp.Proto)
	fmt.Println()
}

func sessionExample() {
	fmt.Println("--- Example 2: Session with Configuration ---")
	// Create a session with base URL and default headers
	session := requests.NewSession().
		WithBaseURL("https://httpbin.org").
		WithHeader("X-Custom-Header", "custom-value").
		WithBearerToken("my-secret-token").
		WithTimeout(30 * time.Second).
		WithHTTP2(true)
	defer session.Close()
	// Make request using session with the new Builder API
	req, _ := requests.NewGet("/headers").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("Session Error: %v\n", err)
		return
	}
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", resp.Text()[:200])
	fmt.Println()
}

func genericMethodsExample() {
	fmt.Println("--- Example 3: Generic HTTP Methods ---")
	// GET with automatic JSON parsing - returns Result[T] which wraps both data and response
	result, err := requests.GetJSON[HTTPBinResponse](
		"https://httpbin.org/get",
		requests.WithQuery("name", "John"),
		requests.WithHeader("Accept", "application/json"),
	)
	if err != nil {
		log.Printf("GetJSON Error: %v\n", err)
		return
	}
	// Access parsed data via Data() method
	fmt.Printf("Origin: %s\n", result.Data().Origin)
	fmt.Printf("Args: %v\n", result.Data().Args)
	// Access response metadata via Result methods
	fmt.Printf("Status: %d, IsSuccess: %v\n", result.StatusCode(), result.IsSuccess())
	// POST with JSON body and automatic response parsing
	postData := map[string]string{"name": "John", "email": "john@example.com"}
	postResult, err := requests.PostJSON[HTTPBinResponse](
		"https://httpbin.org/post",
		postData,
		requests.WithContentType("application/json"),
	)
	if err != nil {
		log.Printf("PostJSON Error: %v\n", err)
		return
	}
	fmt.Printf("POST URL: %s\n", postResult.Data().URL)
	fmt.Println()
}

func builderExample() {
	fmt.Println("--- Example 4: RequestBuilder ---")
	// Build a complex request using the builder pattern
	req, err := requests.NewRequestBuilder(requests.MethodPost, "https://httpbin.org/post").
		WithHeader("Content-Type", "application/json").
		WithHeader("Accept", "application/json").
		WithQuery("version", "v1").
		WithBearerToken("my-token").
		WithJSON(map[string]any{
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   30,
		}).
		WithTimeout(10 * time.Second).
		Build()
	if err != nil {
		log.Printf("Builder Error: %v\n", err)
		return
	}
	fmt.Printf("Built Request: %s %s\n", req.Method, req.URL)
	fmt.Printf("Headers: %v\n", req.Headers)
	fmt.Println()
}

func middlewareExample() {
	fmt.Println("--- Example 5: Middleware ---")
	// Create a middleware chain
	chain := requests.NewMiddlewareChain()
	// Add logging middleware
	chain.UseFunc(func(req *requests.Request, next requests.Handler) (*requests.Response, error) {
		fmt.Printf("  [Middleware] Request: %s %s\n", req.Method, req.URL)
		resp, err := next(req)
		if resp != nil {
			fmt.Printf("  [Middleware] Response: %d\n", resp.StatusCode)
		}
		return resp, err
	})
	// Add header middleware
	chain.Use(requests.HeaderMiddleware(map[string]string{
		"X-Request-ID": "12345",
	}))
	// Execute with middleware using DoJSON for type-safe response
	builder := requests.NewGet("https://httpbin.org/get")
	result, err := requests.DoJSON[HTTPBinResponse](builder)
	if err != nil {
		log.Printf("Middleware Error: %v\n", err)
		return
	}
	fmt.Printf("Final Status: %d\n", result.StatusCode())
	fmt.Println()
}

func retryExample() {
	fmt.Println("--- Example 6: Retry Mechanism ---")
	// Create a retry policy
	policy := requests.ExponentialRetryPolicy(3, 100*time.Millisecond, 1*time.Second)
	// Create retry executor
	executor := requests.NewRetryExecutor(policy)
	var attempts int
	// Use DoJSON with retry for type-safe response
	result, err := executor.Execute(context.TODO(), func() (*requests.Response, error) {
		attempts++
		fmt.Printf("  Attempt %d\n", attempts)
		// Use the builder pattern with Do method
		builder := requests.NewGet("https://httpbin.org/get")
		return builder.Do()
	})
	if err != nil {
		log.Printf("Retry Error: %v\n", err)
		return
	}
	fmt.Printf("Success after %d attempt(s), Status: %d\n", attempts, result.StatusCode)
	fmt.Println()
}

func hooksExample() {
	fmt.Println("--- Example 7: Hooks for Observability ---")
	// Create hooks
	hooks := requests.NewHooks()
	// Add request hook
	hooks.OnRequest(func(req *requests.Request) {
		fmt.Printf("  [Hook] Sending request to: %s\n", req.URL)
	})
	// Add response hook
	hooks.OnResponse(func(req *requests.Request, resp *requests.Response, duration time.Duration) {
		fmt.Printf("  [Hook] Received response: %d in %v\n", resp.StatusCode, duration)
	})
	// Add error hook
	hooks.OnError(func(req *requests.Request, err error, duration time.Duration) {
		fmt.Printf("  [Hook] Error: %v in %v\n", err, duration)
	})
	// Use hooks middleware with DoJSON for type-safe response
	fmt.Println("  Making request with hooks...")
	result, err := requests.DoJSON[HTTPBinResponse](requests.NewGet("https://httpbin.org/get"))
	if err != nil {
		log.Printf("Hooks Error: %v\n", err)
		return
	}
	fmt.Printf("Final Status: %d\n", result.StatusCode())
	fmt.Println()
	// Using MetricsHook
	fmt.Println("--- Metrics Hook Example ---")
	metrics := requests.NewMetricsHook()
	metricsHooks := requests.NewHooks().
		OnRequest(metrics.RequestHook()).
		OnResponse(metrics.ResponseHook()).
		OnError(metrics.ErrorHook())
	// Note: metricsHooks would be used with middleware chain
	_ = metricsHooks
	// Make a few requests using DoJSON
	for i := 0; i < 3; i++ {
		requests.DoJSON[HTTPBinResponse](requests.NewGet("https://httpbin.org/get"))
	}
	totalRequests, responses, errors, avgDuration := metrics.Stats()
	fmt.Printf("Metrics: %d requests, %d responses, %d errors, avg duration: %v\n",
		totalRequests, responses, errors, avgDuration)
}
