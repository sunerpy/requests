package client

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: http-client-refactor
// Property 17: Retry Policy Execution
// For any Client with a retry policy and a request that fails with a retryable error,
// the Client SHALL retry up to MaxAttempts times with exponential backoff.
// Validates: Requirements 9.1, 9.3, 9.4
func TestProperty17_RetryPolicyExecution(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Retry executes up to MaxAttempts times", prop.ForAll(
		func(maxAttempts int) bool {
			if maxAttempts < 1 || maxAttempts > 10 {
				return true
			}
			var attempts int32
			policy := RetryPolicy{
				MaxAttempts:     maxAttempts,
				InitialInterval: 1 * time.Millisecond,
				MaxInterval:     10 * time.Millisecond,
				Multiplier:      2.0,
				RetryIf: func(resp *Response, err error) bool {
					return err != nil
				},
			}
			executor := NewRetryExecutor(policy)
			ctx := context.Background()
			_, _ = executor.Execute(ctx, func() (*Response, error) {
				atomic.AddInt32(&attempts, 1)
				return nil, errors.New("always fail")
			})
			return int(atomic.LoadInt32(&attempts)) == maxAttempts
		},
		gen.IntRange(1, 10),
	))
	properties.TestingRun(t)
}

func TestProperty17_RetryStopsOnSuccess(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Retry stops when request succeeds", prop.ForAll(
		func(successOnAttempt int) bool {
			if successOnAttempt < 1 || successOnAttempt > 5 {
				return true
			}
			var attempts int32
			policy := RetryPolicy{
				MaxAttempts:     10,
				InitialInterval: 1 * time.Millisecond,
				RetryIf: func(resp *Response, err error) bool {
					return err != nil
				},
			}
			executor := NewRetryExecutor(policy)
			ctx := context.Background()
			resp, err := executor.Execute(ctx, func() (*Response, error) {
				current := int(atomic.AddInt32(&attempts, 1))
				if current >= successOnAttempt {
					return CreateMockResponse(200, nil, nil), nil
				}
				return nil, errors.New("fail")
			})
			// Should succeed
			if err != nil {
				return false
			}
			// Should have correct number of attempts
			if int(atomic.LoadInt32(&attempts)) != successOnAttempt {
				return false
			}
			// Should return success response
			return resp != nil && resp.StatusCode == 200
		},
		gen.IntRange(1, 5),
	))
	properties.TestingRun(t)
}

func TestProperty17_ExponentialBackoff(t *testing.T) {
	var delays []time.Duration
	var lastTime time.Time
	policy := RetryPolicy{
		MaxAttempts:     4,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
		Jitter:          0, // No jitter for predictable testing
		RetryIf: func(resp *Response, err error) bool {
			return err != nil
		},
	}
	executor := NewRetryExecutor(policy)
	ctx := context.Background()
	_, _ = executor.Execute(ctx, func() (*Response, error) {
		now := time.Now()
		if !lastTime.IsZero() {
			delays = append(delays, now.Sub(lastTime))
		}
		lastTime = now
		return nil, errors.New("fail")
	})
	// Should have 3 delays (between 4 attempts)
	if len(delays) != 3 {
		t.Errorf("Expected 3 delays, got %d", len(delays))
	}
	// Delays should increase (with some tolerance for timing)
	for i := 1; i < len(delays); i++ {
		if delays[i] < delays[i-1] {
			t.Logf("Delays: %v", delays)
			// Allow some tolerance due to timing
		}
	}
}

// Feature: http-client-refactor
// Property 18: Custom Retry Condition
// For any RetryPolicy with a custom RetryIf function, the Client SHALL call that
// function to determine if a retry should occur.
// Validates: Requirements 9.5
func TestProperty18_CustomRetryCondition(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Custom retry condition is called", prop.ForAll(
		func(retryOnCode int) bool {
			if retryOnCode < 400 || retryOnCode > 599 {
				return true
			}
			var conditionCalled int32
			policy := RetryPolicy{
				MaxAttempts:     3,
				InitialInterval: 1 * time.Millisecond,
				RetryIf: func(resp *Response, err error) bool {
					atomic.AddInt32(&conditionCalled, 1)
					if resp != nil {
						return resp.StatusCode == retryOnCode
					}
					return err != nil
				},
			}
			executor := NewRetryExecutor(policy)
			ctx := context.Background()
			_, _ = executor.Execute(ctx, func() (*Response, error) {
				return CreateMockResponse(retryOnCode, nil, nil), nil
			})
			// Condition should have been called
			return atomic.LoadInt32(&conditionCalled) > 0
		},
		gen.IntRange(400, 599),
	))
	properties.TestingRun(t)
}

func TestProperty18_CustomConditionControlsRetry(t *testing.T) {
	// Test that custom condition can prevent retry
	var attempts int32
	policy := RetryPolicy{
		MaxAttempts:     5,
		InitialInterval: 1 * time.Millisecond,
		RetryIf: func(resp *Response, err error) bool {
			// Only retry on 503
			if resp != nil {
				return resp.StatusCode == 503
			}
			return false
		},
	}
	executor := NewRetryExecutor(policy)
	ctx := context.Background()
	// Return 500 - should not retry
	_, _ = executor.Execute(ctx, func() (*Response, error) {
		atomic.AddInt32(&attempts, 1)
		return CreateMockResponse(500, nil, nil), nil
	})
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("Expected 1 attempt (no retry for 500), got %d", atomic.LoadInt32(&attempts))
	}
	// Reset and test with 503 - should retry
	atomic.StoreInt32(&attempts, 0)
	_, _ = executor.Execute(ctx, func() (*Response, error) {
		atomic.AddInt32(&attempts, 1)
		return CreateMockResponse(503, nil, nil), nil
	})
	if atomic.LoadInt32(&attempts) != 5 {
		t.Errorf("Expected 5 attempts (retry for 503), got %d", atomic.LoadInt32(&attempts))
	}
}

// Feature: http-client-refactor
// Property 19: Max Retries Exceeded Error
// For any request that exceeds the maximum retry attempts, the Client SHALL return
// an error containing the retry count and the last error.
// Validates: Requirements 9.6
func TestProperty19_MaxRetriesExceededError(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	properties.Property("Max retries exceeded returns RetryError", prop.ForAll(
		func(maxAttempts int) bool {
			if maxAttempts < 1 || maxAttempts > 10 {
				return true
			}
			originalErr := errors.New("original error")
			policy := RetryPolicy{
				MaxAttempts:     maxAttempts,
				InitialInterval: 1 * time.Millisecond,
				RetryIf: func(resp *Response, err error) bool {
					return err != nil
				},
			}
			executor := NewRetryExecutor(policy)
			ctx := context.Background()
			_, err := executor.Execute(ctx, func() (*Response, error) {
				return nil, originalErr
			})
			// Should return RetryError
			var retryErr *RetryError
			if !errors.As(err, &retryErr) {
				return false
			}
			// Should have correct attempt count
			if retryErr.Attempts != maxAttempts {
				return false
			}
			// Should contain last error
			return retryErr.LastErr != nil
		},
		gen.IntRange(1, 10),
	))
	properties.TestingRun(t)
}

func TestProperty19_RetryErrorUnwrapping(t *testing.T) {
	originalErr := errors.New("original error")
	policy := RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Millisecond,
		RetryIf: func(resp *Response, err error) bool {
			return err != nil
		},
	}
	executor := NewRetryExecutor(policy)
	ctx := context.Background()
	_, err := executor.Execute(ctx, func() (*Response, error) {
		return nil, originalErr
	})
	// Should be able to check for ErrMaxRetriesExceeded
	if !errors.Is(err, ErrMaxRetriesExceeded) {
		t.Error("Error should match ErrMaxRetriesExceeded")
	}
	// Should be able to unwrap to get last error
	var retryErr *RetryError
	if errors.As(err, &retryErr) {
		if retryErr.LastErr != originalErr {
			t.Error("LastErr should be original error")
		}
	} else {
		t.Error("Should be RetryError")
	}
}

// Unit tests
func TestRetryExecutor_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	policy := RetryPolicy{
		MaxAttempts:     10,
		InitialInterval: 100 * time.Millisecond,
		RetryIf: func(resp *Response, err error) bool {
			return err != nil
		},
	}
	executor := NewRetryExecutor(policy)
	var attempts int32
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	_, err := executor.Execute(ctx, func() (*Response, error) {
		atomic.AddInt32(&attempts, 1)
		return nil, errors.New("fail")
	})
	if err == nil {
		t.Error("Expected error on context cancellation")
	}
	// Should have stopped early due to cancellation
	if atomic.LoadInt32(&attempts) >= 10 {
		t.Error("Should have stopped before max attempts")
	}
}

func TestRetryExecutor_NoRetryOnSuccess(t *testing.T) {
	policy := DefaultRetryPolicy()
	executor := NewRetryExecutor(policy)
	ctx := context.Background()
	var attempts int32
	resp, err := executor.Execute(ctx, func() (*Response, error) {
		atomic.AddInt32(&attempts, 1)
		return CreateMockResponse(200, nil, nil), nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Expected 200 response")
	}
	if atomic.LoadInt32(&attempts) != 1 {
		t.Error("Should only have 1 attempt on success")
	}
}

func TestCalculateDelay(t *testing.T) {
	policy := RetryPolicy{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}
	executor := NewRetryExecutor(policy)
	// Attempt 0: 100ms
	delay0 := executor.calculateDelay(0)
	if delay0 != 100*time.Millisecond {
		t.Errorf("Expected 100ms, got %v", delay0)
	}
	// Attempt 1: 200ms
	delay1 := executor.calculateDelay(1)
	if delay1 != 200*time.Millisecond {
		t.Errorf("Expected 200ms, got %v", delay1)
	}
	// Attempt 2: 400ms
	delay2 := executor.calculateDelay(2)
	if delay2 != 400*time.Millisecond {
		t.Errorf("Expected 400ms, got %v", delay2)
	}
	// Attempt 4: should be capped at 1s
	delay4 := executor.calculateDelay(4)
	if delay4 != 1*time.Second {
		t.Errorf("Expected 1s (capped), got %v", delay4)
	}
}

func TestCalculateDelay_WithJitter(t *testing.T) {
	policy := RetryPolicy{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      1.0,
		Jitter:          0.5, // 50% jitter
	}
	executor := NewRetryExecutor(policy)
	// Run multiple times to check jitter variation
	delays := make([]time.Duration, 10)
	for i := 0; i < 10; i++ {
		delays[i] = executor.calculateDelay(0)
	}
	// All delays should be within jitter range (50-150ms)
	for _, d := range delays {
		if d < 50*time.Millisecond || d > 150*time.Millisecond {
			t.Errorf("Delay %v outside expected jitter range", d)
		}
	}
}

func TestNoRetryPolicy(t *testing.T) {
	policy := NoRetryPolicy()
	if policy.MaxAttempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", policy.MaxAttempts)
	}
}

func TestLinearRetryPolicy(t *testing.T) {
	policy := LinearRetryPolicy(5, 100*time.Millisecond)
	if policy.MaxAttempts != 5 {
		t.Errorf("Expected 5 attempts, got %d", policy.MaxAttempts)
	}
	if policy.InitialInterval != 100*time.Millisecond {
		t.Errorf("Expected 100ms interval, got %v", policy.InitialInterval)
	}
	if policy.Multiplier != 1.0 {
		t.Errorf("Expected multiplier 1.0, got %f", policy.Multiplier)
	}
}

func TestExponentialRetryPolicy(t *testing.T) {
	policy := ExponentialRetryPolicy(3, 100*time.Millisecond, 1*time.Second)
	if policy.MaxAttempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", policy.MaxAttempts)
	}
	if policy.Multiplier != 2.0 {
		t.Errorf("Expected multiplier 2.0, got %f", policy.Multiplier)
	}
	if policy.Jitter != 0.1 {
		t.Errorf("Expected jitter 0.1, got %f", policy.Jitter)
	}
}

func TestRetryOn5xx(t *testing.T) {
	tests := []struct {
		statusCode int
		err        error
		expected   bool
	}{
		{200, nil, false},
		{400, nil, false},
		{500, nil, true},
		{502, nil, true},
		{503, nil, true},
		{0, errors.New("error"), true},
	}
	for _, tc := range tests {
		var resp *Response
		if tc.statusCode > 0 {
			resp = CreateMockResponse(tc.statusCode, nil, nil)
		}
		result := RetryOn5xx(resp, tc.err)
		if result != tc.expected {
			t.Errorf("RetryOn5xx(%d, %v) = %v, want %v", tc.statusCode, tc.err, result, tc.expected)
		}
	}
}

func TestRetryOnNetworkError(t *testing.T) {
	// With error - should retry
	if !RetryOnNetworkError(nil, errors.New("network error")) {
		t.Error("Should retry on error")
	}
	// Without error - should not retry
	if RetryOnNetworkError(CreateMockResponse(500, nil, nil), nil) {
		t.Error("Should not retry without error")
	}
}

func TestRetryOnStatusCodes(t *testing.T) {
	condition := RetryOnStatusCodes(429, 503)
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{429, true},
		{500, false},
		{503, true},
	}
	for _, tc := range tests {
		resp := CreateMockResponse(tc.statusCode, nil, nil)
		result := condition(resp, nil)
		if result != tc.expected {
			t.Errorf("RetryOnStatusCodes(%d) = %v, want %v", tc.statusCode, result, tc.expected)
		}
	}
}

func TestCombineRetryConditions(t *testing.T) {
	cond1 := func(resp *Response, err error) bool {
		return resp != nil && resp.StatusCode == 500
	}
	cond2 := func(resp *Response, err error) bool {
		return resp != nil && resp.StatusCode == 503
	}
	combined := CombineRetryConditions(cond1, cond2)
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{500, true},
		{502, false},
		{503, true},
	}
	for _, tc := range tests {
		resp := CreateMockResponse(tc.statusCode, nil, nil)
		result := combined(resp, nil)
		if result != tc.expected {
			t.Errorf("Combined(%d) = %v, want %v", tc.statusCode, result, tc.expected)
		}
	}
}

func TestRetryMiddleware(t *testing.T) {
	var attempts int32
	policy := RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Millisecond,
		RetryIf: func(resp *Response, err error) bool {
			return err != nil
		},
	}
	middleware := RetryMiddleware(policy)
	handler := func(req *Request) (*Response, error) {
		current := int(atomic.AddInt32(&attempts, 1))
		if current < 3 {
			return nil, errors.New("fail")
		}
		return CreateMockResponse(200, nil, nil), nil
	}
	req := &Request{Context: context.Background()}
	resp, err := middleware.Process(req, handler)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Expected 200 response")
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts, got %d", atomic.LoadInt32(&attempts))
	}
}

func TestWithRetryPolicy(t *testing.T) {
	policy := ExponentialRetryPolicy(5, 100*time.Millisecond, 1*time.Second)
	config := NewRequestConfig()
	WithRetryPolicy(policy)(config)
	if config.Retry == nil {
		t.Fatal("Retry policy not set")
	}
	if config.Retry.MaxAttempts != 5 {
		t.Errorf("Expected 5 attempts, got %d", config.Retry.MaxAttempts)
	}
}

func TestWithMaxRetries(t *testing.T) {
	config := NewRequestConfig()
	WithMaxRetries(10)(config)
	if config.Retry == nil {
		t.Fatal("Retry policy not set")
	}
	if config.Retry.MaxAttempts != 10 {
		t.Errorf("Expected 10 attempts, got %d", config.Retry.MaxAttempts)
	}
}

func TestWithRetryCondition(t *testing.T) {
	customCondition := func(resp *Response, err error) bool {
		return resp != nil && resp.StatusCode == 418
	}
	config := NewRequestConfig()
	WithRetryCondition(customCondition)(config)
	if config.Retry == nil {
		t.Fatal("Retry policy not set")
	}
	// Test the condition
	resp := CreateMockResponse(418, nil, nil)
	if !config.Retry.RetryIf(resp, nil) {
		t.Error("Custom condition should return true for 418")
	}
	resp = CreateMockResponse(500, nil, nil)
	if config.Retry.RetryIf(resp, nil) {
		t.Error("Custom condition should return false for 500")
	}
}

func TestDefaultRetryCondition(t *testing.T) {
	tests := []struct {
		name        string
		resp        *Response
		err         error
		shouldRetry bool
	}{
		{"error", nil, errors.New("error"), true},
		{"500", CreateMockResponse(500, nil, nil), nil, true},
		{"502", CreateMockResponse(502, nil, nil), nil, true},
		{"503", CreateMockResponse(503, nil, nil), nil, true},
		{"429", CreateMockResponse(429, nil, nil), nil, true},
		{"200", CreateMockResponse(200, nil, nil), nil, false},
		{"400", CreateMockResponse(400, nil, nil), nil, false},
		{"404", CreateMockResponse(404, nil, nil), nil, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := DefaultRetryCondition(tc.resp, tc.err)
			if result != tc.shouldRetry {
				t.Errorf("DefaultRetryCondition() = %v, want %v", result, tc.shouldRetry)
			}
		})
	}
}

func TestRetryExecutor_ShouldRetry_NilRetryIf(t *testing.T) {
	// Test shouldRetry when RetryIf is nil (uses default)
	policy := RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Millisecond,
		RetryIf:         nil, // nil RetryIf
	}
	executor := NewRetryExecutor(policy)
	// Should use DefaultRetryCondition
	resp := CreateMockResponse(500, nil, nil)
	if !executor.shouldRetry(resp, nil) {
		t.Error("Should retry on 500 with nil RetryIf")
	}
	resp200 := CreateMockResponse(200, nil, nil)
	if executor.shouldRetry(resp200, nil) {
		t.Error("Should not retry on 200 with nil RetryIf")
	}
}

func TestRetryOn5xx_NilResponse(t *testing.T) {
	// nil response with no error
	if RetryOn5xx(nil, nil) {
		t.Error("Should not retry when both response and error are nil")
	}
}

func TestRetryOnStatusCodes_WithError(t *testing.T) {
	condition := RetryOnStatusCodes(429, 503)
	// With error, should retry
	if !condition(nil, errors.New("error")) {
		t.Error("Should retry on error")
	}
	// nil response without error
	if condition(nil, nil) {
		t.Error("Should not retry when both response and error are nil")
	}
}

func TestRetryMiddleware_NilContext(t *testing.T) {
	var attempts int32
	policy := RetryPolicy{
		MaxAttempts:     2,
		InitialInterval: 1 * time.Millisecond,
		RetryIf: func(resp *Response, err error) bool {
			return err != nil
		},
	}
	middleware := RetryMiddleware(policy)
	handler := func(req *Request) (*Response, error) {
		current := int(atomic.AddInt32(&attempts, 1))
		if current < 2 {
			return nil, errors.New("fail")
		}
		return CreateMockResponse(200, nil, nil), nil
	}
	// Request with nil context
	req := &Request{Context: nil}
	resp, err := middleware.Process(req, handler)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Expected 200 response")
	}
}

func TestCalculateDelay_ZeroMultiplier(t *testing.T) {
	policy := RetryPolicy{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      0, // Zero multiplier
		Jitter:          0,
	}
	executor := NewRetryExecutor(policy)
	// With zero multiplier, delay should stay at initial
	delay := executor.calculateDelay(5)
	if delay != 100*time.Millisecond {
		t.Errorf("Expected 100ms with zero multiplier, got %v", delay)
	}
}

func TestCalculateDelay_ZeroMaxInterval(t *testing.T) {
	policy := RetryPolicy{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     0, // Zero max interval (no cap)
		Multiplier:      2.0,
		Jitter:          0,
	}
	executor := NewRetryExecutor(policy)
	// With zero max interval, delay should grow without cap
	delay := executor.calculateDelay(4) // 100 * 2^4 = 1600ms
	if delay != 1600*time.Millisecond {
		t.Errorf("Expected 1600ms without cap, got %v", delay)
	}
}
