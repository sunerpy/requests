package client

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/sunerpy/requests/internal/models"
)

// RequestHook is called before a request is sent.
type (
	RequestHook func(req *Request)
	// ResponseHook is called after a response is received.
	ResponseHook func(req *Request, resp *models.Response, duration time.Duration)
	// ErrorHook is called when an error occurs.
	ErrorHook func(req *Request, err error, duration time.Duration)

	// hooksData holds the immutable hook slices for atomic swap
	hooksData struct {
		onRequest  []RequestHook
		onResponse []ResponseHook
		onError    []ErrorHook
	}

	// Hooks manages request/response hooks using atomic operations for better performance.
	Hooks struct {
		data atomic.Value // holds *hooksData
		mu   sync.Mutex   // only used for write operations
	}
)

// NewHooks creates a new Hooks instance.
func NewHooks() *Hooks {
	h := &Hooks{}
	h.data.Store(&hooksData{
		onRequest:  make([]RequestHook, 0),
		onResponse: make([]ResponseHook, 0),
		onError:    make([]ErrorHook, 0),
	})
	return h
}

// getData returns the current hooks data (lock-free read).
func (h *Hooks) getData() *hooksData {
	return h.data.Load().(*hooksData)
}

// OnRequest registers a request hook.
func (h *Hooks) OnRequest(hook RequestHook) *Hooks {
	h.mu.Lock()
	defer h.mu.Unlock()

	old := h.getData()
	newHooks := make([]RequestHook, len(old.onRequest)+1)
	copy(newHooks, old.onRequest)
	newHooks[len(old.onRequest)] = hook

	h.data.Store(&hooksData{
		onRequest:  newHooks,
		onResponse: old.onResponse,
		onError:    old.onError,
	})
	return h
}

// OnResponse registers a response hook.
func (h *Hooks) OnResponse(hook ResponseHook) *Hooks {
	h.mu.Lock()
	defer h.mu.Unlock()

	old := h.getData()
	newHooks := make([]ResponseHook, len(old.onResponse)+1)
	copy(newHooks, old.onResponse)
	newHooks[len(old.onResponse)] = hook

	h.data.Store(&hooksData{
		onRequest:  old.onRequest,
		onResponse: newHooks,
		onError:    old.onError,
	})
	return h
}

// OnError registers an error hook.
func (h *Hooks) OnError(hook ErrorHook) *Hooks {
	h.mu.Lock()
	defer h.mu.Unlock()

	old := h.getData()
	newHooks := make([]ErrorHook, len(old.onError)+1)
	copy(newHooks, old.onError)
	newHooks[len(old.onError)] = hook

	h.data.Store(&hooksData{
		onRequest:  old.onRequest,
		onResponse: old.onResponse,
		onError:    newHooks,
	})
	return h
}

// TriggerRequest calls all registered request hooks (lock-free).
func (h *Hooks) TriggerRequest(req *Request) {
	data := h.getData()
	for _, hook := range data.onRequest {
		hook(req)
	}
}

// TriggerResponse calls all registered response hooks (lock-free).
func (h *Hooks) TriggerResponse(req *Request, resp *models.Response, duration time.Duration) {
	data := h.getData()
	for _, hook := range data.onResponse {
		hook(req, resp, duration)
	}
}

// TriggerError calls all registered error hooks (lock-free).
func (h *Hooks) TriggerError(req *Request, err error, duration time.Duration) {
	data := h.getData()
	for _, hook := range data.onError {
		hook(req, err, duration)
	}
}

// Clone creates a copy of the hooks.
func (h *Hooks) Clone() *Hooks {
	old := h.getData()

	clone := &Hooks{}
	newData := &hooksData{
		onRequest:  make([]RequestHook, len(old.onRequest)),
		onResponse: make([]ResponseHook, len(old.onResponse)),
		onError:    make([]ErrorHook, len(old.onError)),
	}
	copy(newData.onRequest, old.onRequest)
	copy(newData.onResponse, old.onResponse)
	copy(newData.onError, old.onError)
	clone.data.Store(newData)

	return clone
}

// Clear removes all registered hooks.
func (h *Hooks) Clear() *Hooks {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.data.Store(&hooksData{
		onRequest:  make([]RequestHook, 0),
		onResponse: make([]ResponseHook, 0),
		onError:    make([]ErrorHook, 0),
	})
	return h
}

// Len returns the total number of registered hooks.
func (h *Hooks) Len() int {
	data := h.getData()
	return len(data.onRequest) + len(data.onResponse) + len(data.onError)
}

// HooksMiddleware creates a middleware that triggers hooks.
func HooksMiddleware(hooks *Hooks) Middleware {
	return MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
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
	return func(req *Request, resp *models.Response, duration time.Duration) {
		logger("Response: %s %s -> %d (%v)", req.Method, req.URL, resp.StatusCode, duration)
	}
}

// ErrorLoggingHook creates an error hook that logs errors.
func ErrorLoggingHook(logger func(format string, args ...any)) ErrorHook {
	return func(req *Request, err error, duration time.Duration) {
		logger("Error: %s %s -> %v (%v)", req.Method, req.URL, err, duration)
	}
}

// MetricsHook creates hooks for collecting metrics using atomic operations.
type MetricsHook struct {
	requestCount  atomic.Int64
	responseCount atomic.Int64
	errorCount    atomic.Int64
	totalDuration atomic.Int64 // stored as nanoseconds
}

// NewMetricsHook creates a new metrics hook.
func NewMetricsHook() *MetricsHook {
	return &MetricsHook{}
}

// RequestHook returns a request hook for counting requests.
func (m *MetricsHook) RequestHook() RequestHook {
	return func(req *Request) {
		m.requestCount.Add(1)
	}
}

// ResponseHook returns a response hook for counting responses.
func (m *MetricsHook) ResponseHook() ResponseHook {
	return func(req *Request, resp *models.Response, duration time.Duration) {
		m.responseCount.Add(1)
		m.totalDuration.Add(int64(duration))
	}
}

// ErrorHook returns an error hook for counting errors.
func (m *MetricsHook) ErrorHook() ErrorHook {
	return func(req *Request, err error, duration time.Duration) {
		m.errorCount.Add(1)
		m.totalDuration.Add(int64(duration))
	}
}

// Stats returns the current metrics.
func (m *MetricsHook) Stats() (requests, responses, errors int64, avgDuration time.Duration) {
	requests = m.requestCount.Load()
	responses = m.responseCount.Load()
	errors = m.errorCount.Load()
	total := responses + errors
	if total > 0 {
		avgDuration = time.Duration(m.totalDuration.Load() / total)
	}
	return requests, responses, errors, avgDuration
}

// Reset resets the metrics.
func (m *MetricsHook) Reset() {
	m.requestCount.Store(0)
	m.responseCount.Store(0)
	m.errorCount.Store(0)
	m.totalDuration.Store(0)
}
