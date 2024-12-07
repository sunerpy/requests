package url

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCookies(t *testing.T) {
	cookies := NewCookies()
	assert.NotNil(t, cookies)
	assert.Empty(t, cookies.cookies)
}

func TestCookies_Add(t *testing.T) {
	cookies := NewCookies()
	cookie := &http.Cookie{
		Name:  "session",
		Value: "12345",
	}
	cookies.Add(cookie)
	assert.Len(t, cookies.cookies, 1)
	assert.Equal(t, "session", cookies.cookies[0].Name)
	assert.Equal(t, "12345", cookies.cookies[0].Value)
}

func TestGetCookies(t *testing.T) {
	t.Run("Valid URL with Cookies", func(t *testing.T) {
		rawURL := "https://example.com"
		cookies := []*http.Cookie{
			{Name: "token", Value: "abcdef"},
			{Name: "user", Value: "john_doe"},
		}
		gotCookies := GetCookies(rawURL, cookies)
		assert.Len(t, gotCookies, 2)
		assert.Equal(t, "token", gotCookies[0].Name)
		assert.Equal(t, "abcdef", gotCookies[0].Value)
		assert.Equal(t, "user", gotCookies[1].Name)
		assert.Equal(t, "john_doe", gotCookies[1].Value)
	})
	t.Run("Invalid URL", func(t *testing.T) {
		rawURL := "://invalid-url"
		cookies := []*http.Cookie{
			{Name: "token", Value: "abcdef"},
		}
		gotCookies := GetCookies(rawURL, cookies)
		assert.Nil(t, gotCookies)
	})
}
