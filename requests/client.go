package requests

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/sunerpy/requests/models"
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
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
				DisableKeepAlives:   false,
				ForceAttemptHTTP2:   false,
				TLSNextProto:        make(map[string]func(string, *tls.Conn) http.RoundTripper),
			}
		},
	}
	http2TransportPool = sync.Pool{
		New: func() interface{} {
			return &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
				DisableKeepAlives:   false,
				ForceAttemptHTTP2:   true,
				// TLSNextProto 保持默认，支持 HTTP/2
			}
		},
	}
	defaultSess = NewSession()
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
	WithHeader(key, value string) Session
	WithHTTP2(enabled bool) Session
	WithKeepAlive(enabled bool) Session
	WithMaxIdleConns(maxIdle int) Session
	Close() error
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
		maxIdleConns: 100,
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
	if tr, ok := s.client.Transport.(*http.Transport); ok {
		tr.DisableKeepAlives = !enabled
	}
	return s
}

func (s *defaultSession) WithMaxIdleConns(maxIdle int) Session {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.maxIdleConns = maxIdle
	if tr, ok := s.client.Transport.(*http.Transport); ok {
		tr.MaxIdleConns = maxIdle
		tr.MaxIdleConnsPerHost = maxIdle
	}
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
	// 创建新的 http.Client
	newClient := &http.Client{
		Transport: transport,
	}
	// 复制并深拷贝 headers
	newHeaders := s.headers.Clone()
	// 创建新的 Session
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
	// 如果不使用 HTTP/2，清空 TLSNextProto
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
		s.client.Transport = nil // 防止再次使用
	}
	return nil
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
