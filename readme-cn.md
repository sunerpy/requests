# Go Requests

[![Language](https://img.shields.io/badge/language-golang-brightgreen)](https://github.com/sunerpy/requests) [![Last Commit](https://img.shields.io/github/last-commit/sunerpy/requests)](https://github.com/sunerpy/requests) [![CI](https://github.com/sunerpy/requests/workflows/Go/badge.svg)](https://github.com/sunerpy/requests/actions) [![codecov](https://codecov.io/gh/sunerpy/requests/branch/main/graph/badge.svg)](https://codecov.io/gh/sunerpy/requests) [![Benchmark](https://github.com/sunerpy/requests/actions/workflows/benchmark.yml/badge.svg)](https://sunerpy.github.io/requests/dev/bench) [![Go Reference](https://pkg.go.dev/badge/github.com/sunerpy/requests.svg)](https://pkg.go.dev/github.com/sunerpy/requests) [![Go Report Card](https://goreportcard.com/badge/github.com/sunerpy/requests)](https://goreportcard.com/report/github.com/sunerpy/requests) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](readme.md) | 简体中文

Go Requests 是一个现代化、类型安全的 Go HTTP 客户端库，灵感来源于 Python 的 Requests 库。它利用 Go 1.18+ 泛型提供编译时类型安全和流畅的 API。

## 特性

* **类型安全的泛型方法**: `GetJSON[T]`、`PostJSON[T]` 等返回 `Result[T]` 自动解析响应
* **流畅的请求构建器**: 使用 `NewGet`、`NewPost` 等链式方法构建复杂请求
* **统一的 Result 类型**: `Result[T]` 包装解析后的数据和响应元数据
* **中间件支持**: 添加日志、认证、重试等横切关注点
* **Session 管理**: 连接池、Cookie 持久化、默认请求头
* **可插拔编解码器**: 内置 JSON、XML，支持自定义编解码器
* **重试机制**: 可配置的重试策略，支持指数退避
* **请求/响应钩子**: 观察生命周期事件用于日志和指标
* **文件上传**: 支持进度跟踪的多部分上传
* **HTTP/2 支持**: 自动 HTTP/2，回退到 HTTP/1.1
* **连接池优化**: 可配置的池大小和空闲超时
* **高性能**: 对象池、零拷贝优化，接近 net/http 原生性能

## 安装

```bash
go get github.com/sunerpy/requests
```

## 快速开始

### 简单 GET 请求

```go
import "github.com/sunerpy/requests"

// 基本 GET 请求
resp, err := requests.Get("https://api.github.com")
if err != nil {
    log.Fatal(err)
}
fmt.Println("状态:", resp.StatusCode)
fmt.Println("响应:", resp.Text())

// 带查询参数的 GET
resp, err := requests.Get("https://api.github.com/search/repos",
    requests.WithQuery("q", "golang"),
    requests.WithQuery("sort", "stars"),
)
```

### 泛型 JSON 响应与 Result[T]

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// GetJSON 返回 Result[T]，包装数据和响应元数据
result, err := requests.GetJSON[User]("https://api.example.com/user/1")
if err != nil {
    log.Fatal(err)
}

// 访问解析后的数据
fmt.Printf("用户: %+v\n", result.Data())

// 访问响应元数据
fmt.Printf("状态码: %d\n", result.StatusCode())
fmt.Printf("成功: %v\n", result.IsSuccess())
```

### POST JSON 请求

```go
type CreateUser struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponse struct {
    ID int `json:"id"`
}

// PostJSON 返回 Result[T]
result, err := requests.PostJSON[UserResponse](
    "https://api.example.com/users",
    CreateUser{Name: "John", Email: "john@example.com"},
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("创建的用户 ID: %d\n", result.Data().ID)
```

### 使用请求构建器

```go
// 使用便捷构造函数
req, err := requests.NewPost("https://api.example.com/users").
    WithHeader("X-Custom-Header", "value").
    WithQuery("version", "v2").
    WithJSON(map[string]string{"name": "John"}).
    WithTimeout(10 * time.Second).
    Build()

// 使用 DoJSON 自动解析 JSON
result, err := requests.DoJSON[UserResponse](
    requests.NewPost("https://api.example.com/users").
        WithJSON(CreateUser{Name: "John", Email: "john@example.com"}),
)
```

## 高级用法

### 带默认配置的 Session

```go
session := requests.NewSession().
    WithBaseURL("https://api.example.com").
    WithTimeout(30 * time.Second).
    WithHeader("Authorization", "Bearer token").
    WithHTTP2(true).
    WithMaxIdleConns(100)

defer session.Close()

req, _ := requests.NewGet("/users").Build()
resp, err := session.Do(req)
```

### 使用 Context 控制超时和取消

```go
import "context"

// 创建带超时的 context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

session := requests.NewSession()
defer session.Close()

req, _ := requests.NewGet("https://api.example.com/data").Build()

// 使用 context 执行请求 - 支持超时和取消
resp, err := session.DoWithContext(ctx, req)
if err != nil {
    if err == context.DeadlineExceeded {
        log.Println("请求超时")
    } else if err == context.Canceled {
        log.Println("请求被取消")
    }
    return
}
fmt.Println("响应:", resp.Text())
```

### Context 取消示例

```go
ctx, cancel := context.WithCancel(context.Background())

// 2 秒后取消请求
go func() {
    time.Sleep(2 * time.Second)
    cancel()
}()

session := requests.NewSession()
defer session.Close()

req, _ := requests.NewGet("https://api.example.com/slow-endpoint").Build()
resp, err := session.DoWithContext(ctx, req)
if err != nil {
    // 处理取消
    log.Printf("请求被取消: %v", err)
}
```

### 带重试策略的 Session

```go
session := requests.NewSession().
    WithBaseURL("https://api.example.com").
    WithTimeout(30 * time.Second).
    WithRetry(requests.RetryPolicy{
        MaxAttempts:     3,                        // 最大重试次数
        InitialInterval: 100 * time.Millisecond,   // 初始重试间隔
        MaxInterval:     5 * time.Second,          // 最大重试间隔
        Multiplier:      2.0,                      // 退避乘数
        RetryIf: func(resp *requests.Response, err error) bool {
            if err != nil {
                return true // 网络错误时重试
            }
            // 5xx 错误和 429 太多请求时重试
            return resp.StatusCode >= 500 || resp.StatusCode == 429
        },
    })

defer session.Close()

req, _ := requests.NewGet("/users").Build()
resp, err := session.Do(req) // 失败时自动重试
```

### 带中间件的 Session

```go
// 创建日志中间件
loggingMiddleware := requests.MiddlewareFunc(func(req *requests.Request, next requests.Handler) (*requests.Response, error) {
    start := time.Now()
    log.Printf("请求: %s %s", req.Method, req.URL)
    
    resp, err := next(req)
    
    duration := time.Since(start)
    if resp != nil {
        log.Printf("响应: %d 耗时 %v", resp.StatusCode, duration)
    }
    return resp, err
})

// 创建带中间件的 session
session := requests.NewSession().
    WithBaseURL("https://api.example.com").
    WithMiddleware(loggingMiddleware)

defer session.Close()

req, _ := requests.NewGet("/users").Build()
resp, err := session.Do(req) // 中间件会记录请求/响应
```

### 中间件

```go
// 日志中间件
chain := requests.NewMiddlewareChain()
chain.UseFunc(func(req *requests.Request, next requests.Handler) (*requests.Response, error) {
    log.Printf("请求: %s %s", req.Method, req.URL)
    resp, err := next(req)
    if resp != nil {
        log.Printf("响应: %d", resp.StatusCode)
    }
    return resp, err
})
```

### 重试策略

```go
policy := requests.ExponentialRetryPolicy(3, 100*time.Millisecond, 10*time.Second)
executor := requests.NewRetryExecutor(policy)

resp, err := executor.Execute(nil, func() (*requests.Response, error) {
    return requests.Get("https://api.example.com/data")
})
```

## API 参考

### HTTP 方法

| 函数 | 描述 |
| ------ | ------ |
| `Get(url, opts...)` | 发送 GET 请求 |
| `Post(url, body, opts...)` | 发送 POST 请求 |
| `Put(url, body, opts...)` | 发送 PUT 请求 |
| `Delete(url, opts...)` | 发送 DELETE 请求 |
| `GetJSON[T](url, opts...)` | GET 并解析 JSON，返回 `(Result[T], error)` |
| `PostJSON[T](url, data, opts...)` | POST JSON，返回 `(Result[T], error)` |
| `NewSession()` | 创建新的 Session |
| `DefaultSession()` | 获取默认 Session |

### Session 方法

| 方法 | 描述 |
| ------ | ------ |
| `Do(req)` | 执行请求 |
| `DoWithContext(ctx, req)` | 使用 context 执行请求（支持超时/取消） |
| `Clone()` | 创建 Session 副本 |
| `Close()` | 关闭 Session 并释放资源 |
| `Clear()` | 重置 Session 到默认状态 |

### Session 配置方法

| 方法 | 描述 |
| ------ | ------ |
| `WithBaseURL(url)` | 设置所有请求的基础 URL |
| `WithTimeout(duration)` | 设置请求超时 |
| `WithProxy(url)` | 设置代理 URL |
| `WithDNS(servers)` | 设置自定义 DNS 服务器 |
| `WithHeader(key, value)` | 添加默认请求头 |
| `WithHeaders(map)` | 添加多个默认请求头 |
| `WithBasicAuth(user, pass)` | 设置 Basic 认证 |
| `WithBearerToken(token)` | 设置 Bearer Token |
| `WithHTTP2(enabled)` | 启用/禁用 HTTP/2 |
| `WithKeepAlive(enabled)` | 启用/禁用 Keep-Alive |
| `WithMaxIdleConns(n)` | 设置最大空闲连接数 |
| `WithIdleTimeout(duration)` | 设置空闲连接超时 |
| `WithRetry(policy)` | 设置重试策略 |
| `WithMiddleware(m)` | 添加中间件 |
| `WithCookieJar(jar)` | 设置 Cookie Jar |

### Result[T] 方法

| 方法 | 描述 |
| ------ | ------ |
| `Data()` | 获取解析后的数据（类型 T） |
| `StatusCode()` | 获取 HTTP 状态码 |
| `Headers()` | 获取响应头 |
| `IsSuccess()` | 检查状态码是否为 2xx |
| `IsError()` | 检查状态码是否为 4xx 或 5xx |

### 请求选项

| 选项 | 描述 |
| ------ | ------ |
| `WithHeader(key, value)` | 添加请求头 |
| `WithQuery(key, value)` | 添加查询参数 |
| `WithBasicAuth(user, pass)` | 设置 Basic 认证 |
| `WithBearerToken(token)` | 设置 Bearer Token |
| `WithTimeout(duration)` | 设置超时 |

## 迁移指南

### 从 v1.x 到 v2.x

1. **Result[T] 类型**: 泛型方法现在返回 `(Result[T], error)` 而不是 `(T, *Response, error)`
2. **请求构建器**: 使用 `NewGet`、`NewPost` 等代替 `NewRequest(method, url, params, body)`
3. **删除的方法**: `*WithResponse` 方法已删除，使用 `Result[T]` 代替
4. **单一导入**: 只需导入 `github.com/sunerpy/requests`

## 性能优化

本库针对高性能场景进行了优化：

* **对象池**: `RequestBuilder`、`RequestConfig`、`FastValues` 对象池减少内存分配
* **Header 所有权转移**: 避免不必要的克隆操作
* **响应缓冲池**: 分级缓冲池 (4KB/32KB) 用于响应读取
* **接近 net/http**: 仅比原生 net/http 多约 6 次内存分配

```bash
BenchmarkBasicRequests/NetHTTP_Get     40μs, 69 allocs, 6.3KB
BenchmarkBasicRequests/Requests_Get    44μs, 75 allocs, 6.3KB  (仅多 6 次分配)
```

对于性能敏感场景，可使用 `FastValues` 构建 URL 参数：

```go
import "github.com/sunerpy/requests/url"

// 无锁 FastValues（仅限单线程场景）
fv := url.AcquireFastValues()
defer url.ReleaseFastValues(fv)
fv.Add("page", "1")
fv.Add("limit", "10")
```

详见 [性能优化指南](docs/guides/performance-optimization.md)。

## 文档

详细文档请参阅：

* [入门指南](docs/guides/getting-started.md)
* [API 参考](docs/guides/api-reference.md)
* [性能优化指南](docs/guides/performance-optimization.md)
* [迁移指南](docs/guides/migration-guide.md)
* [示例](docs/examples/)

## 项目结构

```bash
requests/
├── internal/           # 内部包（不对外暴露）
│   ├── client/         # HTTP 客户端核心实现
│   ├── models/         # 内部数据模型
│   └── utils/          # 内部工具
├── codec/              # 编解码器实现
├── url/                # URL 工具和文件上传
├── docs/               # 文档
│   ├── guides/         # 使用指南
│   └── examples/       # 代码示例
├── test/               # 性能测试
├── example/            # 示例程序
├── session.go          # Session 管理
├── methods.go          # HTTP 方法函数
├── types.go            # 类型导出
└── Makefile            # 构建自动化
```

## 贡献

欢迎贡献代码！请：

1. Fork 本仓库
2. 创建功能分支
3. 编写测试并提交更改
4. 提交 Pull Request

## 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE)
