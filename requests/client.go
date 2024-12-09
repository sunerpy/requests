package requests

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/sunerpy/requests/models"
)

const (
	idleConnTimeout         = 90 * time.Second
	dnsResloveTimeout       = 15 * time.Second
	defaultDisableKeepAlive = false
	defaultMaxIdleConns     = 100
)

var (
	http1TransportPool  sync.Pool
	http2TransportPool  sync.Pool
	defaultHTTP2Enabled = false
	defaultHTTP2Lock    sync.Mutex
	defaultSess         Session
)

// 初始化 Transport 池并创建默认 Session
func init() {
	http1TransportPool = sync.Pool{
		New: func() interface{} {
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
		New: func() interface{} {
			return &http.Transport{
				MaxIdleConns:        defaultMaxIdleConns,
				MaxIdleConnsPerHost: defaultMaxIdleConns,
				IdleConnTimeout:     idleConnTimeout,
				DisableKeepAlives:   defaultDisableKeepAlive,
				ForceAttemptHTTP2:   true,
				// TLSNextProto 保持默认，支持 HTTP/2
			}
		},
	}
	defaultSess = NewSession()
}

func GetTransport(enableHTTP2 bool) *http.Transport {
	if enableHTTP2 {
		return http2TransportPool.Get().(*http.Transport)
	}
	return http1TransportPool.Get().(*http.Transport)
}

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

// SetHTTP2Enabled 设置全局 HTTP/2 启用状态，并更新默认 Session
func SetHTTP2Enabled(enabled bool) {
	defaultHTTP2Lock.Lock()
	defer defaultHTTP2Lock.Unlock()
	defaultHTTP2Enabled = enabled
	defaultSess = defaultSess.WithHTTP2(enabled)
}

// IsHTTP2Enabled 获取全局 HTTP/2 启用状态
func IsHTTP2Enabled() bool {
	defaultHTTP2Lock.Lock()
	defer defaultHTTP2Lock.Unlock()
	return defaultHTTP2Enabled
}

type Client interface {
	Do(req *Request) (*models.Response, error)
	Clone() Client
}
type Session interface {
	Client
	WithBaseURL(base string) Session
	WithTimeout(d time.Duration) Session
	WithProxy(proxyURL string) Session
	WithDNS(dnsServers []string) Session
	WithHeader(key, value string) Session
	WithBasicAuth(username, password string) Session
	WithHTTP2(enabled bool) Session
	WithKeepAlive(enabled bool) Session
	WithMaxIdleConns(maxIdle int) Session
	Close() error
	Clear() Session
}
type defaultSession struct {
	baseURL      string
	timeout      time.Duration
	proxyURL     *url.URL
	headers      http.Header
	client       *http.Client
	useHTTP2     bool
	keepAlive    bool
	maxIdleConns int
	clientLock   sync.Mutex
	dnsServers   []string
	authHeader   string
}

// NewSession 创建一个新的 Session，使用对应的 Transport 池
func NewSession() Session {
	defaultHTTP2Lock.Lock()
	defer defaultHTTP2Lock.Unlock()
	var transport *http.Transport
	if defaultHTTP2Enabled {
		transport = http2TransportPool.Get().(*http.Transport)
	} else {
		transport = http1TransportPool.Get().(*http.Transport)
	}
	return &defaultSession{
		headers:      http.Header{},
		client:       &http.Client{Transport: transport},
		useHTTP2:     defaultHTTP2Enabled,
		keepAlive:    true,
		maxIdleConns: defaultMaxIdleConns,
	}
}

func (s *defaultSession) WithBaseURL(base string) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.baseURL = base
	return s
}

// This function `WithTimeout` in the `defaultSession` struct is setting the timeout duration for the
// session. It acquires a lock on the client, sets the timeout duration to the provided value `d`, and
// the timeout duration is a request timeout.
func (s *defaultSession) WithTimeout(d time.Duration) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.timeout = d
	if s.client != nil {
		s.client.Timeout = d
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
	// 设置自定义的 DNS 服务器
	s.dnsServers = dnsServers
	if tr, ok := s.client.Transport.(*http.Transport); ok {
		tr.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			if len(s.dnsServers) > 0 {
				return customDial(ctx, network, address, s.dnsServers)
			}
			dialer := &net.Dialer{
				Timeout:   dnsResloveTimeout,
				DualStack: true,
			}
			return dialer.DialContext(ctx, network, address)
		}
	}
	return s
}

// customDial 使用自定义 DNS 服务器解析并连接目标地址
func customDial(ctx context.Context, network, address string, dnsServers []string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   dnsResloveTimeout,
		DualStack: true,
	}
	// 创建自定义的 Resolver
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
	return nil, fmt.Errorf("failed to connect to any resolved IPs for %s", address)
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

func (s *defaultSession) WithHeader(key, value string) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.headers.Set(key, value)
	return s
}

// WithHTTP2 设置当前 session 是否使用 HTTP/2
func (s *defaultSession) WithHTTP2(enabled bool) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	if s.useHTTP2 == enabled {
		return s
	}
	s.useHTTP2 = enabled
	if tr, ok := s.client.Transport.(*http.Transport); ok {
		// 将当前 Transport 放回相应的池中
		if tr.ForceAttemptHTTP2 {
			http2TransportPool.Put(tr)
		} else {
			http1TransportPool.Put(tr)
		}
		// 获取新的 Transport
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

func (s *defaultSession) Clone() Client {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	// 获取新的 Transport
	var transport *http.Transport
	if s.useHTTP2 {
		transport = http2TransportPool.Get().(*http.Transport)
	} else {
		transport = http1TransportPool.Get().(*http.Transport)
	}
	newClient := &http.Client{
		Transport: transport,
	}
	newHeaders := s.headers.Clone()
	newSession := &defaultSession{
		baseURL:      s.baseURL,
		timeout:      s.timeout,
		proxyURL:     s.proxyURL,
		headers:      newHeaders,
		client:       newClient,
		useHTTP2:     s.useHTTP2,
		keepAlive:    s.keepAlive,
		maxIdleConns: s.maxIdleConns,
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

// Clear 清空当前 Session 的所有配置并恢复到默认状态
func (s *defaultSession) Clear() Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.baseURL = ""
	s.useHTTP2 = false
	s.keepAlive = true
	s.maxIdleConns = defaultMaxIdleConns
	s.timeout = idleConnTimeout
	s.dnsServers = nil
	s.client = &http.Client{
		Transport: GetTransport(false), // 重置为默认 Transport
	}
	return s
}

func (s *defaultSession) Do(req *Request) (*models.Response, error) {
	if req.URL == nil {
		return nil, errors.New("request URL cannot be nil")
	}
	finalURL := req.URL
	if s.baseURL != "" {
		base, err := url.Parse(s.baseURL)
		if err != nil {
			return nil, err
		}
		finalURL = base.ResolveReference(req.URL)
	}
	httpReq, err := http.NewRequest(req.Method, finalURL.String(), req.Body)
	if err != nil {
		return nil, err
	}
	for k, vals := range s.headers {
		for _, v := range vals {
			httpReq.Header.Add(k, v)
		}
	}
	for k, vals := range req.Headers {
		for _, v := range vals {
			httpReq.Header.Add(k, v)
		}
	}
	ctx := httpReq.Context()
	if s.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
		httpReq = httpReq.WithContext(ctx)
	}
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	return models.NewResponse(resp, finalURL.String())
}
