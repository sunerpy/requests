package url

import (
	"net/url"
	"strings"
	"sync"
)

type Values struct {
	data     map[string][]string
	mu       sync.RWMutex
	indexKey []string
}

func NewValues() *Values {
	return &Values{
		data: make(map[string][]string),
	}
}

func NewForm() *Values {
	return &Values{
		data: make(map[string][]string),
	}
}

func NewURLParams() *Values {
	return &Values{
		data: make(map[string][]string),
	}
}

func (v *Values) Add(key, value string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.data[key] = append(v.data[key], value)
	if !contains(v.indexKey, key) {
		v.indexKey = append(v.indexKey, key)
	}
}

func (v *Values) Set(key, value string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.data[key] = []string{value}
	if !contains(v.indexKey, key) {
		v.indexKey = append(v.indexKey, key)
	}
}

func (v *Values) Get(key string) string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if val, ok := v.data[key]; ok && len(val) > 0 {
		return val[0]
	}
	return ""
}

func (v *Values) GetAll(key string) []string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if val, ok := v.data[key]; ok {
		return val
	}
	return nil
}

func (v *Values) Del(key string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	delete(v.data, key)
	index := searchStrings(v.indexKey, key)
	if index != -1 {
		v.indexKey = append(v.indexKey[:index], v.indexKey[index+1:]...)
	}
}

func (v *Values) Encode() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	var sb strings.Builder
	first := true
	for key, values := range v.data {
		for _, value := range values {
			if !first {
				sb.WriteByte('&')
			}
			sb.WriteString(url.QueryEscape(key))
			sb.WriteByte('=')
			sb.WriteString(url.QueryEscape(value))
			first = false
		}
	}
	return sb.String()
}

func (v *Values) Keys() []string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.indexKey
}

func (v *Values) Values() map[string][]string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	result := make(map[string][]string)
	for key, values := range v.data {
		result[key] = append([]string(nil), values...)
	}
	return result
}

func (v *Values) Merge(other *Values) {
	other.mu.RLock()
	defer other.mu.RUnlock()
	v.mu.Lock()
	defer v.mu.Unlock()
	for key, values := range other.data {
		v.data[key] = append(v.data[key], values...)
		if !contains(v.indexKey, key) {
			v.indexKey = append(v.indexKey, key)
		}
	}
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func searchStrings(slice []string, str string) int {
	for i, v := range slice {
		if v == str {
			return i
		}
	}
	return -1
}
