// Package requests provides a simple and easy-to-use HTTP client library for Go.
package requests

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"github.com/sunerpy/requests/internal/client"
	"github.com/sunerpy/requests/internal/models"
)

const (
	idleConnTimeout         = 90 * time.Second
	dnsResolveTimeout       = 15 * time.Second
	defaultDisableKeepAlive = false
	defaultMaxIdleConns     = 100
)

var (
	http1TransportPool  sync.Pool
	http2TransportPool  sync.Pool
	sessionPool         sync.Pool
	defaultHTTP2Enabled = false
	defaultHTTP2Lock    sync.Mutex
	defaultSess         Session
)

// Initialize Transport pools and create default Session
func init() {
	http1TransportPool = sync.Pool{
		New: func() any {
			return &http.Transport{
				MaxIdleConns:        defaultMaxIdleConns,
				MaxIdleConnsPerHost: defaultMaxIdleConns,
				IdleConnTimeout:     idleConnTimeout,
				DisableKeepAlives:   defaultDisableKeepAlive,
				ForceAttemptHTTP2:   false,
				TLSNextProto:        make(map[string]func(string, *tls.Conn) http.RoundTripper),
			}
		},
	}
	http2TransportPool = sync.Pool{
		New: func() any {
			return &http.Transport{
				MaxIdleConns:        defaultMaxIdleConns,
				MaxIdleConnsPerHost: defaultMaxIdleConns,
				IdleConnTimeout:     idleConnTimeout,
				DisableKeepAlives:   defaultDisableKeepAlive,
				ForceAttemptHTTP2:   true,
			}
		},
	}
	sessionPool = sync.Pool{
		New: func() any {
			return &defaultSession{
				headers:      make(http.Header, 4),
				keepAlive:    true,
				maxIdleConns: defaultMaxIdleConns,
				idleTimeout:  idleConnTimeout,
			}
		},
	}
	defaultSess = NewSession()
}

// GetTransport returns a transport from the pool.
func GetTransport(enableHTTP2 bool) *http.Transport {
	if enableHTTP2 {
		return http2TransportPool.Get().(*http.Transport)
	}
	return http1TransportPool.Get().(*http.Transport)
}

// PutTransport returns a transport to the pool.
func PutTransport(transport *http.Transport) {
	if transport == nil {
		return
	}
	if transport.ForceAttemptHTTP2 {
		http2TransportPool.Put(transport)
	} else {
		http1TransportPool.Put(transport)
	}
}

// AcquireSession gets a Session from the pool.
// Remember to call ReleaseSession when done.
func AcquireSession() Session {
	s := sessionPool.Get().(*defaultSession)
	defaultHTTP2Lock.Lock()
	useHTTP2 := defaultHTTP2Enabled
	defaultHTTP2Lock.Unlock()
	s.useHTTP2 = useHTTP2
	if useHTTP2 {
		s.client = &http.Client{Transport: http2TransportPool.Get().(*http.Transport)}
	} else {
		s.client = &http.Client{Transport: http1TransportPool.Get().(*http.Transport)}
	}
	return s
}

// ReleaseSession returns a Session to the pool.
func ReleaseSession(sess Session) {
	if sess == nil {
		return
	}
	s, ok := sess.(*defaultSession)
	if !ok {
		return
	}
	// Return transport to pool
	if tr, ok := s.client.Transport.(*http.Transport); ok {
		if tr.ForceAttemptHTTP2 {
			http2TransportPool.Put(tr)
		} else {
			http1TransportPool.Put(tr)
		}
	}
	// Reset session state
	s.baseURL = ""
	s.timeout = 0
	s.idleTimeout = idleConnTimeout
	s.proxyURL = nil
	s.dnsServers = nil
	s.authHeader = ""
	s.bearerToken = ""
	s.keepAlive = true
	s.maxIdleConns = defaultMaxIdleConns
	// Clear headers but keep the map
	for k := range s.headers {
		delete(s.headers, k)
	}
	s.client = nil
	sessionPool.Put(s)
}

// SetHTTP2Enabled sets the global HTTP/2 enabled state.
func SetHTTP2Enabled(enabled bool) {
	defaultHTTP2Lock.Lock()
	defer defaultHTTP2Lock.Unlock()
	defaultHTTP2Enabled = enabled
	defaultSess = defaultSess.WithHTTP2(enabled)
}

// IsHTTP2Enabled returns the global HTTP/2 enabled state.
func IsHTTP2Enabled() bool {
	defaultHTTP2Lock.Lock()
	defer defaultHTTP2Lock.Unlock()
	return defaultHTTP2Enabled
}

// Client is the HTTP client interface.
type Client interface {
	Do(req *Request) (*models.Response, error)
	Clone() Client
}

// Session extends Client with session management capabilities.
type Session interface {
	Client
	WithBaseURL(base string) Session
	WithTimeout(d time.Duration) Session
	WithProxy(proxyURL string) Session
	WithDNS(dnsServers []string) Session
	WithHeader(key, value string) Session
	WithHeaders(headers map[string]string) Session
	WithBasicAuth(username, password string) Session
	WithBearerToken(token string) Session
	WithHTTP2(enabled bool) Session
	WithKeepAlive(enabled bool) Session
	WithMaxIdleConns(maxIdle int) Session
	WithIdleTimeout(d time.Duration) Session
	WithCookieJar(jar http.CookieJar) Session
	Close() error
	Clear() Session
}

// defaultSession implements the Session interface.
type defaultSession struct {
	baseURL      string
	timeout      time.Duration
	idleTimeout  time.Duration
	proxyURL     *url.URL
	headers      http.Header
	client       *http.Client
	useHTTP2     bool
	keepAlive    bool
	maxIdleConns int
	clientLock   sync.Mutex
	dnsServers   []string
	authHeader   string
	bearerToken  string
}

// NewSession creates a new Session with default settings.
func NewSession() Session {
	defaultHTTP2Lock.Lock()
	defer defaultHTTP2Lock.Unlock()
	var transport *http.Transport
	if defaultHTTP2Enabled {
		transport = http2TransportPool.Get().(*http.Transport)
	} else {
		transport = http1TransportPool.Get().(*http.Transport)
	}
	jar, _ := cookiejar.New(nil)
	return &defaultSession{
		headers:      http.Header{},
		client:       &http.Client{Transport: transport, Jar: jar},
		useHTTP2:     defaultHTTP2Enabled,
		keepAlive:    true,
		maxIdleConns: defaultMaxIdleConns,
		idleTimeout:  idleConnTimeout,
	}
}

func (s *defaultSession) WithBaseURL(base string) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.baseURL = base
	return s
}

func (s *defaultSession) WithTimeout(d time.Duration) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.timeout = d
	if s.client != nil {
		s.client.Timeout = d
	}
	return s
}

func (s *defaultSession) WithIdleTimeout(d time.Duration) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.idleTimeout = d
	if tr, ok := s.client.Transport.(*http.Transport); ok {
		tr.IdleConnTimeout = d
	}
	return s
}

func (s *defaultSession) WithProxy(proxyURL string) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	u, err := url.Parse(proxyURL)
	if err == nil && u.Scheme != "" && u.Host != "" {
		s.proxyURL = u
		if tr, ok := s.client.Transport.(*http.Transport); ok {
			tr.Proxy = http.ProxyURL(u)
		}
	} else {
		s.proxyURL = nil
		if tr, ok := s.client.Transport.(*http.Transport); ok {
			tr.Proxy = nil
		}
	}
	return s
}

func (s *defaultSession) WithDNS(dnsServers []string) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.dnsServers = dnsServers
	if tr, ok := s.client.Transport.(*http.Transport); ok {
		tr.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			if len(s.dnsServers) > 0 {
				return customDial(ctx, network, address, s.dnsServers)
			}
			dialer := &net.Dialer{
				Timeout:   dnsResolveTimeout,
				DualStack: true,
			}
			return dialer.DialContext(ctx, network, address)
		}
	}
	return s
}

func customDial(ctx context.Context, network, address string, dnsServers []string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   dnsResolveTimeout,
		DualStack: true,
	}
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			for _, dns := range dnsServers {
				conn, err := dialer.DialContext(ctx, network, dns+":53")
				if err == nil {
					return conn, nil
				}
			}
			return nil, fmt.Errorf("failed to connect to any DNS servers: %v", dnsServers)
		},
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %v", err)
	}
	ips, err := resolver.LookupHost(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("DNS resolution failed for %s: %v", host, err)
	}
	for _, ip := range ips {
		conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
		if err == nil {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("failed to connect to any resolved IPs for %s, address: %v", address, ips)
}

func (s *defaultSession) WithBasicAuth(username, password string) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	auth := fmt.Sprintf("%s:%s", username, password)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	s.authHeader = fmt.Sprintf("Basic %s", encodedAuth)
	s.headers.Set("Authorization", s.authHeader)
	return s
}

func (s *defaultSession) WithBearerToken(token string) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.bearerToken = token
	s.headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return s
}

func (s *defaultSession) WithHeader(key, value string) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.headers.Set(key, value)
	return s
}

func (s *defaultSession) WithHeaders(headers map[string]string) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	for k, v := range headers {
		s.headers.Set(k, v)
	}
	return s
}

func (s *defaultSession) WithHTTP2(enabled bool) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	if s.useHTTP2 == enabled {
		return s
	}
	s.useHTTP2 = enabled
	if tr, ok := s.client.Transport.(*http.Transport); ok {
		if tr.ForceAttemptHTTP2 {
			http2TransportPool.Put(tr)
		} else {
			http1TransportPool.Put(tr)
		}
		if enabled {
			tr = http2TransportPool.Get().(*http.Transport)
		} else {
			tr = http1TransportPool.Get().(*http.Transport)
		}
		s.client.Transport = tr
	}
	return s
}

func (s *defaultSession) WithKeepAlive(enabled bool) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.keepAlive = enabled
	var tr *http.Transport
	if existingTransport, ok := s.client.Transport.(*http.Transport); ok {
		tr = existingTransport
	} else {
		tr = GetTransport(s.useHTTP2)
		s.client.Transport = tr
	}
	tr.DisableKeepAlives = !s.keepAlive
	return s
}

func (s *defaultSession) WithMaxIdleConns(maxIdle int) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.maxIdleConns = maxIdle
	var tr *http.Transport
	if existingTransport, ok := s.client.Transport.(*http.Transport); ok {
		tr = existingTransport
	} else {
		tr = GetTransport(s.useHTTP2)
		s.client.Transport = tr
	}
	tr.MaxIdleConns = maxIdle
	tr.MaxIdleConnsPerHost = maxIdle
	return s
}

func (s *defaultSession) WithCookieJar(jar http.CookieJar) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.client.Jar = jar
	return s
}

func (s *defaultSession) Clone() Client {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	var transport *http.Transport
	if s.useHTTP2 {
		transport = http2TransportPool.Get().(*http.Transport)
	} else {
		transport = http1TransportPool.Get().(*http.Transport)
	}
	jar, _ := cookiejar.New(nil)
	newClient := &http.Client{
		Transport: transport,
		Jar:       jar,
	}
	newHeaders := s.headers.Clone()
	newSession := &defaultSession{
		baseURL:      s.baseURL,
		timeout:      s.timeout,
		idleTimeout:  s.idleTimeout,
		proxyURL:     s.proxyURL,
		headers:      newHeaders,
		client:       newClient,
		useHTTP2:     s.useHTTP2,
		keepAlive:    s.keepAlive,
		maxIdleConns: s.maxIdleConns,
		bearerToken:  s.bearerToken,
	}
	if !s.useHTTP2 {
		newSession.client.Transport.(*http.Transport).TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	}
	return newSession
}

func (s *defaultSession) Close() error {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	if tr, ok := s.client.Transport.(*http.Transport); ok {
		if tr.ForceAttemptHTTP2 {
			http2TransportPool.Put(tr)
		} else {
			http1TransportPool.Put(tr)
		}
		s.client.Transport = nil
	}
	return nil
}

func (s *defaultSession) Clear() Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.baseURL = ""
	s.useHTTP2 = false
	s.keepAlive = true
	s.maxIdleConns = defaultMaxIdleConns
	s.timeout = 0
	s.idleTimeout = idleConnTimeout
	s.dnsServers = nil
	s.bearerToken = ""
	s.authHeader = ""
	s.headers = http.Header{}
	jar, _ := cookiejar.New(nil)
	s.client = &http.Client{
		Transport: GetTransport(false),
		Jar:       jar,
	}
	return s
}

// resolveURL resolves the request URL against the session's base URL.
func (s *defaultSession) resolveURL(reqURL *url.URL) (*url.URL, error) {
	if s.baseURL == "" {
		return reqURL, nil
	}
	base, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}
	return base.ResolveReference(reqURL), nil
}

// applyHeaders copies headers from source to destination.
func applyHeaders(dst, src http.Header) {
	for k, vals := range src {
		for _, v := range vals {
			dst.Add(k, v)
		}
	}
}

func (s *defaultSession) Do(req *Request) (*models.Response, error) {
	if req.URL == nil {
		return nil, errors.New("request URL cannot be nil")
	}
	if req.Context == nil {
		req.Context = context.Background()
	}
	finalURL, err := s.resolveURL(req.URL)
	if err != nil {
		return nil, err
	}
	ctx := req.Context
	if s.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
	}
	httpReq, err := http.NewRequestWithContext(ctx, req.Method.String(), finalURL.String(), req.Body)
	if err != nil {
		return nil, err
	}
	// Fast path: skip header iteration if no session headers
	if len(s.headers) > 0 {
		applyHeaders(httpReq.Header, s.headers)
	}
	// Fast path: skip header iteration if no request headers
	if len(req.Headers) > 0 {
		applyHeaders(httpReq.Header, req.Headers)
	}
	resp, err := s.client.Do(httpReq) //nolint:bodyclose // body is closed in models.NewResponse
	if err != nil {
		return nil, err
	}
	return models.NewResponse(resp, finalURL.String())
}

// DoFast is an optimized version of Do for simple requests without baseURL or timeout.
// Use this when you need maximum performance and don't need session-level configuration.
func (s *defaultSession) DoFast(req *Request) (*models.Response, error) {
	httpReq, err := http.NewRequestWithContext(req.Context, req.Method.String(), req.URL.String(), req.Body)
	if err != nil {
		return nil, err
	}
	if len(req.Headers) > 0 {
		httpReq.Header = req.Headers
	}
	resp, err := s.client.Do(httpReq) //nolint:bodyclose // body is closed in models.NewResponse
	if err != nil {
		return nil, err
	}
	return models.NewResponse(resp, req.URL.String())
}

// Request is an alias for client.Request - the unified request type.
// Use NewRequestBuilder or convenience constructors (NewGet, NewPost, etc.) to create requests.
type Request = client.Request
