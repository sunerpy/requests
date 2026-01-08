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
	// Example 8: Context with timeout and cancellation
	contextExample()
	// Example 9: Session with retry policy
	sessionRetryExample()
	// Example 10: Session with middleware
	sessionMiddlewareExample()
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

func contextExample() {
	fmt.Println("--- Example 8: Context with Timeout and Cancellation ---")

	// Example 8a: Context with timeout
	fmt.Println("  8a: Context with timeout")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session := requests.NewSession()
	defer session.Close()

	req, _ := requests.NewGet("https://httpbin.org/get").Build()

	// DoWithContext respects context timeout and cancellation
	resp, err := session.DoWithContext(ctx, req)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			fmt.Println("  Request timed out")
		case context.Canceled:
			fmt.Println("  Request was canceled")
		default:
			log.Printf("  Context Error: %v\n", err)
		}
		return
	}
	fmt.Printf("  Status: %d\n", resp.StatusCode)

	// Example 8b: Context cancellation
	fmt.Println("  8b: Context cancellation (simulated)")
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	// In a real scenario, you might cancel based on user action or other events
	go func() {
		time.Sleep(100 * time.Millisecond)
		// cancelFunc() // Uncomment to test cancellation
		_ = cancelFunc // Avoid unused variable warning
	}()

	req2, _ := requests.NewGet("https://httpbin.org/delay/1").Build()
	resp2, err := session.DoWithContext(cancelCtx, req2)
	if err != nil {
		fmt.Printf("  Canceled or error: %v\n", err)
	} else {
		fmt.Printf("  Completed with status: %d\n", resp2.StatusCode)
	}
	cancelFunc() // Clean up

	fmt.Println()
}

func sessionRetryExample() {
	fmt.Println("--- Example 9: Session with Retry Policy ---")

	// Create session with built-in retry support
	session := requests.NewSession().
		WithBaseURL("https://httpbin.org").
		WithTimeout(30 * time.Second).
		WithRetry(requests.RetryPolicy{
			MaxAttempts:     3,
			InitialInterval: 100 * time.Millisecond,
			MaxInterval:     2 * time.Second,
			Multiplier:      2.0,
			Jitter:          0.1,
			RetryIf: func(resp *requests.Response, err error) bool {
				if err != nil {
					fmt.Println("  Retrying due to error...")
					return true
				}
				if resp != nil && (resp.StatusCode >= 500 || resp.StatusCode == 429) {
					fmt.Printf("  Retrying due to status %d...\n", resp.StatusCode)
					return true
				}
				return false
			},
		})
	defer session.Close()

	// Make request - will automatically retry on failure
	req, _ := requests.NewGet("/get").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("  Session Retry Error: %v\n", err)
		return
	}
	fmt.Printf("  Success! Status: %d\n", resp.StatusCode)
	fmt.Println()
}

func sessionMiddlewareExample() {
	fmt.Println("--- Example 10: Session with Middleware ---")

	// Create a logging middleware
	loggingMiddleware := requests.MiddlewareFunc(func(req *requests.Request, next requests.Handler) (*requests.Response, error) {
		start := time.Now()
		fmt.Printf("  [Middleware] Starting request: %s %s\n", req.Method, req.URL)

		resp, err := next(req)

		duration := time.Since(start)
		if resp != nil {
			fmt.Printf("  [Middleware] Completed: %d in %v\n", resp.StatusCode, duration)
		} else if err != nil {
			fmt.Printf("  [Middleware] Failed: %v in %v\n", err, duration)
		}
		return resp, err
	})

	// Create an auth middleware
	authMiddleware := requests.MiddlewareFunc(func(req *requests.Request, next requests.Handler) (*requests.Response, error) {
		req.SetHeader("X-Auth-Token", "secret-token-123")
		fmt.Println("  [Auth Middleware] Added auth header")
		return next(req)
	})

	// Create session with multiple middlewares
	session := requests.NewSession().
		WithBaseURL("https://httpbin.org").
		WithMiddleware(loggingMiddleware).
		WithMiddleware(authMiddleware)
	defer session.Close()

	// Make request - middlewares will be executed in order
	req, _ := requests.NewGet("/headers").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("  Session Middleware Error: %v\n", err)
		return
	}
	fmt.Printf("  Response preview: %s...\n", resp.Text()[:100])
	fmt.Println()
}
