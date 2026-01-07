package client

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/sunerpy/requests/internal/models"
)

// RetryExecutor handles retry logic for HTTP requests.
type RetryExecutor struct {
	policy RetryPolicy
}

// NewRetryExecutor creates a new retry executor with the given policy.
func NewRetryExecutor(policy RetryPolicy) *RetryExecutor {
	return &RetryExecutor{policy: policy}
}

// Execute executes the given function with retry logic.
func (e *RetryExecutor) Execute(ctx context.Context, fn func() (*models.Response, error)) (*models.Response, error) {
	var lastErr error
	var lastResp *models.Response
	for attempt := 0; attempt < e.policy.MaxAttempts; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, &RequestError{Op: "Retry", Err: ctx.Err()}
		default:
		}
		resp, err := fn()
		lastResp = resp
		lastErr = err
		// Check if we should retry
		if !e.shouldRetry(resp, err) {
			return resp, err
		}
		// Don't sleep after the last attempt
		if attempt < e.policy.MaxAttempts-1 {
			delay := e.calculateDelay(attempt)
			select {
			case <-ctx.Done():
				return nil, &RequestError{Op: "Retry", Err: ctx.Err()}
			case <-time.After(delay):
			}
		}
	}
	return lastResp, &RetryError{
		Attempts: e.policy.MaxAttempts,
		LastErr:  lastErr,
	}
}

// shouldRetry determines if a request should be retried.
func (e *RetryExecutor) shouldRetry(resp *models.Response, err error) bool {
	if e.policy.RetryIf != nil {
		return e.policy.RetryIf(resp, err)
	}
	return DefaultRetryCondition(resp, err)
}

// calculateDelay calculates the delay before the next retry attempt.
func (e *RetryExecutor) calculateDelay(attempt int) time.Duration {
	delay := e.policy.InitialInterval
	// Apply exponential backoff
	if e.policy.Multiplier > 0 {
		multiplier := math.Pow(e.policy.Multiplier, float64(attempt))
		delay = time.Duration(float64(delay) * multiplier)
	}
	// Apply max interval cap
	if e.policy.MaxInterval > 0 && delay > e.policy.MaxInterval {
		delay = e.policy.MaxInterval
	}
	// Apply jitter
	if e.policy.Jitter > 0 {
		jitterRange := float64(delay) * e.policy.Jitter
		jitter := (rand.Float64()*2 - 1) * jitterRange
		delay = time.Duration(float64(delay) + jitter)
	}
	// Ensure delay is not negative
	if delay < 0 {
		delay = 0
	}
	return delay
}

// NoRetryPolicy returns a policy that never retries.
func NoRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts: 1,
	}
}

// LinearRetryPolicy returns a policy with linear backoff.
func LinearRetryPolicy(maxAttempts int, interval time.Duration) RetryPolicy {
	return RetryPolicy{
		MaxAttempts:     maxAttempts,
		InitialInterval: interval,
		MaxInterval:     interval,
		Multiplier:      1.0,
		Jitter:          0,
		RetryIf:         DefaultRetryCondition,
	}
}

// ExponentialRetryPolicy returns a policy with exponential backoff.
func ExponentialRetryPolicy(maxAttempts int, initialInterval, maxInterval time.Duration) RetryPolicy {
	return RetryPolicy{
		MaxAttempts:     maxAttempts,
		InitialInterval: initialInterval,
		MaxInterval:     maxInterval,
		Multiplier:      2.0,
		Jitter:          0.1,
		RetryIf:         DefaultRetryCondition,
	}
}

// RetryOn5xx returns a retry condition that retries on 5xx errors.
func RetryOn5xx(resp *models.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp != nil {
		return resp.StatusCode >= 500
	}
	return false
}

// RetryOnNetworkError returns a retry condition that only retries on network errors.
func RetryOnNetworkError(resp *models.Response, err error) bool {
	return err != nil
}

// RetryOnStatusCodes returns a retry condition that retries on specific status codes.
func RetryOnStatusCodes(codes ...int) func(*models.Response, error) bool {
	codeSet := make(map[int]bool)
	for _, code := range codes {
		codeSet[code] = true
	}
	return func(resp *models.Response, err error) bool {
		if err != nil {
			return true
		}
		if resp != nil {
			return codeSet[resp.StatusCode]
		}
		return false
	}
}

// CombineRetryConditions combines multiple retry conditions with OR logic.
func CombineRetryConditions(conditions ...func(*models.Response, error) bool) func(*models.Response, error) bool {
	return func(resp *models.Response, err error) bool {
		for _, cond := range conditions {
			if cond(resp, err) {
				return true
			}
		}
		return false
	}
}

// RetryMiddleware creates a middleware that adds retry logic.
func RetryMiddleware(policy RetryPolicy) Middleware {
	return MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		executor := NewRetryExecutor(policy)
		ctx := req.Context
		if ctx == nil {
			ctx = context.Background()
		}
		return executor.Execute(ctx, func() (*models.Response, error) {
			return next(req)
		})
	})
}

// WithRetryPolicy returns a request option that sets the retry policy.
func WithRetryPolicy(policy RetryPolicy) RequestOption {
	return func(c *RequestConfig) {
		c.Retry = &policy
	}
}

// WithMaxRetries returns a request option that sets the maximum retry attempts.
func WithMaxRetries(maxAttempts int) RequestOption {
	return func(c *RequestConfig) {
		if c.Retry == nil {
			policy := DefaultRetryPolicy()
			c.Retry = &policy
		}
		c.Retry.MaxAttempts = maxAttempts
	}
}

// WithRetryCondition returns a request option that sets a custom retry condition.
func WithRetryCondition(condition func(*models.Response, error) bool) RequestOption {
	return func(c *RequestConfig) {
		if c.Retry == nil {
			policy := DefaultRetryPolicy()
			c.Retry = &policy
		}
		c.Retry.RetryIf = condition
	}
}
