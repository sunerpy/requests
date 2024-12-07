package url

import "net/http"

func NewRequest(method, rawURL string) (*http.Request, error) {
	return http.NewRequest(method, rawURL, nil)
}
