package requests

import (
	"io"

	"github.com/sunerpy/requests/models"
	"github.com/sunerpy/requests/requests"
	"github.com/sunerpy/requests/url"
)

func Get(baseURL string, params *url.Values) (*models.Response, error) {
	return requests.Get(baseURL, params)
}

// Post 发送 POST 请求
func Post(baseURL string, form *url.Values) (*models.Response, error) {
	return requests.Post(baseURL, form)
}

// Put 发送 PUT 请求
func Put(baseURL string, form *url.Values) (*models.Response, error) {
	return requests.Put(baseURL, form)
}

// Delete 发送 DELETE 请求
func Delete(baseURL string, params *url.Values) (*models.Response, error) {
	return requests.Delete(baseURL, params)
}

// Patch 发送 PATCH 请求
func Patch(baseURL string, form *url.Values) (*models.Response, error) {
	return requests.Patch(baseURL, form)
}

func SetHTTP2Enabled(enabled bool) {
	requests.SetHTTP2Enabled(enabled)
}

func NewSession() requests.Session {
	return requests.NewSession()
}

func NewRequest(method, rawURL string, params *url.Values, body io.Reader) (*requests.Request, error) {
	return requests.NewRequest(method, rawURL, params, body)
}
