// Package examples demonstrates Session configuration patterns with the requests library.
//
// This file shows how to configure Session with:
// - Retry policies
// - Middleware
// - Method chaining
// - Various configuration options
package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/sunerpy/requests"
)

// BasicSessionExample demonstrates creating and configuring a basic session.
func BasicSessionExample() {
	// Create a session with common configuration
	session := requests.NewSession().
		WithBaseURL("https://api.example.com").
		WithTimeout(30*time.Second).
		WithHeader("User-Agent", "MyApp/1.0").
		WithHeader("Accept", "application/json")

	defer session.Close()

	// All requests will use the base URL and default headers
	req, _ := requests.NewGet("/users").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
}

// RetryPolicyExample demonstrates configuring retry behavior.
func RetryPolicyExample() {
	// Create a session with retry policy
	session := requests.NewSession().
		WithBaseURL("https://api.example.com").
		WithRetry(requests.RetryPolicy{
			MaxAttempts:     3,                      // Try up to 3 times
			InitialInterval: 100 * time.Millisecond, // Start with 100ms delay
			MaxInterval:     5 * time.Second,        // Cap delay at 5 seconds
			Multiplier:      2.0,                    // Double delay each retry
			Jitter:          0.1,                    // Add 10% randomness
			RetryIf: func(resp *requests.Response, err error) bool {
				// Retry on network errors
				if err != nil {
					return true
				}
				// Retry on server errors (5xx) and rate limiting (429)
				if resp != nil {
					return resp.StatusCode >= 500 || resp.StatusCode == 429
				}
				return false
			},
		})

	defer session.Close()

	req, _ := requests.NewGet("/data").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("Request failed after retries: %v", err)
		return
	}

	fmt.Printf("Success! Status: %d\n", resp.StatusCode)
}

// CustomRetryConditionExample shows how to create custom retry conditions.
func CustomRetryConditionExample() {
	// Retry only on specific status codes
	retryOnSpecificCodes := func(resp *requests.Response, err error) bool {
		if err != nil {
			return true
		}
		if resp != nil {
			// Only retry on 502, 503, 504
			switch resp.StatusCode {
			case 502, 503, 504:
				return true
			}
		}
		return false
	}

	session := requests.NewSession().
		WithRetry(requests.RetryPolicy{
			MaxAttempts:     5,
			InitialInterval: 200 * time.Millisecond,
			MaxInterval:     10 * time.Second,
			Multiplier:      1.5,
			RetryIf:         retryOnSpecificCodes,
		})

	defer session.Close()

	req, _ := requests.NewGet("https://api.example.com/data").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
}

// MiddlewareExample demonstrates adding middleware to a session.
func MiddlewareExample() {
	// Create a logging middleware
	loggingMiddleware := requests.MiddlewareFunc(func(req *requests.Request, next requests.Handler) (*requests.Response, error) {
		start := time.Now()
		fmt.Printf("[LOG] %s %s\n", req.Method, req.URL)

		resp, err := next(req)

		duration := time.Since(start)
		if resp != nil {
			fmt.Printf("[LOG] %d %s (%v)\n", resp.StatusCode, req.URL, duration)
		}
		return resp, err
	})

	// Create an auth middleware
	authMiddleware := requests.MiddlewareFunc(func(req *requests.Request, next requests.Handler) (*requests.Response, error) {
		// Add authentication header
		req.SetHeader("Authorization", "Bearer my-secret-token")
		return next(req)
	})

	// Create session with middlewares (executed in order)
	session := requests.NewSession().
		WithBaseURL("https://api.example.com").
		WithMiddleware(loggingMiddleware).
		WithMiddleware(authMiddleware)

	defer session.Close()

	req, _ := requests.NewGet("/protected/data").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
}

// MethodChainingExample demonstrates fluent configuration.
func MethodChainingExample() {
	// All configuration methods return Session for chaining
	session := requests.NewSession().
		WithBaseURL("https://api.example.com").
		WithTimeout(30*time.Second).
		WithIdleTimeout(90*time.Second).
		WithMaxIdleConns(100).
		WithKeepAlive(true).
		WithHTTP2(true).
		WithHeader("User-Agent", "MyApp/1.0").
		WithHeader("Accept", "application/json").
		WithBearerToken("my-token").
		WithRetry(requests.RetryPolicy{
			MaxAttempts:     3,
			InitialInterval: 100 * time.Millisecond,
			MaxInterval:     5 * time.Second,
			Multiplier:      2.0,
		})

	defer session.Close()

	// Session is fully configured and ready to use
	req, _ := requests.NewGet("/users").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
}

// SessionPoolExample demonstrates using session pooling for high-performance scenarios.
func SessionPoolExample() {
	// Acquire a session from the pool
	session := requests.AcquireSession()

	// Configure the session
	session = session.
		WithBaseURL("https://api.example.com").
		WithTimeout(10 * time.Second)

	// Make requests
	req, _ := requests.NewGet("/data").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
	} else {
		fmt.Printf("Status: %d\n", resp.StatusCode)
	}

	// Return session to pool when done
	requests.ReleaseSession(session)
}

// ProxyAndDNSExample demonstrates proxy and custom DNS configuration.
func ProxyAndDNSExample() {
	session := requests.NewSession().
		WithProxy("http://proxy.example.com:8080").
		WithDNS([]string{"8.8.8.8", "8.8.4.4"})

	defer session.Close()

	req, _ := requests.NewGet("https://api.example.com/data").Build()
	resp, err := session.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
}

// AuthenticationExample demonstrates various authentication methods.
func AuthenticationExample() {
	// Basic Auth
	basicSession := requests.NewSession().
		WithBasicAuth("username", "password")
	defer basicSession.Close()

	// Bearer Token
	bearerSession := requests.NewSession().
		WithBearerToken("my-jwt-token")
	defer bearerSession.Close()

	// Custom Auth Header
	customSession := requests.NewSession().
		WithHeader("X-API-Key", "my-api-key")
	defer customSession.Close()

	// Make requests with different auth methods
	req, _ := requests.NewGet("https://api.example.com/data").Build()

	resp1, _ := basicSession.Do(req)
	fmt.Printf("Basic Auth: %d\n", resp1.StatusCode)

	resp2, _ := bearerSession.Do(req)
	fmt.Printf("Bearer Token: %d\n", resp2.StatusCode)

	resp3, _ := customSession.Do(req)
	fmt.Printf("API Key: %d\n", resp3.StatusCode)
}

// CloneSessionExample demonstrates cloning a session.
func CloneSessionExample() {
	// Create a base session with common configuration
	baseSession := requests.NewSession().
		WithBaseURL("https://api.example.com").
		WithTimeout(30*time.Second).
		WithHeader("User-Agent", "MyApp/1.0")

	defer baseSession.Close()

	// Clone creates an independent copy
	clonedSession := baseSession.Clone()

	// The cloned session can be used independently
	// Note: Clone returns Client interface, not Session
	req, _ := requests.NewGet("/data").Build()
	resp, err := clonedSession.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
}

// ClearSessionExample demonstrates resetting a session.
func ClearSessionExample() {
	session := requests.NewSession().
		WithBaseURL("https://api.example.com").
		WithTimeout(30 * time.Second).
		WithBearerToken("token")

	// Make some requests...
	req, _ := requests.NewGet("/data").Build()
	_, _ = session.Do(req)

	// Clear resets the session to default state
	session = session.Clear()

	// Session is now reset - no base URL, no auth, default timeout
	defer session.Close()
}
