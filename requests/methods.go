package requests

import (
	"io"
	"strings"

	"github.com/sunerpy/requests/models"
	"github.com/sunerpy/requests/url"
)

const (
	contentKey      = "Content-Type"
	formContentType = "application/x-www-form-urlencoded"
)

// Get 发送 GET 请求
func Get(baseURL string, params *url.Values) (*models.Response, error) {
	u, err := url.BuildURL(baseURL, params)
	if err != nil {
		return nil, err
	}
	req, err := NewRequest("GET", u, params, nil)
	if err != nil {
		return nil, err
	}
	return defaultSess.Do(req)
}

func Post(baseURL string, form *url.Values) (*models.Response, error) {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := NewRequest("POST", baseURL, nil, body)
	if err != nil {
		return nil, err
	}
	req.AddHeader(contentKey, formContentType)
	return defaultSess.Do(req)
}

// Put 发送 PUT 请求
func Put(baseURL string, form *url.Values) (*models.Response, error) {
	body := ""
	if form != nil {
		body = form.Encode()
	}
	req, err := NewRequest("PUT", baseURL, nil, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.AddHeader(contentKey, formContentType)
	return defaultSess.Do(req)
}

// Delete 发送 DELETE 请求
func Delete(baseURL string, params *url.Values) (*models.Response, error) {
	u, err := url.BuildURL(baseURL, params)
	if err != nil {
		return nil, err
	}
	req, err := NewRequest("DELETE", u, params, nil)
	if err != nil {
		return nil, err
	}
	return defaultSess.Do(req)
}

// Patch 发送 PATCH 请求
func Patch(baseURL string, form *url.Values) (*models.Response, error) {
	body := ""
	if form != nil {
		body = form.Encode()
	}
	req, err := NewRequest("PATCH", baseURL, nil, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.AddHeader(contentKey, formContentType)
	return defaultSess.Do(req)
}
