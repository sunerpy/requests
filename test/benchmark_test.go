package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sunerpy/requests"
	"github.com/sunerpy/requests/url"
)

// ============================================================================
// Test Server Setup
// ============================================================================
// newTestServer creates a test server that returns JSON responses
func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message": "test response"}`))
	}))
}

// TestResponse is the response type for JSON benchmarks
type TestResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ============================================================================
// Basic HTTP Request Benchmarks
// ============================================================================
func BenchmarkBasicRequests(b *testing.B) {
	server := newTestServer()
	defer server.Close()
	b.Run("NetHTTP_Get", func(b *testing.B) {
		b.ReportAllocs()
		client := &http.Client{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", server.URL+"/get", nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	})
	b.Run("Requests_Get", func(b *testing.B) {
		b.ReportAllocs()
		session := requests.NewSession()
		defer session.Close()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := requests.NewGet(server.URL + "/get").Build()
			resp, err := session.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Bytes()
		}
	})
	b.Run("NetHTTP_Get_WithQuery", func(b *testing.B) {
		b.ReportAllocs()
		client := &http.Client{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", server.URL+"/get?key=value", nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	})
	b.Run("Requests_Get_WithQuery", func(b *testing.B) {
		b.ReportAllocs()
		session := requests.NewSession()
		defer session.Close()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := requests.NewGet(server.URL+"/get").
				WithQuery("key", "value").
				Build()
			resp, err := session.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Bytes()
		}
	})
	b.Run("NetHTTP_Get_PrebuiltURL", func(b *testing.B) {
		b.ReportAllocs()
		client := &http.Client{}
		targetURL := server.URL + "/get"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", targetURL, nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	})
	b.Run("Requests_Get_PrebuiltURL", func(b *testing.B) {
		b.ReportAllocs()
		session := requests.NewSession()
		defer session.Close()
		// Pre-build the URL once
		targetURL := server.URL + "/get"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := requests.NewGet(targetURL).Build()
			resp, err := session.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Bytes()
		}
	})
}

// ============================================================================
// JSON Parsing Benchmarks
// ============================================================================
func BenchmarkJSONParsing(b *testing.B) {
	server := newTestServer()
	defer server.Close()
	b.Run("NetHTTP_GetJSON", func(b *testing.B) {
		b.ReportAllocs()
		client := &http.Client{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", server.URL+"/get", nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			var result TestResponse
			json.Unmarshal(body, &result)
			_ = result
		}
	})
	b.Run("Requests_GetJSON", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := requests.GetJSON[TestResponse](server.URL + "/get")
			if err != nil {
				b.Fatal(err)
			}
			_ = result.Data()
		}
	})
	b.Run("NetHTTP_GetJSON_WithQuery", func(b *testing.B) {
		b.ReportAllocs()
		client := &http.Client{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", server.URL+"/get?key=value", nil)
			req.Header.Set("X-Custom", "header")
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			var result TestResponse
			json.Unmarshal(body, &result)
			_ = result
		}
	})
	b.Run("Requests_GetJSON_WithOptions", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := requests.GetJSON[TestResponse](
				server.URL+"/get",
				requests.WithQuery("key", "value"),
				requests.WithHeader("X-Custom", "header"),
			)
			if err != nil {
				b.Fatal(err)
			}
			_ = result.Data()
		}
	})
	b.Run("NetHTTP_Manual_Parse", func(b *testing.B) {
		b.ReportAllocs()
		client := &http.Client{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", server.URL+"/get", nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			var result TestResponse
			json.Unmarshal(body, &result)
		}
	})
	b.Run("Requests_Manual_Parse", func(b *testing.B) {
		b.ReportAllocs()
		session := requests.NewSession()
		defer session.Close()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := requests.NewGet(server.URL + "/get").Build()
			resp, err := session.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			var result TestResponse
			json.Unmarshal(resp.Bytes(), &result)
		}
	})
}

// ============================================================================
// Session Benchmarks
// ============================================================================
func BenchmarkSession(b *testing.B) {
	server := newTestServer()
	defer server.Close()
	b.Run("NetHTTP_Client_Create", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			client := &http.Client{}
			_ = client
		}
	})
	b.Run("Requests_Session_Create", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			session := requests.NewSession()
			session.Close()
		}
	})
	b.Run("Requests_Session_WithConfig", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			session := requests.NewSession().
				WithBaseURL(server.URL).
				WithHeader("X-Custom", "value").
				WithMaxIdleConns(100)
			session.Close()
		}
	})
	b.Run("NetHTTP_Client_Reuse", func(b *testing.B) {
		b.ReportAllocs()
		client := &http.Client{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", server.URL+"/get", nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	})
	b.Run("Requests_Session_Reuse", func(b *testing.B) {
		b.ReportAllocs()
		session := requests.NewSession()
		defer session.Close()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := requests.NewGet(server.URL + "/get").Build()
			resp, err := session.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Bytes()
		}
	})
}

// ============================================================================
// RequestBuilder Benchmarks
// ============================================================================
func BenchmarkRequestBuilder(b *testing.B) {
	b.Run("NetHTTP_NewRequest_Simple", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = http.NewRequest("GET", "https://example.com/api", nil)
		}
	})
	b.Run("Requests_Builder_Simple", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = requests.NewGet("https://example.com/api").Build()
		}
	})
	b.Run("NetHTTP_NewRequest_WithQuery", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = http.NewRequest("GET", "https://example.com/api?key1=value1&key2=value2", nil)
		}
	})
	b.Run("Requests_Builder_WithQuery", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = requests.NewGet("https://example.com/api").
				WithQuery("key1", "value1").
				WithQuery("key2", "value2").
				Build()
		}
	})
	b.Run("NetHTTP_NewRequest_WithHeaders", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", "https://example.com/api", nil)
			req.Header.Set("X-Header-1", "value1")
			req.Header.Set("X-Header-2", "value2")
			req.Header.Set("Authorization", "Bearer token")
		}
	})
	b.Run("Requests_Builder_WithHeaders", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = requests.NewGet("https://example.com/api").
				WithHeader("X-Header-1", "value1").
				WithHeader("X-Header-2", "value2").
				WithHeader("Authorization", "Bearer token").
				Build()
		}
	})
	b.Run("NetHTTP_NewRequest_Complex", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("POST", "https://example.com/api?version=v1", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer my-token")
		}
	})
	b.Run("Requests_Builder_Complex", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = requests.NewPost("https://example.com/api").
				WithQuery("version", "v1").
				WithHeader("Content-Type", "application/json").
				WithHeader("Accept", "application/json").
				WithBearerToken("my-token").
				WithJSON(map[string]string{"key": "value"}).
				Build()
		}
	})
}

// ============================================================================
// Result[T] Benchmarks
// ============================================================================
func BenchmarkResultType(b *testing.B) {
	server := newTestServer()
	defer server.Close()
	b.Run("Result_Access", func(b *testing.B) {
		b.ReportAllocs()
		result, _ := requests.GetJSON[TestResponse](server.URL + "/get")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = result.Data()
			_ = result.StatusCode()
			_ = result.IsSuccess()
			_ = result.Headers()
		}
	})
}

// ============================================================================
// HTTP/2 Benchmarks
// ============================================================================
func BenchmarkHTTP2(b *testing.B) {
	server := newTestServer()
	defer server.Close()
	b.Run("NetHTTP_Default", func(b *testing.B) {
		b.ReportAllocs()
		client := &http.Client{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", server.URL+"/get", nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	})
	b.Run("Requests_HTTP1", func(b *testing.B) {
		b.ReportAllocs()
		requests.SetHTTP2Enabled(false)
		session := requests.NewSession()
		defer session.Close()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := requests.NewGet(server.URL + "/get").Build()
			resp, err := session.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Bytes()
		}
	})
	b.Run("Requests_HTTP2", func(b *testing.B) {
		b.ReportAllocs()
		requests.SetHTTP2Enabled(true)
		session := requests.NewSession()
		defer session.Close()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := requests.NewGet(server.URL + "/get").Build()
			resp, err := session.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Bytes()
		}
	})
	b.Run("Requests_HTTP2_WithMaxIdleConns", func(b *testing.B) {
		b.ReportAllocs()
		requests.SetHTTP2Enabled(true)
		session := requests.NewSession().WithMaxIdleConns(300)
		defer session.Close()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := requests.NewGet(server.URL + "/get").Build()
			resp, err := session.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Bytes()
		}
	})
}

// ============================================================================
// Parallel Benchmarks
// ============================================================================
func BenchmarkParallel(b *testing.B) {
	server := newTestServer()
	defer server.Close()
	b.Run("NetHTTP_Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			client := &http.Client{}
			for pb.Next() {
				req, _ := http.NewRequest("GET", server.URL+"/get", nil)
				resp, err := client.Do(req)
				if err != nil {
					b.Fatal(err)
				}
				io.ReadAll(resp.Body)
				resp.Body.Close()
			}
		})
	})
	b.Run("Requests_Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				resp, err := requests.Get(server.URL+"/get", requests.WithQuery("key", "value"))
				if err != nil {
					b.Fatal(err)
				}
				resp.Bytes()
			}
		})
	})
	b.Run("Requests_Parallel_Session", func(b *testing.B) {
		b.ReportAllocs()
		session := requests.NewSession().WithMaxIdleConns(300)
		defer session.Close()
		// Pre-build URL with params (use unsafe version for single-threaded URL building)
		params := url.NewURLParamsUnsafe()
		params.Set("key", "value")
		uri, _ := url.FastBuildURL(server.URL+"/get", params)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				req, _ := requests.NewGet(uri).Build()
				resp, err := session.Do(req)
				if err != nil {
					b.Fatal(err)
				}
				resp.Bytes()
			}
		})
	})
}

// ============================================================================
// URL Building Benchmarks
// ============================================================================
func BenchmarkURLBuilding(b *testing.B) {
	b.Run("NetHTTP_URLParse", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			u := "https://example.com/api?key1=value1&key2=value2&key3=value3"
			_, _ = http.NewRequest("GET", u, nil)
		}
	})
	b.Run("Requests_URLParams_Build", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Use unsafe version for maximum performance in single-threaded benchmark
			params := url.NewURLParamsUnsafe()
			params.Set("key1", "value1")
			params.Set("key2", "value2")
			params.Set("key3", "value3")
			url.FastBuildURL("https://example.com/api", params)
		}
	})
	b.Run("Requests_URLParams_Build_Safe", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Use thread-safe version
			params := url.NewURLParams()
			params.Set("key1", "value1")
			params.Set("key2", "value2")
			params.Set("key3", "value3")
			url.BuildURL("https://example.com/api", params)
		}
	})
}
