# Go Requests

[![Language](https://img.shields.io/badge/language-golang-brightgreen)](https://github.com/sunerpy/requests) [![Last Commit](https://img.shields.io/github/last-commit/sunerpy/requests)](https://github.com/sunerpy/requests) [![CI](https://github.com/sunerpy/requests/workflows/Go/badge.svg)](https://github.com/sunerpy/requests/actions) [![codecov](https://codecov.io/gh/sunerpy/requests/branch/main/graph/badge.svg)](https://codecov.io/gh/sunerpy/requests) [![Benchmark](https://github.com/sunerpy/requests/actions/workflows/benchmark.yml/badge.svg)](https://sunerpy.github.io/requests/dev/bench) [![Go Reference](https://pkg.go.dev/badge/github.com/sunerpy/requests.svg)](https://pkg.go.dev/github.com/sunerpy/requests) [![Go Report Card](https://goreportcard.com/badge/github.com/sunerpy/requests)](https://goreportcard.com/report/github.com/sunerpy/requests) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](readme.md) | 简体中文

Go Requests 是一个简单且易于使用的 Go 语言 HTTP 请求库，灵感来源于 Python 的 Requests 库。它提供了易于理解和使用的 API，旨在简化 Go 中的 HTTP 请求。

## 特性

* 简单易用的 HTTP 请求 API
* 支持 GET、POST、PUT、DELETE 等常见 HTTP 请求方法
* 自动处理请求头、URL 编码、JSON 编码/解码等
* 支持 HTTP 请求超时设置
* 支持文件上传与下载
* 支持代理设置
* 支持自定义DNS服务器
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
    resp, err := requests.Get("https://api.github.com",nil)
    if err != nil {
        log.Fatal(err)
    }
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
    resp, err := requests.Post("https://httpbin.org/post", nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Status Code:", resp.StatusCode)
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
    resp, err := requests.Post("https://httpbin.org/post", nil)
    if err != nil {
        log.Fatal(err)
    }

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
    session := requests.NewSession()
    session.SetTimeout(5 * time.Second)
	req, err := requests.NewRequest("GET", "https://httpbin.org/get", nil, nil)
	if err != nil {
        log.Fatal(err)
    }
    resp, err := session.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", resp.Text())
}
```

## API

#### `requests.Get(baseURL string, params *url.Values) (*models.Response, error)`

发起一个 GET 请求。

* `url`：请求的 URL。

#### `requests.Post(baseURL string, params *url.Values) (*models.Response, error)`

发起一个 POST 请求。

* `url`：请求的 URL。

#### `requests.NewSession() Session`

创建一个新的Session，可以设置请求的默认配置，例如超时、代理等。

#### `(*Session) SetTimeout(timeout time.Duration)`

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
