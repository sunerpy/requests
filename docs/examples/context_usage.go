// Package examples demonstrates context usage patterns with the requests library.
//
// This file shows how to use context.Context for:
// - Request timeout control
// - Request cancellation
// - Deadline management
package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sunerpy/requests"
)

// TimeoutExample demonstrates using context with timeout.
// The request will be canceled if it takes longer than the specified timeout.
func TimeoutExample() {
	// Create a context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Always call cancel to release resources

	session := requests.NewSession()
	defer session.Close()

	req, err := requests.NewGet("https://api.example.com/data").Build()
	if err != nil {
		log.Fatal(err)
	}

	// DoWithContext will respect the context timeout
	resp, err := session.DoWithContext(ctx, req)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("Request timed out after 5 seconds")
			return
		}
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response body: %s\n", resp.Text())
}

// CancellationExample demonstrates canceling a request programmatically.
// This is useful when you need to cancel a request based on user action or other events.
func CancellationExample() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	session := requests.NewSession()
	defer session.Close()

	// Simulate cancellation after 2 seconds (e.g., user clicks cancel button)
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Canceling request...")
		cancel()
	}()

	req, err := requests.NewGet("https://api.example.com/slow-endpoint").Build()
	if err != nil {
		log.Fatal(err)
	}

	// This request will be canceled after 2 seconds
	resp, err := session.DoWithContext(ctx, req)
	if err != nil {
		if err == context.Canceled {
			fmt.Println("Request was canceled by user")
			return
		}
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
}

// DeadlineExample demonstrates using context with a specific deadline.
// The request must complete before the deadline.
func DeadlineExample() {
	// Set a deadline 10 seconds from now
	deadline := time.Now().Add(10 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	session := requests.NewSession()
	defer session.Close()

	req, err := requests.NewGet("https://api.example.com/data").Build()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := session.DoWithContext(ctx, req)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Printf("Request did not complete before deadline: %v\n", deadline)
			return
		}
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
}

// ContextWithSessionTimeoutExample shows how context timeout interacts with session timeout.
// The shorter timeout wins.
func ContextWithSessionTimeoutExample() {
	// Session has 30 second timeout
	session := requests.NewSession().
		WithTimeout(30 * time.Second)
	defer session.Close()

	// But context has 5 second timeout - this will be used
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := requests.NewGet("https://api.example.com/data").Build()
	if err != nil {
		log.Fatal(err)
	}

	// The 5 second context timeout will be respected
	resp, err := session.DoWithContext(ctx, req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
}

// ParallelRequestsWithCancellationExample demonstrates canceling multiple parallel requests.
func ParallelRequestsWithCancellationExample() {
	// Create a cancellable context for all requests
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := requests.NewSession()
	defer session.Close()

	urls := []string{
		"https://api.example.com/users",
		"https://api.example.com/posts",
		"https://api.example.com/comments",
	}

	results := make(chan string, len(urls))
	errors := make(chan error, len(urls))

	// Start parallel requests
	for _, url := range urls {
		go func(url string) {
			req, err := requests.NewGet(url).Build()
			if err != nil {
				errors <- err
				return
			}

			resp, err := session.DoWithContext(ctx, req)
			if err != nil {
				errors <- err
				return
			}

			results <- fmt.Sprintf("%s: %d", url, resp.StatusCode)
		}(url)
	}

	// Wait for first result or error
	select {
	case result := <-results:
		fmt.Printf("First result: %s\n", result)
		// Cancel remaining requests
		cancel()
	case err := <-errors:
		fmt.Printf("First error: %v\n", err)
		// Cancel remaining requests
		cancel()
	case <-time.After(10 * time.Second):
		fmt.Println("Overall timeout")
		cancel()
	}
}

// ContextValueExample demonstrates passing values through context.
// Note: This is for demonstration - the requests library doesn't use context values internally.
func ContextValueExample() {
	type requestIDKey struct{}

	// Create context with request ID
	ctx := context.WithValue(context.Background(), requestIDKey{}, "req-12345")

	session := requests.NewSession()
	defer session.Close()

	// You can use the request ID for logging or tracing
	requestID := ctx.Value(requestIDKey{}).(string)
	fmt.Printf("Making request with ID: %s\n", requestID)

	req, err := requests.NewGet("https://api.example.com/data").
		WithHeader("X-Request-ID", requestID).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := session.DoWithContext(ctx, req)
	if err != nil {
		log.Printf("Request %s failed: %v", requestID, err)
		return
	}

	fmt.Printf("Request %s completed with status: %d\n", requestID, resp.StatusCode)
}
