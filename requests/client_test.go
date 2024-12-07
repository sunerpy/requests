package requests

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockTransport struct {
	Response *http.Response
	Err      error
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
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
	session := NewSession().WithBaseURL("https://example.com").WithTimeout(5 * time.Second)
	clonedSession := session.Clone()
	clonedSess, ok := clonedSession.(*defaultSession)
	assert.True(t, ok)
	assert.Equal(t, "https://example.com", clonedSess.baseURL)
	assert.Equal(t, 5*time.Second, clonedSess.timeout)
}

func TestDefaultSession_Do(t *testing.T) {
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
}
