package url

import "net/http"

type Headers struct {
	headers http.Header
}

func NewHeaders() *Headers {
	return &Headers{headers: http.Header{}}
}

func (h *Headers) Add(key, value string) {
	h.headers.Add(key, value)
}

func (h *Headers) Set(key, value string) {
	h.headers.Set(key, value)
}

func (h *Headers) Get(key string) string {
	return h.headers.Get(key)
}
