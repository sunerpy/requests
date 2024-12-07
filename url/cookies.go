package url

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// Cookies 管理 HTTP Cookies
type Cookies struct {
	cookies []*http.Cookie
}

// NewCookies 创建一个 Cookies 实例
func NewCookies() *Cookies {
	return &Cookies{}
}

// Add 添加一个 Cookie
func (c *Cookies) Add(cookie *http.Cookie) {
	c.cookies = append(c.cookies, cookie)
}

// GetCookies 获取指定 URL 的 Cookies
func GetCookies(rawURL string, cookies []*http.Cookie) []*http.Cookie {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	// 初始化 DefaultClient 的 Jar
	if http.DefaultClient.Jar == nil {
		http.DefaultClient.Jar = newMemoryCookieJar()
	}
	http.DefaultClient.Jar.SetCookies(u, cookies)
	return http.DefaultClient.Jar.Cookies(u)
}

// newMemoryCookieJar 创建一个新的 CookieJar 实例
func newMemoryCookieJar() http.CookieJar {
	jar, _ := cookiejar.New(nil)
	return jar
}
