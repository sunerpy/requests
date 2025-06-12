# Basic HTTP Requests

## GET 请求

```go
// 简单 GET
resp, err := requests.Get("https://api.example.com/users")

// 带查询参数
resp, err := requests.Get("https://api.example.com/users",
    requests.WithQuery("page", "1"),
    requests.WithQuery("limit", "10"),
)

// 带请求头
resp, err := requests.Get("https://api.example.com/users",
    requests.WithHeader("Authorization", "Bearer token"),
    requests.WithHeader("Accept", "application/json"),
)
```

## POST 请求

```go
// JSON 请求体
data := map[string]string{"name": "John", "email": "john@example.com"}
resp, err := requests.Post("https://api.example.com/users", data)

// 表单请求体
form := url.NewForm()
form.Set("username", "john")
form.Set("password", "secret")
resp, err := requests.Post("https://api.example.com/login", form)

// 字符串请求体
resp, err := requests.Post("https://api.example.com/data", "raw data")
```

## PUT 请求

```go
data := map[string]string{"name": "John Updated"}
resp, err := requests.Put("https://api.example.com/users/1", data)
```

## DELETE 请求

```go
resp, err := requests.Delete("https://api.example.com/users/1")
```

## PATCH 请求

```go
data := map[string]string{"status": "active"}
resp, err := requests.Patch("https://api.example.com/users/1", data)
```

## 处理响应

```go
resp, err := requests.Get("https://api.example.com/users")
if err != nil {
    log.Fatal(err)
}

// 状态码
fmt.Println("Status:", resp.StatusCode)

// 响应头
fmt.Println("Content-Type:", resp.Headers.Get("Content-Type"))

// 响应体
fmt.Println("Body:", resp.Text())      // 字符串
fmt.Println("Bytes:", resp.Bytes())    // 字节数组
```
