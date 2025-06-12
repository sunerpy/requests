# Getting Started

本指南将帮助您快速开始使用 Go Requests 库。

## 安装

```bash
go get github.com/sunerpy/requests
```

## 快速开始

### 基本 GET 请求

```go
package main

import (
    "fmt"
    "log"

    "github.com/sunerpy/requests"
)

func main() {
    resp, err := requests.Get("https://api.github.com")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Status:", resp.StatusCode)
    fmt.Println("Body:", resp.Text())
}
```

### 带参数的 GET 请求

```go
resp, err := requests.Get("https://api.github.com/search/repos",
    requests.WithQuery("q", "golang"),
    requests.WithQuery("sort", "stars"),
)
```

### JSON 响应解析

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

result, err := requests.GetJSON[User]("https://api.example.com/user/1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User: %+v\n", result.Data())
fmt.Printf("Status: %d\n", result.StatusCode())
```

### POST 请求

```go
type CreateUser struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponse struct {
    ID int `json:"id"`
}

result, err := requests.PostJSON[UserResponse](
    "https://api.example.com/users",
    CreateUser{Name: "John", Email: "john@example.com"},
)
```

## 基本概念

### Result[T] 类型

泛型方法（如 `GetJSON[T]`）返回 `Result[T]` 类型，它包装了解析后的数据和响应元数据：

```go
result, err := requests.GetJSON[User](url)

// 访问解析后的数据
user := result.Data()

// 访问响应元数据
statusCode := result.StatusCode()
headers := result.Headers()
isSuccess := result.IsSuccess()
```

### RequestBuilder

使用 Builder 模式构建复杂请求：

```go
req, err := requests.NewPost("https://api.example.com/users").
    WithHeader("X-Custom-Header", "value").
    WithQuery("version", "v2").
    WithJSON(map[string]string{"name": "John"}).
    WithTimeout(10 * time.Second).
    Build()
```

### Session

Session 提供连接池、Cookie 持久化和默认配置：

```go
session := requests.NewSession().
    WithBaseURL("https://api.example.com").
    WithTimeout(30 * time.Second).
    WithHeader("Authorization", "Bearer token")

defer session.Close()

req, _ := requests.NewGet("/users").Build()
resp, err := session.Do(req)
```

## 下一步

- 查看 [API 参考](api-reference.md) 了解完整 API
- 查看 [示例](../examples/) 了解更多用例
- 查看 [迁移指南](migration-guide.md) 从旧版本升级
