package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sunerpy/requests"
	"github.com/sunerpy/requests/url"
)

func BenchmarkHTTPLibraries(b *testing.B) {
	// 创建一个模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟一个简单的响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close() // 测试结束时关闭模拟服务器
	benchmarks := []struct {
		name     string
		testFunc func(b *testing.B)
	}{
		{
			name: "NetHTTP",
			testFunc: func(b *testing.B) {
				client := &http.Client{}
				for i := 0; i < b.N; i++ {
					// 使用模拟服务器的地址
					req, err := http.NewRequest("GET", server.URL+"/get?key=value", nil)
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
			},
		},
		{
			name: "Requests",
			testFunc: func(b *testing.B) {
				params := url.NewURLParams()
				params.Set("key", "value")
				session := requests.NewSession()
				for i := 0; i < b.N; i++ {
					// 使用模拟服务器的地址
					req, err := requests.NewRequest("GET", server.URL+"/get", params, nil)
					if err != nil {
						b.Fatalf("requests request creation failed: %v", err)
					}
					resp, err := session.Do(req)
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
				session := requests.NewSession()
				params := url.NewURLParams()
				params.Set("key", "value")
				for i := 0; i < b.N; i++ {
					// 使用模拟服务器的地址
					req, err := requests.NewRequest("GET", server.URL+"/get", params, nil)
					if err != nil {
						b.Fatalf("requests request creation failed: %v", err)
					}
					resp, err := session.Do(req)
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
				session := requests.NewSession().WithMaxIdleConns(300)
				params := url.NewURLParams()
				params.Set("key", "value")
				for i := 0; i < b.N; i++ {
					// 使用模拟服务器的地址
					req, err := requests.NewRequest("GET", server.URL+"/get", params, nil)
					if err != nil {
						b.Fatalf("requests request creation failed: %v", err)
					}
					resp, err := session.Do(req)
					if err != nil {
						b.Fatalf("requests HTTP/2 GET failed: %v", err)
					}
					resp.Bytes()
				}
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, bm.testFunc)
	}
}

func BenchmarkHTTPLibrariesParallel(b *testing.B) {
	// 创建一个模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟一个简单的响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close() // 测试结束时关闭模拟服务器
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
						// 使用模拟服务器的地址
						req, err := http.NewRequest("GET", server.URL+"/get?key=value", nil)
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
						// 使用模拟服务器的地址
						resp, err := requests.Get(server.URL+"/get", params)
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
						// 使用模拟服务器的地址
						resp, err := requests.Get(server.URL+"/get", params)
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
					session := requests.NewSession().WithMaxIdleConns(300)
					params := url.NewURLParams()
					params.Set("key", "value")
					uri := server.URL + "/get"
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
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, bm.testFunc)
	}
}
