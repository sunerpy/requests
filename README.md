# Go Requests

![Codecov](https://img.shields.io/codecov/c/github/sunerpy/requests)

Go Requests 是一个简单且易于使用的 Go 语言 HTTP 请求库，灵感来源于 Python 的 Requests 库。它提供了易于理解和使用的 API，旨在简化 Go 中的 HTTP 请求。

## 特性

* 简单易用的 HTTP 请求 API
* 支持 GET、POST、PUT、DELETE 等常见 HTTP 请求方法
* 自动处理请求头、URL 编码、JSON 编码/解码等
* 支持 HTTP 请求超时设置
* 支持文件上传与下载
* 支持代理设置
* 支持请求重定向与 cookie 管理

## 安装

通过 Go modules 安装：

```bash
go get github.com/sunerpy/requests
```

## 示例

### 发起 GET 请求

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
    defer resp.Body.Close()

    fmt.Println("Status Code:", resp.StatusCode)
    fmt.Println("Response Body:", resp.Text())
}
```

### 发起 POST 请求

```go
package main

import (
    "fmt"
    "log"

    "github.com/sunerpy/requests"
)

func main() {
    resp, err := requests.Post("https://httpbin.org/post", requests.Params{"key": "value"})
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Status Code:", resp.StatusCode)
    fmt.Println("Response JSON:", resp.JSON())
}
```

### 发起带 JSON 数据的 POST 请求

```go
package main

import (
    "fmt"
    "log"

    "github.com/sunerpy/requests"
)

func main() {
    resp, err := requests.Post("https://httpbin.org/post", requests.JSON{"name": "John", "age": 30})
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Response JSON:", resp.JSON())
}
```

### 文件上传

```go
package main

import (
    "fmt"
    "log"

    "github.com/sunerpy/requests"
)

func main() {
    resp, err := requests.Post("https://httpbin.org/post", requests.File("file", "/path/to/file.txt"))
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Response JSON:", resp.JSON())
}
```

### 设置请求超时

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/sunerpy/requests"
)

func main() {
    client := requests.NewClient()
    client.SetTimeout(5 * time.Second)

    resp, err := client.Get("https://httpbin.org/get")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Response:", resp.Text())
}
```

## API

#### `requests.Get(url string, options ...Option) (*Response, error)`

发起一个 GET 请求。

* `url`：请求的 URL。
* `options`：可选参数，例如请求头、超时设置等。

#### `requests.Post(url string, options ...Option) (*Response, error)`

发起一个 POST 请求。

* `url`：请求的 URL。
* `options`：可选参数，例如表单数据、JSON 数据、文件等。

#### `requests.NewClient() *Client`

创建一个新的客户端，可以设置请求的默认配置，例如超时、代理等。

#### `(*Client) SetTimeout(timeout time.Duration)`

设置请求的超时时间。

#### `(*Response) JSON() interface{}`

解析响应的 JSON 数据并返回。

#### `(*Response) Text() string`

返回响应的文本内容。

## 贡献

欢迎贡献代码！如果你有兴趣贡献，请按照以下步骤：

1. Fork 本仓库。
2. 在你的 Fork 上进行修改。
3. 提交 Pull Request。

## 许可证

Go Requests 是一个开源项目，采用 MIT 许可证。
