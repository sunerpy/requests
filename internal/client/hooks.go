package client

import (
	"sync"
	"time"
)

// RequestHook is called before a request is sent.
type (
	RequestHook func(req *Request)
	// ResponseHook is called after a response is received.
	ResponseHook func(req *Request, resp *Response, duration time.Duration)
	// ErrorHook is called when an error occurs.
	ErrorHook func(req *Request, err error, duration time.Duration)
	// Hooks manages request/response hooks.
	Hooks struct {
		mu         sync.RWMutex
		onRequest  []RequestHook
		onResponse []ResponseHook
		onError    []ErrorHook
	}
)

// NewHooks creates a new Hooks instance.
func NewHooks() *Hooks {
	return &Hooks{
		onRequest:  make([]RequestHook, 0),
		onResponse: make([]ResponseHook, 0),
		onError:    make([]ErrorHook, 0),
	}
}

// OnRequest registers a request hook.
func (h *Hooks) OnRequest(hook RequestHook) *Hooks {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onRequest = append(h.onRequest, hook)
	return h
}

// OnResponse registers a response hook.
func (h *Hooks) OnResponse(hook ResponseHook) *Hooks {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onResponse = append(h.onResponse, hook)
	return h
}

// OnError registers an error hook.
func (h *Hooks) OnError(hook ErrorHook) *Hooks {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onError = append(h.onError, hook)
	return h
}

// TriggerRequest calls all registered request hooks.
func (h *Hooks) TriggerRequest(req *Request) {
	h.mu.RLock()
	hooks := make([]RequestHook, len(h.onRequest))
	copy(hooks, h.onRequest)
	h.mu.RUnlock()
	for _, hook := range hooks {
		hook(req)
	}
}

// TriggerResponse calls all registered response hooks.
func (h *Hooks) TriggerResponse(req *Request, resp *Response, duration time.Duration) {
	h.mu.RLock()
	hooks := make([]ResponseHook, len(h.onResponse))
	copy(hooks, h.onResponse)
	h.mu.RUnlock()
	for _, hook := range hooks {
		hook(req, resp, duration)
	}
}

// TriggerError calls all registered error hooks.
func (h *Hooks) TriggerError(req *Request, err error, duration time.Duration) {
	h.mu.RLock()
	hooks := make([]ErrorHook, len(h.onError))
	copy(hooks, h.onError)
	h.mu.RUnlock()
	for _, hook := range hooks {
		hook(req, err, duration)
	}
}

// Clone creates a copy of the hooks.
func (h *Hooks) Clone() *Hooks {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clone := NewHooks()
	clone.onRequest = make([]RequestHook, len(h.onRequest))
	copy(clone.onRequest, h.onRequest)
	clone.onResponse = make([]ResponseHook, len(h.onResponse))
	copy(clone.onResponse, h.onResponse)
	clone.onError = make([]ErrorHook, len(h.onError))
	copy(clone.onError, h.onError)
	return clone
}

// Clear removes all registered hooks.
func (h *Hooks) Clear() *Hooks {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onRequest = make([]RequestHook, 0)
	h.onResponse = make([]ResponseHook, 0)
	h.onError = make([]ErrorHook, 0)
	return h
}

// Len returns the total number of registered hooks.
func (h *Hooks) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.onRequest) + len(h.onResponse) + len(h.onError)
}

// HooksMiddleware creates a middleware that triggers hooks.
func HooksMiddleware(hooks *Hooks) Middleware {
	return MiddlewareFunc(func(req *Request, next Handler) (*Response, error) {
		start := time.Now()
		// Trigger request hooks
		hooks.TriggerRequest(req)
		// Execute request
		resp, err := next(req)
		duration := time.Since(start)
		// Trigger appropriate hooks
		if err != nil {
			hooks.TriggerError(req, err, duration)
		} else {
			hooks.TriggerResponse(req, resp, duration)
		}
		return resp, err
	})
}

// LoggingHook creates a request hook that logs requests.
func LoggingHook(logger func(format string, args ...any)) RequestHook {
	return func(req *Request) {
		logger("Request: %s %s", req.Method, req.URL)
	}
}

// ResponseLoggingHook creates a response hook that logs responses.
func ResponseLoggingHook(logger func(format string, args ...any)) ResponseHook {
	return func(req *Request, resp *Response, duration time.Duration) {
		logger("Response: %s %s -> %d (%v)", req.Method, req.URL, resp.StatusCode, duration)
	}
}

// ErrorLoggingHook creates an error hook that logs errors.
func ErrorLoggingHook(logger func(format string, args ...any)) ErrorHook {
	return func(req *Request, err error, duration time.Duration) {
		logger("Error: %s %s -> %v (%v)", req.Method, req.URL, err, duration)
	}
}

// MetricsHook creates hooks for collecting metrics.
type MetricsHook struct {
	mu            sync.Mutex
	RequestCount  int64
	ResponseCount int64
	ErrorCount    int64
	TotalDuration time.Duration
}

// NewMetricsHook creates a new metrics hook.
func NewMetricsHook() *MetricsHook {
	return &MetricsHook{}
}

// RequestHook returns a request hook for counting requests.
func (m *MetricsHook) RequestHook() RequestHook {
	return func(req *Request) {
		m.mu.Lock()
		m.RequestCount++
		m.mu.Unlock()
	}
}

// ResponseHook returns a response hook for counting responses.
func (m *MetricsHook) ResponseHook() ResponseHook {
	return func(req *Request, resp *Response, duration time.Duration) {
		m.mu.Lock()
		m.ResponseCount++
		m.TotalDuration += duration
		m.mu.Unlock()
	}
}

// ErrorHook returns an error hook for counting errors.
func (m *MetricsHook) ErrorHook() ErrorHook {
	return func(req *Request, err error, duration time.Duration) {
		m.mu.Lock()
		m.ErrorCount++
		m.TotalDuration += duration
		m.mu.Unlock()
	}
}

// Stats returns the current metrics.
func (m *MetricsHook) Stats() (requests, responses, errors int64, avgDuration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	requests = m.RequestCount
	responses = m.ResponseCount
	errors = m.ErrorCount
	total := responses + errors
	if total > 0 {
		avgDuration = m.TotalDuration / time.Duration(total)
	}
	return requests, responses, errors, avgDuration
}

// Reset resets the metrics.
func (m *MetricsHook) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RequestCount = 0
	m.ResponseCount = 0
	m.ErrorCount = 0
	m.TotalDuration = 0
}
