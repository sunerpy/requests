package test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	thirdrequest "github.com/gospider007/requests"
	"github.com/sunerpy/requests"
	"github.com/sunerpy/requests/url"
)

func BenchmarkHTTPLibraries(b *testing.B) {
	benchmarks := []struct {
		name     string
		testFunc func(b *testing.B)
	}{
		{
			name: "NetHTTP",
			testFunc: func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					req, err := http.NewRequest("GET", "http://192.168.82.149:22222/get?key=value", nil)
					if err != nil {
						b.Fatalf("net/http request creation failed: %v", err)
					}
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						b.Fatalf("net/http GET failed: %v", err)
					}
					_, err = io.ReadAll(resp.Body)
					if err != nil {
						b.Fatalf("Failed to read response body: %v", err)
					}
					resp.Body.Close()
				}
			},
		},
		{
			name: "Requests",
			testFunc: func(b *testing.B) {
				params := url.NewURLParams()
				params.Set("key", "value")
				for i := 0; i < b.N; i++ {
					resp, err := requests.Get("http://192.168.82.149:22222/get", params)
					if err != nil {
						b.Fatalf("requests GET failed: %v", err)
					}
					resp.Bytes()
				}
			},
		},
		{
			name: "Requests_HTTP2",
			testFunc: func(b *testing.B) {
				requests.SetHTTP2Enabled(true)
				params := url.NewURLParams()
				params.Set("key", "value")
				for i := 0; i < b.N; i++ {
					resp, err := requests.Get("http://192.168.82.149:22222/get", params)
					if err != nil {
						b.Fatalf("requests HTTP/2 GET failed: %v", err)
					}
					resp.Bytes()
				}
			},
		},
		{
			name: "Requests_HTTP2_withmax",
			testFunc: func(b *testing.B) {
				requests.SetHTTP2Enabled(true)
				requests.NewSession().WithMaxIdleConns(300)
				params := url.NewURLParams()
				params.Set("key", "value")
				for i := 0; i < b.N; i++ {
					resp, err := requests.Get("http://192.168.82.149:22222/get", params)
					if err != nil {
						b.Fatalf("requests HTTP/2 GET failed: %v", err)
					}
					resp.Bytes()
				}
			},
		},
		{
			name: "ThirdRequest",
			testFunc: func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					resp, err := thirdrequest.Get(nil, "http://192.168.82.149:22222/get", thirdrequest.RequestOption{Params: map[string]string{"key": "value"}})
					if err != nil {
						b.Fatalf("thirdrequest GET failed: %v", err)
					}
					resp.Content()
				}
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, bm.testFunc)
	}
}

func BenchmarkHTTPLibrariesParallel(b *testing.B) {
	benchmarks := []struct {
		name     string
		testFunc func(b *testing.B)
	}{
		{
			name: "NetHTTP",
			testFunc: func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					client := &http.Client{}
					for pb.Next() {
						req, err := http.NewRequest("GET", "http://192.168.82.149:22222/get?key=value", nil)
						if err != nil {
							b.Fatalf("net/http request creation failed: %v", err)
						}
						resp, err := client.Do(req)
						if err != nil {
							b.Fatalf("net/http GET failed: %v", err)
						}
						_, err = io.ReadAll(resp.Body)
						if err != nil {
							b.Fatalf("Failed to read response body: %v", err)
						}
						resp.Body.Close()
					}
				})
			},
		},
		{
			name: "Requests",
			testFunc: func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					params := url.NewURLParams()
					params.Set("key", "value")
					for pb.Next() {
						resp, err := requests.Get("http://192.168.82.149:22222/get", params)
						if err != nil {
							b.Fatalf("requests GET failed: %v", err)
						}
						resp.Bytes()
					}
				})
			},
		},
		{
			name: "Requests_HTTP2",
			testFunc: func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					requests.SetHTTP2Enabled(true)
					params := url.NewURLParams()
					params.Set("key", "value")
					for pb.Next() {
						resp, err := requests.Get("http://192.168.82.149:22222/get", params)
						if err != nil {
							b.Fatalf("requests HTTP/2 GET failed: %v", err)
						}
						resp.Bytes()
					}
				})
			},
		},
		{
			name: "Requests_HTTP2_withmax",
			testFunc: func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					requests.SetHTTP2Enabled(true)
					session := requests.NewSession().WithMaxIdleConns(500)
					params := url.NewURLParams()
					params.Set("key", "value")
					uri := "http://192.168.82.149:22222/get"
					u, err := url.BuildURL(uri, params)
					if err != nil {
						b.Fatalf("url.BuildURL failed: %v", err)
					}
					for pb.Next() {
						req, err := requests.NewRequest("GET", u, params, nil)
						if err != nil {
							b.Fatalf("request creation failed: %v", err)
						}
						resp, err := session.Do(req)
						if err != nil {
							b.Fatalf("requests HTTP/2 GET failed: %v", err)
						}
						resp.Bytes()
					}
				})
			},
		},
		{
			name: "ThirdRequest",
			testFunc: func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						resp, err := thirdrequest.Get(nil, "http://192.168.82.149:22222/get", thirdrequest.RequestOption{Params: map[string]string{"key": "value"}})
						if err != nil {
							b.Fatalf("thirdrequest GET failed: %v", err)
						}
						resp.Content()
					}
				})
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, bm.testFunc)
	}
}

func BenchmarkLargeDataTransfer(b *testing.B) {
	testCases := []struct {
		size int
		name string
	}{
		{1 * 1024 * 1024, "1MB File"},
		{10 * 1024 * 1024, "10MB File"},
	}
	for _, tc := range testCases {
		b.Run(tc.name+"/NetHTTP", func(b *testing.B) {
			content := strings.Repeat("a", tc.size)
			reader := strings.NewReader(content)
			for i := 0; i < b.N; i++ {
				req, _ := http.NewRequest("POST", "http://192.168.82149:22222/post", reader)
				client := &http.Client{}
				resp, _ := client.Do(req)
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
		})
		b.Run(tc.name+"/Requests", func(b *testing.B) {
			content := strings.Repeat("a", tc.size)
			form := url.NewForm()
			form.Set("file", content)
			for i := 0; i < b.N; i++ {
				resp, _ := requests.Post("http://192.168.82149:22222/post", form)
				resp.Bytes()
			}
		})
	}
}

func BenchmarkLargeHeaders(b *testing.B) {
	headers := map[string]string{
		"X-Test-Header-1": "Value1",
		"X-Test-Header-2": "Value2",
	}
	content := strings.Repeat("test", 1024)
	b.Run("NetHTTP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("POST", "http://192.168.82149:22222/post", strings.NewReader(content))
			for key, value := range headers {
				req.Header.Set(key, value)
			}
			client := &http.Client{}
			resp, _ := client.Do(req)
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	})
	b.Run("Requests", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			form := url.NewForm()
			form.Set("key", content)
			session := requests.NewSession().WithHTTP2(true)
			for key, value := range headers {
				session.WithHeader(key, value)
			}
			req, _ := requests.NewRequest("POST", "http://192.168.82149:22222/post", nil, strings.NewReader(content))
			resp, _ := session.Do(req)
			resp.Bytes()
		}
	})
}
