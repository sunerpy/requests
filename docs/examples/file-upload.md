# File Upload Examples

## 基本文件上传

```go
import "github.com/sunerpy/requests/url"

// 创建文件上传
upload, err := url.NewFileUpload("file", "/path/to/document.pdf")
if err != nil {
    log.Fatal(err)
}

// 发送请求
resp, err := requests.Post("https://api.example.com/upload", upload)
```

## 带进度回调的上传

```go
upload, err := url.NewFileUpload("file", "/path/to/large-file.zip")
if err != nil {
    log.Fatal(err)
}

// 设置进度回调
upload.WithProgress(func(uploaded, total int64) {
    percent := float64(uploaded) / float64(total) * 100
    fmt.Printf("\rUploading: %.1f%% (%d/%d bytes)", percent, uploaded, total)
})

resp, err := requests.Post("https://api.example.com/upload", upload)
fmt.Println("\nUpload complete!")
```

## 多文件上传

```go
// 创建多部分表单
form := url.NewMultipartForm()

// 添加文件
file1, _ := url.NewFileUpload("files", "/path/to/file1.jpg")
file2, _ := url.NewFileUpload("files", "/path/to/file2.jpg")
form.AddFile(file1)
form.AddFile(file2)

// 添加其他字段
form.AddField("description", "My photos")
form.AddField("album", "vacation")

resp, err := requests.Post("https://api.example.com/upload", form)
```

## 从内存上传

```go
// 从字节数组上传
data := []byte("file content here")
upload := url.NewBytesUpload("file", "document.txt", data)

resp, err := requests.Post("https://api.example.com/upload", upload)
```

## 带认证的上传

```go
upload, _ := url.NewFileUpload("file", "/path/to/file.pdf")

resp, err := requests.Post(
    "https://api.example.com/upload",
    upload,
    requests.WithBearerToken("your-token"),
    requests.WithHeader("X-Custom-Header", "value"),
)
```

## 使用 Session 上传

```go
session := requests.NewSession().
    WithBaseURL("https://api.example.com").
    WithBearerToken("token").
    WithTimeout(5 * time.Minute)  // 大文件需要更长超时
defer session.Close()

upload, _ := url.NewFileUpload("file", "/path/to/large-file.zip")

req, _ := requests.NewPost("/upload").
    WithBody(upload).
    Build()

resp, err := session.Do(req)
```
