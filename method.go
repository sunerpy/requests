package requests

import "github.com/sunerpy/requests/internal/client"

// Method is an alias for client.Method - represents an HTTP method.
type Method = client.Method

// HTTP methods as constants - re-exported from client package.
const (
	MethodGet     = client.MethodGet
	MethodPost    = client.MethodPost
	MethodPut     = client.MethodPut
	MethodDelete  = client.MethodDelete
	MethodPatch   = client.MethodPatch
	MethodHead    = client.MethodHead
	MethodOptions = client.MethodOptions
	MethodConnect = client.MethodConnect
	MethodTrace   = client.MethodTrace
)
