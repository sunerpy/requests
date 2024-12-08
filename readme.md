# Go Requests

[![Language](https://img.shields.io/badge/language-golang-brightgreen)](https://github.com/sunerpy/requests) [![Last Commit](https://img.shields.io/github/last-commit/sunerpy/requests)](https://github.com/sunerpy/requests) [![CI](https://github.com/sunerpy/requests/workflows/Go/badge.svg)](https://github.com/sunerpy/requests/actions) [![codecov](https://codecov.io/gh/sunerpy/requests/branch/main/graph/badge.svg)](https://codecov.io/gh/sunerpy/requests) [![Benchmark](https://github.com/sunerpy/requests/actions/workflows/benchmark.yml/badge.svg)](https://sunerpy.github.io/requests/dev/bench) [![Go Reference](https://pkg.go.dev/badge/github.com/sunerpy/requests.svg)](https://pkg.go.dev/github.com/sunerpy/requests) [![Go Report Card](https://goreportcard.com/badge/github.com/sunerpy/requests)](https://goreportcard.com/report/github.com/sunerpy/requests) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

English | [简体中文](readme-cn.md)

Go Requests is a simple and easy-to-use HTTP request library for Go, inspired by Python's Requests library. It provides an easy-to-understand and use API designed to simplify HTTP requests in Go.

## Features

* Easy-to-use HTTP request API
* Supports common HTTP request methods such as GET, POST, PUT, DELETE, etc.
* Automatically handles headers, URL encoding, JSON encoding/decoding, etc.
* Supports HTTP request timeout configuration
* Supports file upload and download
* Supports proxy settings
* Supports DNS server settings
* Supports request redirection and cookie management

## Install

Install via Go modules：

```bash
go get github.com/sunerpy/requests
```

## Example

### Making a GET Request

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

### Making a POST Request

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

### File Upload

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

### Setting Request Timeout

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

Initiates a GET request。

* `url`：The URL of the request。

#### `requests.Post(baseURL string, params *url.Values) (*models.Response, error)`

Initiates a POST request。

* `url`：The URL of the request。

#### `requests.NewSession() Session`

Creates a new session, allowing you to set default configurations for requests, such as timeouts, proxies, etc.

#### `(*Session) SetTimeout(timeout time.Duration)`

Sets the timeout for the request.

#### `(*Response) JSON() interface{}`

Parses the response JSON data and returns it.

#### `(*Response) Text() string`

Returns the text content of the response.

## Contributions

Contributions are welcome! If you're interested in contributing, please follow these steps:

1. Fork this repository.
2. Make changes in your fork.
3. Submit a pull request.

## License

Go Requests is an open-source project licensed under the MIT license.
