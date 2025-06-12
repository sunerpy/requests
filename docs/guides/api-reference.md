# API Reference

## HTTP 方法

### 基础方法

| 函数 | 描述 |
| ------ | ------ |
| `Get(url, opts...)` | 发送 GET 请求 |
| `Post(url, body, opts...)` | 发送 POST 请求 |
| `Put(url, body, opts...)` | 发送 PUT 请求 |
| `Delete(url, opts...)` | 发送 DELETE 请求 |
| `Patch(url, body, opts...)` | 发送 PATCH 请求 |
| `Head(url, opts...)` | 发送 HEAD 请求 |
| `Options(url, opts...)` | 发送 OPTIONS 请求 |

### 泛型方法

| 函数 | 描述 |
| ------ | ------ |
| `GetJSON[T](url, opts...)` | GET 并解析 JSON，返回 `(Result[T], error)` |
| `PostJSON[T](url, data, opts...)` | POST JSON，返回 `(Result[T], error)` |
| `PutJSON[T](url, data, opts...)` | PUT JSON，返回 `(Result[T], error)` |
| `DeleteJSON[T](url, opts...)` | DELETE 并解析 JSON，返回 `(Result[T], error)` |
| `PatchJSON[T](url, data, opts...)` | PATCH JSON，返回 `(Result[T], error)` |

### 便捷方法

| 函数 | 描述 |
| ------ | ------ |
| `GetString(url, opts...)` | GET 并返回字符串 |
| `GetBytes(url, opts...)` | GET 并返回字节数组 |
| `DoJSON[T](builder)` | 执行 Builder 并解析 JSON |
| `DoXML[T](builder)` | 执行 Builder 并解析 XML |

## RequestBuilder

### 构造函数

| 函数 | 描述 |
| ------ | ------ |
| `NewGet(url)` | 创建 GET 请求 Builder |
| `NewPost(url)` | 创建 POST 请求 Builder |
| `NewPut(url)` | 创建 PUT 请求 Builder |
| `NewDeleteBuilder(url)` | 创建 DELETE 请求 Builder |
| `NewPatch(url)` | 创建 PATCH 请求 Builder |
| `NewRequestBuilder(method, url)` | 创建任意方法的 Builder |

### Builder 方法

```go
builder.WithHeader(key, value)      // 添加请求头
builder.WithHeaders(map)            // 添加多个请求头
builder.WithQuery(key, value)       // 添加查询参数
builder.WithQueryParams(map)        // 添加多个查询参数
builder.WithBody(reader)            // 设置请求体
builder.WithJSON(data)              // 设置 JSON 请求体
builder.WithForm(values)            // 设置表单请求体
builder.WithBasicAuth(user, pass)   // 设置 Basic 认证
builder.WithBearerToken(token)      // 设置 Bearer Token
builder.WithTimeout(duration)       // 设置超时
builder.WithContext(ctx)            // 设置 Context
builder.Build()                     // 构建请求
builder.Do()                        // 构建并执行请求
```

## Result[T]

### 方法

| 方法 | 描述 |
| ------ | ------ |
| `Data()` | 获取解析后的数据（类型 T） |
| `Response()` | 获取底层 Response 对象 |
| `StatusCode()` | 获取 HTTP 状态码 |
| `Headers()` | 获取响应头 |
| `IsSuccess()` | 检查状态码是否为 2xx |
| `IsError()` | 检查状态码是否为 4xx 或 5xx |
| `Cookies()` | 获取响应 Cookie |

## Request Options

| 选项 | 描述 |
| ------ | ------ |
| `WithHeader(key, value)` | 添加请求头 |
| `WithHeaders(map)` | 添加多个请求头 |
| `WithQuery(key, value)` | 添加查询参数 |
| `WithQueryParams(map)` | 添加多个查询参数 |
| `WithBasicAuth(user, pass)` | 设置 Basic 认证 |
| `WithBearerToken(token)` | 设置 Bearer Token |
| `WithTimeout(duration)` | 设置超时 |
| `WithContentType(type)` | 设置 Content-Type |
| `WithContext(ctx)` | 设置 Context |
| `WithAccept(type)` | 设置 Accept 头 |

## Session

### 创建和配置

```go
session := requests.NewSession()

session.WithBaseURL(url)           // 设置基础 URL
session.WithTimeout(duration)      // 设置超时
session.WithHeader(key, value)     // 添加默认请求头
session.WithHeaders(map)           // 添加多个默认请求头
session.WithBasicAuth(user, pass)  // 设置 Basic 认证
session.WithBearerToken(token)     // 设置 Bearer Token
session.WithHTTP2(enabled)         // 启用/禁用 HTTP/2
session.WithKeepAlive(enabled)     // 启用/禁用 Keep-Alive
session.WithMaxIdleConns(n)        // 设置最大空闲连接数
session.WithProxy(url)             // 设置代理
session.WithDNS(servers)           // 设置自定义 DNS
session.Clone()                    // 克隆 Session
session.Close()                    // 关闭 Session
```

## Middleware

### 创建中间件链

```go
chain := requests.NewMiddlewareChain()

chain.Use(middleware)              // 添加中间件
chain.UseFunc(func)                // 添加函数中间件
chain.Execute(req, handler)        // 执行请求
```

### 内置中间件

```go
requests.HeaderMiddleware(headers) // 添加请求头的中间件
requests.HooksMiddleware(hooks)    // Hooks 中间件
```

## Hooks

```go
hooks := requests.NewHooks()

hooks.OnRequest(func(req))                    // 请求前回调
hooks.OnResponse(func(req, resp, duration))   // 响应后回调
hooks.OnError(func(req, err, duration))       // 错误回调
```

## Retry

### 重试策略

```go
requests.NoRetryPolicy()                              // 不重试
requests.LinearRetryPolicy(maxRetries, delay)         // 线性退避
requests.ExponentialRetryPolicy(maxRetries, min, max) // 指数退避
```

### 重试执行器

```go
executor := requests.NewRetryExecutor(policy)
resp, err := executor.Execute(ctx, func() (*Response, error) {
    return requests.Get(url)
})
```

## 性能优化 API

### FastValues

无锁的 URL 参数构建器，适用于单线程高性能场景：

```go
import "github.com/sunerpy/requests/url"

// 创建 FastValues
fv := url.NewFastValues()
fv.Add("key", "value")
fv.Add("foo", "bar")
encoded := fv.Encode()  // "foo=bar&key=value"

// 使用对象池（推荐）
fv := url.AcquireFastValues()
defer url.ReleaseFastValues(fv)
fv.Add("page", "1")
```

**注意**: `FastValues` 不是线程安全的，仅适用于单线程场景。并发场景请使用 `url.NewValues()`。

### 内部优化

本库内部已实现以下优化，用户无需额外配置：

- **RequestBuilder 对象池**: 减少构建器分配
- **RequestConfig 对象池**: 减少配置对象分配  
- **Header 所有权转移**: 避免 Header map 克隆
- **响应缓冲池**: 分级缓冲池 (4KB/32KB) 复用

详见 [性能优化指南](performance-optimization.md)。

## 错误处理

### 错误类型

| 类型 | 描述 |
| ------ | ------ |
| `RequestError` | 请求创建或发送错误 |
| `ResponseError` | 响应处理错误 |
| `TimeoutError` | 超时错误 |
| `ConnectionError` | 连接错误 |
| `DecodeError` | 解码错误 |
| `EncodeError` | 编码错误 |
| `RetryError` | 重试耗尽错误 |

### 错误检查

```go
requests.IsTimeout(err)         // 检查是否超时
requests.IsConnectionError(err) // 检查是否连接错误
requests.IsResponseError(err)   // 检查是否响应错误
requests.IsTemporary(err)       // 检查是否临时错误
```
