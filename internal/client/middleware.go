package client

import "github.com/sunerpy/requests/internal/models"

// MiddlewareChain manages a chain of middleware.
type MiddlewareChain struct {
	middlewares []Middleware
}

// NewMiddlewareChain creates a new middleware chain.
func NewMiddlewareChain(middlewares ...Middleware) *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: middlewares,
	}
}

// Use adds a middleware to the chain.
func (c *MiddlewareChain) Use(m Middleware) *MiddlewareChain {
	c.middlewares = append(c.middlewares, m)
	return c
}

// UseFunc adds a middleware function to the chain.
func (c *MiddlewareChain) UseFunc(fn func(req *Request, next Handler) (*models.Response, error)) *MiddlewareChain {
	return c.Use(MiddlewareFunc(fn))
}

// Execute runs the middleware chain with the given handler as the final handler.
func (c *MiddlewareChain) Execute(req *Request, finalHandler Handler) (*models.Response, error) {
	if len(c.middlewares) == 0 {
		return finalHandler(req)
	}
	// Build the chain from the end
	handler := finalHandler
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		m := c.middlewares[i]
		next := handler
		handler = func(r *Request) (*models.Response, error) {
			return m.Process(r, next)
		}
	}
	return handler(req)
}

// Clone creates a copy of the middleware chain.
func (c *MiddlewareChain) Clone() *MiddlewareChain {
	clone := &MiddlewareChain{
		middlewares: make([]Middleware, len(c.middlewares)),
	}
	copy(clone.middlewares, c.middlewares)
	return clone
}

// Len returns the number of middlewares in the chain.
func (c *MiddlewareChain) Len() int {
	return len(c.middlewares)
}

// LoggingMiddleware creates a middleware that logs requests and responses.
func LoggingMiddleware(logger func(format string, args ...any)) Middleware {
	return MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		logger("Request: %s %s", req.Method, req.URL)
		resp, err := next(req)
		if err != nil {
			logger("Error: %v", err)
		} else {
			logger("Response: %d %s", resp.StatusCode, resp.Status)
		}
		return resp, err
	})
}

// HeaderMiddleware creates a middleware that adds headers to all requests.
func HeaderMiddleware(headers map[string]string) Middleware {
	return MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		for k, v := range headers {
			if req.Headers.Get(k) == "" {
				req.Headers.Set(k, v)
			}
		}
		return next(req)
	})
}

// UserAgentMiddleware creates a middleware that sets the User-Agent header.
func UserAgentMiddleware(userAgent string) Middleware {
	return HeaderMiddleware(map[string]string{"User-Agent": userAgent})
}

// AuthMiddleware creates a middleware that adds authorization header.
func AuthMiddleware(authHeader string) Middleware {
	return MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		if req.Headers.Get("Authorization") == "" {
			req.Headers.Set("Authorization", authHeader)
		}
		return next(req)
	})
}

// BearerTokenMiddleware creates a middleware that adds bearer token auth.
func BearerTokenMiddleware(token string) Middleware {
	return AuthMiddleware("Bearer " + token)
}

// BasicAuthMiddleware creates a middleware that adds basic auth.
func BasicAuthMiddleware(username, password string) Middleware {
	return AuthMiddleware("Basic " + basicAuth(username, password))
}

// RecoveryMiddleware creates a middleware that recovers from panics.
func RecoveryMiddleware(onPanic func(req *Request, recovered any)) Middleware {
	return MiddlewareFunc(func(req *Request, next Handler) (resp *models.Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				if onPanic != nil {
					onPanic(req, r)
				}
				err = &RequestError{
					Op:  "Middleware",
					URL: req.URL.String(),
					Err: ErrPanic,
				}
			}
		}()
		return next(req)
	})
}

// ConditionalMiddleware creates a middleware that only executes if condition is true.
func ConditionalMiddleware(condition func(*Request) bool, m Middleware) Middleware {
	return MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		if condition(req) {
			return m.Process(req, next)
		}
		return next(req)
	})
}

// ChainMiddleware combines multiple middlewares into one.
func ChainMiddleware(middlewares ...Middleware) Middleware {
	return MiddlewareFunc(func(req *Request, next Handler) (*models.Response, error) {
		chain := NewMiddlewareChain(middlewares...)
		return chain.Execute(req, next)
	})
}
