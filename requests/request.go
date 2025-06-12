package requests

import (
	"context"
	"io"
	"net/http"
	neturl "net/url"

	customurl "github.com/sunerpy/requests/url"
)

type Request struct {
	Method  string
	URL     *neturl.URL
	Headers http.Header
	Body    io.Reader
	Context context.Context
	Params  *customurl.Values
}

func genRequest(ctx context.Context, method, rawURL string, params *customurl.Values, body io.Reader) (*Request, error) {
	parsedURL, err := neturl.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if params != nil {
		parsedURL.RawQuery = params.Encode()
	}
	return &Request{
		Method:  method,
		URL:     parsedURL,
		Headers: http.Header{},
		Body:    body,
		Context: ctx,
		Params:  params,
	}, nil
}

func NewRequestWithContext(ctx context.Context, method, rawURL string, params *customurl.Values, body io.Reader) (*Request, error) {
	return genRequest(ctx, method, rawURL, params, body)
}

func NewRequest(method, rawURL string, params *customurl.Values, body io.Reader) (*Request, error) {
	return genRequest(context.Background(), method, rawURL, params, body)
}

func (r *Request) AddHeader(key, value string) {
	r.Headers.Add(key, value)
}
