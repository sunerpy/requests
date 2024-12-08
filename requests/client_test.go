package requests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

type MockTransport struct {
	Response *http.Response
	Err      error
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestGetTransport(t *testing.T) {
	transport := GetTransport(true)
	assert.NotNil(t, transport, "GetTransport(true) should not return nil")
	assert.True(t, transport.ForceAttemptHTTP2, "GetTransport(true) should return an HTTP/2 transport")
	transport = GetTransport(false)
	assert.NotNil(t, transport, "GetTransport(false) should not return nil")
	assert.False(t, transport.ForceAttemptHTTP2, "GetTransport(false) should return an HTTP/1 transport")
}

func TestPutTransport(t *testing.T) {
	transport := GetTransport(true)
	PutTransport(transport)
	assert.Equal(t, transport, http2TransportPool.Get(), "PutTransport did not put HTTP/2 transport back in the pool")
	transport = GetTransport(false)
	PutTransport(transport)
	assert.Equal(t, transport, http1TransportPool.Get(), "PutTransport did not put HTTP/1 transport back in the pool")
}

func TestSetHTTP2Enabled(t *testing.T) {
	SetHTTP2Enabled(true)
	assert.True(t, defaultHTTP2Enabled)
	SetHTTP2Enabled(false)
	assert.False(t, defaultHTTP2Enabled)
	assert.False(t, IsHTTP2Enabled())
}

func TestNewSession(t *testing.T) {
	SetHTTP2Enabled(true)
	session := NewSession()
	assert.NotNil(t, session)
	defaultSess, ok := session.(*defaultSession)
	assert.True(t, ok)
	assert.True(t, defaultSess.useHTTP2)
	SetHTTP2Enabled(false)
	session = NewSession()
	assert.NotNil(t, session)
	defaultSess, ok = session.(*defaultSession)
	assert.True(t, ok)
	assert.False(t, defaultSess.useHTTP2)
}

func TestDefaultSession_WithBaseURL(t *testing.T) {
	session := NewSession().WithBaseURL("https://example.com")
	defaultSess, ok := session.(*defaultSession)
	assert.True(t, ok)
	assert.Equal(t, "https://example.com", defaultSess.baseURL)
}

func TestDefaultSession_WithTimeout(t *testing.T) {
	session := NewSession().WithTimeout(5 * time.Second)
	defaultSess, ok := session.(*defaultSession)
	assert.True(t, ok)
	assert.Equal(t, 5*time.Second, defaultSess.timeout)
}

func TestDefaultSession_WithProxy(t *testing.T) {
	proxyURL := "http://localhost:8080"
	session := NewSession().WithProxy(proxyURL)
	defaultSess, ok := session.(*defaultSession)
	assert.True(t, ok)
	assert.NotNil(t, defaultSess.proxyURL)
	assert.Equal(t, proxyURL, defaultSess.proxyURL.String())
}

func TestDefaultSession_WithProxy_InvalidURL(t *testing.T) {
	session := NewSession().WithProxy("invalid-url")
	defaultSess, ok := session.(*defaultSession)
	assert.True(t, ok)
	assert.Nil(t, defaultSess.proxyURL)
}

type mockConn struct{}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 12345,
	}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("8.8.8.8"),
		Port: 53,
	}
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestWithDNS_GoMonkey(t *testing.T) {
	dnsServers := []string{"8.8.8.8", "8.8.4.4"}
	sess := &defaultSession{
		client: &http.Client{
			Transport: &http.Transport{},
		},
	}
	patches := gomonkey.ApplyFunc((*net.Resolver).LookupHost, func(_ *net.Resolver, ctx context.Context, host string) ([]string, error) {
		if host == "example.com" {
			return []string{"93.184.216.34"}, nil
		}
		return nil, fmt.Errorf("DNS resolution failed")
	})
	defer patches.Reset()
	patches.ApplyMethod(reflect.TypeOf(&net.Dialer{}), "DialContext", func(_ *net.Dialer, ctx context.Context, network, address string) (net.Conn, error) {
		if address == "8.8.8.8:53" || address == "8.8.4.4:53" {
			return &mockConn{}, nil
		}
		if address == "93.184.216.34:80" {
			return &mockConn{}, nil
		}
		return nil, fmt.Errorf("failed to connect")
	})
	sess.WithDNS(dnsServers)
	tr, ok := sess.client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", sess.client.Transport)
	}
	ctx := context.Background()
	network := "tcp"
	address := "example.com:80"
	conn, err := tr.DialContext(ctx, network, address)
	if err != nil {
		t.Fatalf("DialContext failed: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil connection, got nil")
	}
	t.Log("DialContext succeeded with custom DNS")
}

func TestDefaultSession_WithHeader(t *testing.T) {
	session := NewSession().WithHeader("Authorization", "Bearer token")
	defaultSess, ok := session.(*defaultSession)
	assert.True(t, ok)
	assert.Equal(t, "Bearer token", defaultSess.headers.Get("Authorization"))
}

func TestDefaultSession_WithHTTP2(t *testing.T) {
	session := NewSession()
	defaultSess, ok := session.(*defaultSession)
	assert.True(t, ok)
	assert.False(t, defaultSess.useHTTP2)
	session.WithHTTP2(true)
	assert.True(t, defaultSess.useHTTP2)
	session.WithHTTP2(true)
	assert.True(t, defaultSess.useHTTP2)
	session.WithHTTP2(false)
	assert.False(t, defaultSess.useHTTP2)
}

func TestDefaultSession_Clone(t *testing.T) {
	t.Run("Clone test", func(t *testing.T) {
		session := NewSession().WithBaseURL("https://example.com").WithTimeout(5 * time.Second)
		clonedSession := session.Clone()
		clonedSess, ok := clonedSession.(*defaultSession)
		assert.True(t, ok)
		assert.Equal(t, "https://example.com", clonedSess.baseURL)
		assert.Equal(t, 5*time.Second, clonedSess.timeout)
	})
	t.Run("Clone with http2 ", func(t *testing.T) {
		session := NewSession().WithBaseURL("https://example.com").WithTimeout(5 * time.Second).WithHTTP2(true)
		clonedSession := session.Clone()
		clonedSess, ok := clonedSession.(*defaultSession)
		assert.True(t, ok)
		assert.True(t, clonedSess.useHTTP2)
		assert.Equal(t, session.(*defaultSession).client.Transport, clonedSess.client.Transport)
	})
}

func TestDefaultSession_Do(t *testing.T) {
	t.Run("Do test", func(t *testing.T) {
		mockTransport := &MockTransport{
			Response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"key": "value"}`))),
			},
		}
		session := &defaultSession{
			client: &http.Client{Transport: mockTransport},
			headers: http.Header{
				"Authorization": []string{"Bearer token"},
			},
		}
		req := &Request{
			Method:  "GET",
			URL:     mustParseURL("https://example.com"),
			Headers: http.Header{"Custom-Header": []string{"CustomValue"}},
		}
		resp, err := session.Do(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, `{"key": "value"}`, resp.Text())
	})
	t.Run("Do with set session baseurl", func(t *testing.T) {
		session := NewSession().WithBaseURL("https://example.com").WithTimeout(5 * time.Second)
		req, err := NewRequest("GET", "path", nil, nil)
		assert.NoError(t, err)
		resp, err := session.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/path", resp.GetURL())
	})
	t.Run("Do with error baseurl", func(t *testing.T) {
		session := NewSession().WithBaseURL("123.3:{}invalid-url")
		req, err := NewRequest("GET", "path", nil, nil)
		assert.NoError(t, err)
		resp, err := session.Do(req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestDefaultSession_Do_WithTimeout(t *testing.T) {
	mockTransport := &MockTransport{
		Err: context.DeadlineExceeded,
	}
	session := &defaultSession{
		client:  &http.Client{Transport: mockTransport},
		timeout: 1 * time.Millisecond,
	}
	req := &Request{
		Method: "GET",
		URL:    mustParseURL("https://example.com"),
	}
	resp, err := session.Do(req)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
}

func TestDefaultSession_Do_WithBaseURL(t *testing.T) {
	mockTransport := &MockTransport{
		Response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`OK`))),
		},
	}
	session := &defaultSession{
		client:  &http.Client{Transport: mockTransport},
		baseURL: "https://example.com",
	}
	req := &Request{
		Method: "GET",
		URL:    mustParseURL("/path"),
	}
	resp, err := session.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "OK", resp.Text())
}

func TestDefaultSession_Do_InvalidURL(t *testing.T) {
	session := NewSession()
	req := &Request{
		Method: "GET",
		URL:    nil,
	}
	resp, err := session.Do(req)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func TestDefaultSession_WithKeepAlive(t *testing.T) {
	t.Run("Set KeepAlive to true", func(t *testing.T) {
		session := NewSession()
		defaultSess, ok := session.(*defaultSession)
		assert.True(t, ok)
		session.WithKeepAlive(true)
		assert.True(t, defaultSess.keepAlive)
		if tr, ok := defaultSess.client.Transport.(*http.Transport); ok {
			assert.False(t, tr.DisableKeepAlives)
		}
	})
	t.Run("Set KeepAlive to false", func(t *testing.T) {
		session := NewSession()
		defaultSess, ok := session.(*defaultSession)
		assert.True(t, ok)
		session.WithKeepAlive(false)
		assert.False(t, defaultSess.keepAlive)
		if tr, ok := defaultSess.client.Transport.(*http.Transport); ok {
			assert.True(t, tr.DisableKeepAlives)
		}
	})
	t.Run("Set transport nil", func(t *testing.T) {
		session := NewSession()
		defaultSess, ok := session.(*defaultSession)
		assert.True(t, ok)
		defaultSess.client.Transport = nil
		session.WithKeepAlive(true)
		assert.True(t, defaultSess.keepAlive)
	})
}

func TestDefaultSession_WithMaxIdleConns(t *testing.T) {
	t.Run("Set MaxIdleConns", func(t *testing.T) {
		session := NewSession()
		defaultSess, ok := session.(*defaultSession)
		assert.True(t, ok)
		session.WithMaxIdleConns(10)
		assert.Equal(t, 10, defaultSess.maxIdleConns)
		if tr, ok := defaultSess.client.Transport.(*http.Transport); ok {
			assert.Equal(t, 10, tr.MaxIdleConns)
			assert.Equal(t, 10, tr.MaxIdleConnsPerHost)
		}
	})
	t.Run("Set transport nil", func(t *testing.T) {
		session := NewSession()
		defaultSess, ok := session.(*defaultSession)
		assert.True(t, ok)
		defaultSess.client.Transport = nil
		session.WithMaxIdleConns(10)
		assert.Equal(t, 10, defaultSess.maxIdleConns)
	})
}

func TestDefaultSession_Close(t *testing.T) {
	t.Run("Close and release Transport", func(t *testing.T) {
		session := NewSession()
		defaultSess, ok := session.(*defaultSession)
		assert.True(t, ok)
		tr := &http.Transport{}
		defaultSess.client.Transport = tr
		err := session.Close()
		assert.NoError(t, err)
		assert.Nil(t, defaultSess.client.Transport)
	})
	t.Run("Close with enable http2", func(t *testing.T) {
		session := NewSession()
		defaultSess, ok := session.(*defaultSession)
		assert.True(t, ok)
		defaultSess.WithHTTP2(true)
		err := session.Close()
		assert.NoError(t, err)
	})
}

func TestClearSession(t *testing.T) {
	t.Run("Clear session", func(t *testing.T) {
		session := NewSession()
		session.WithBaseURL("https://example.com")
		session.Clear()
		assert.Equal(t, "", session.(*defaultSession).baseURL)
	})
}
