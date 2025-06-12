# Session Usage Examples

## 创建 Session

```go
// 基本 Session
session := requests.NewSession()
defer session.Close()

// 带配置的 Session
session := requests.NewSession().
    WithBaseURL("https://api.example.com").
    WithTimeout(30 * time.Second).
    WithHeader("User-Agent", "MyApp/1.0")
defer session.Close()
```

## 使用 Base URL

```go
session := requests.NewSession().
    WithBaseURL("https://api.example.com/v1")
defer session.Close()

// 请求会自动拼接 base URL
req, _ := requests.NewGet("/users").Build()
resp, err := session.Do(req)  // 实际请求: https://api.example.com/v1/users
```

## 认证

```go
// Bearer Token
session := requests.NewSession().
    WithBearerToken("your-api-token")

// Basic Auth
session := requests.NewSession().
    WithBasicAuth("username", "password")
```

## 默认请求头

```go
session := requests.NewSession().
    WithHeader("Accept", "application/json").
    WithHeader("X-API-Key", "your-key").
    WithHeaders(map[string]string{
        "X-Custom-1": "value1",
        "X-Custom-2": "value2",
    })
```

## HTTP/2 支持

```go
// 启用 HTTP/2
session := requests.NewSession().
    WithHTTP2(true)

// 全局启用 HTTP/2
requests.SetHTTP2Enabled(true)
```

## 连接池配置

```go
session := requests.NewSession().
    WithMaxIdleConns(100).        // 最大空闲连接数
    WithIdleTimeout(90 * time.Second).  // 空闲超时
    WithKeepAlive(true)           // 启用 Keep-Alive
```

## 代理设置

```go
session := requests.NewSession().
    WithProxy("http://proxy.example.com:8080")
```

## 自定义 DNS

```go
session := requests.NewSession().
    WithDNS([]string{"8.8.8.8", "8.8.4.4"})
```

## Cookie 持久化

```go
session := requests.NewSession()
defer session.Close()

// 第一个请求 - 服务器设置 Cookie
req1, _ := requests.NewGet("https://api.example.com/login").Build()
session.Do(req1)

// 后续请求自动携带 Cookie
req2, _ := requests.NewGet("https://api.example.com/profile").Build()
session.Do(req2)
```

## 克隆 Session

```go
original := requests.NewSession().
    WithBaseURL("https://api.example.com").
    WithBearerToken("token")

// 克隆 Session（独立的配置）
clone := original.Clone().(requests.Session)

// 修改克隆不影响原始
clone.WithHeader("X-Custom", "value")
```

## 完整示例

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/sunerpy/requests"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    // 创建配置好的 Session
    session := requests.NewSession().
        WithBaseURL("https://api.example.com").
        WithTimeout(30 * time.Second).
        WithBearerToken("your-token").
        WithHTTP2(true).
        WithMaxIdleConns(100)
    defer session.Close()

    // 获取用户列表
    req, _ := requests.NewGet("/users").
        WithQuery("page", "1").
        Build()
    
    resp, err := session.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Status: %d\n", resp.StatusCode)
    fmt.Printf("Body: %s\n", resp.Text())
}
```
