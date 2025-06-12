# JSON API Examples

## 定义数据类型

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type APIResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

## GET JSON

```go
// 获取单个用户
result, err := requests.GetJSON[User]("https://api.example.com/users/1")
if err != nil {
    log.Fatal(err)
}

user := result.Data()
fmt.Printf("User: %s (%s)\n", user.Name, user.Email)

// 获取用户列表
result, err := requests.GetJSON[[]User]("https://api.example.com/users")
if err != nil {
    log.Fatal(err)
}

for _, user := range result.Data() {
    fmt.Printf("- %s\n", user.Name)
}
```

## POST JSON

```go
newUser := CreateUserRequest{
    Name:  "John Doe",
    Email: "john@example.com",
}

result, err := requests.PostJSON[User](
    "https://api.example.com/users",
    newUser,
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created user ID: %d\n", result.Data().ID)
```

## PUT JSON

```go
updateData := map[string]string{
    "name": "John Updated",
}

result, err := requests.PutJSON[User](
    "https://api.example.com/users/1",
    updateData,
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated: %s\n", result.Data().Name)
```

## DELETE JSON

```go
result, err := requests.DeleteJSON[APIResponse](
    "https://api.example.com/users/1",
)
if err != nil {
    log.Fatal(err)
}

if result.Data().Success {
    fmt.Println("User deleted successfully")
}
```

## 使用 Result[T]

```go
result, err := requests.GetJSON[User]("https://api.example.com/users/1")
if err != nil {
    log.Fatal(err)
}

// 访问数据
user := result.Data()

// 检查状态
if result.IsSuccess() {
    fmt.Println("Request successful")
}

// 获取状态码
fmt.Printf("Status: %d\n", result.StatusCode())

// 获取响应头
contentType := result.Headers().Get("Content-Type")
```

## 错误处理

```go
result, err := requests.GetJSON[User]("https://api.example.com/users/1")
if err != nil {
    // 检查错误类型
    if requests.IsTimeout(err) {
        log.Println("Request timed out")
    } else if requests.IsConnectionError(err) {
        log.Println("Connection failed")
    } else {
        log.Printf("Error: %v\n", err)
    }
    return
}

// 检查 HTTP 状态
if result.IsError() {
    log.Printf("HTTP Error: %d\n", result.StatusCode())
    return
}

// 处理成功响应
user := result.Data()
```
