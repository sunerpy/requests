package url

import (
	"net/url"
	"strconv"
	"strings"
)

func BuildURL(base string, params *Values) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}
	return u.String(), nil
}

// FastBuildURL builds a URL with query parameters using FastValues for better performance.
// This is optimized for single-threaded scenarios.
func FastBuildURL(base string, params *FastValues) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	if params != nil && params.Len() > 0 {
		u.RawQuery = params.Encode()
	}
	return u.String(), nil
}

// BuildURLFast builds a URL with query parameters from a map for maximum performance.
// Avoids intermediate allocations by building the query string directly.
func BuildURLFast(base string, params map[string]string) (string, error) {
	if len(params) == 0 {
		return base, nil
	}
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	// Build query string directly
	var sb strings.Builder
	sb.Grow(len(params) * 32) // Estimate ~32 bytes per param
	first := true
	for k, v := range params {
		if !first {
			sb.WriteByte('&')
		}
		sb.WriteString(url.QueryEscape(k))
		sb.WriteByte('=')
		sb.WriteString(url.QueryEscape(v))
		first = false
	}
	u.RawQuery = sb.String()
	return u.String(), nil
}

// NewURLParams creates a new thread-safe Values for URL parameters.
// This is the recommended default for most use cases as it provides
// safe concurrent access using sync.RWMutex.
//
// For high-performance single-threaded scenarios, use NewURLParamsUnsafe instead.
func NewURLParams() *Values {
	return NewValues()
}

// NewURLParamsUnsafe creates a new non-thread-safe FastValues for URL parameters.
// This provides better performance in single-threaded scenarios by avoiding
// lock overhead.
//
// WARNING: This is NOT safe for concurrent access. Use NewURLParams for
// thread-safe operations.
func NewURLParamsUnsafe() *FastValues {
	return NewFastValues()
}

// AcquireURLParams gets a thread-safe Values from the pool for URL parameters.
// Remember to call ReleaseURLParams when done.
// This combines object pooling with thread-safety for optimal performance
// in concurrent scenarios.
func AcquireURLParams() *Values {
	return AcquireValues()
}

// ReleaseURLParams returns a thread-safe Values to the pool.
func ReleaseURLParams(v *Values) {
	ReleaseValues(v)
}

// AcquireURLParamsUnsafe gets a non-thread-safe FastValues from the pool.
// Remember to call ReleaseURLParamsUnsafe when done.
//
// WARNING: This is NOT safe for concurrent access. Use AcquireURLParams for
// thread-safe operations.
func AcquireURLParamsUnsafe() *FastValues {
	return AcquireFastValues()
}

// ReleaseURLParamsUnsafe returns a non-thread-safe FastValues to the pool.
func ReleaseURLParamsUnsafe(v *FastValues) {
	ReleaseFastValues(v)
}

// BuildURLWithQuery builds a URL with query parameters from url.Values.
func BuildURLWithQuery(base string, query url.Values) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	if len(query) > 0 {
		// Merge with existing query parameters
		existingQuery := u.Query()
		for k, vs := range query {
			for _, v := range vs {
				existingQuery.Add(k, v)
			}
		}
		u.RawQuery = existingQuery.Encode()
	}
	return u.String(), nil
}

func ParseParams(data any) *Values {
	return parseValues(data)
}

func parseValues(data any) *Values {
	p := NewValues()
	switch v := data.(type) {
	case string:
		parseStringToValues(v, p)
	case map[string]string:
		parseMapString(v, p)
	case map[string][]string:
		parseMapStringSlice(v, p)
	case map[string]int:
		parseMapInt(v, p)
	case map[string][]int:
		parseMapIntSlice(v, p)
	case map[string]float64:
		parseMapFloat(v, p)
	case map[string][]float64:
		parseMapFloatSlice(v, p)
	case map[string]any:
		parseMapInterface(v, p)
	}
	return p
}

func parseStringToValues(data string, p *Values) {
	if data == "" {
		return
	}
	for _, l := range strings.Split(data, "&") {
		value := strings.SplitN(l, "=", 2)
		if len(value) == 2 {
			p.Add(value[0], value[1])
		}
	}
}

func parseMapString(data map[string]string, p *Values) {
	for key, value := range data {
		p.Set(key, value)
	}
}

func parseMapStringSlice(data map[string][]string, p *Values) {
	for key, values := range data {
		for _, value := range values {
			p.Add(key, value)
		}
	}
}

func parseMapInt(data map[string]int, p *Values) {
	for key, value := range data {
		p.Set(key, strconv.Itoa(value))
	}
}

func parseMapIntSlice(data map[string][]int, p *Values) {
	for key, values := range data {
		for _, value := range values {
			p.Add(key, strconv.Itoa(value))
		}
	}
}

func parseMapFloat(data map[string]float64, p *Values) {
	for key, value := range data {
		p.Set(key, strconv.FormatFloat(value, 'f', -1, 64))
	}
}

func parseMapFloatSlice(data map[string][]float64, p *Values) {
	for key, values := range data {
		for _, value := range values {
			p.Add(key, strconv.FormatFloat(value, 'f', -1, 64))
		}
	}
}

func parseMapInterface(data map[string]any, p *Values) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			p.Add(key, v)
		case []string:
			parseMapStringSlice(map[string][]string{key: v}, p)
		case int:
			p.Add(key, strconv.Itoa(v))
		case []int:
			parseMapIntSlice(map[string][]int{key: v}, p)
		case float64:
			p.Add(key, strconv.FormatFloat(v, 'f', -1, 64))
		case []float64:
			parseMapFloatSlice(map[string][]float64{key: v}, p)
		case bool:
			p.Add(key, strconv.FormatBool(v))
		case []any:
			for _, item := range v {
				parseSingleInterface(key, item, p)
			}
		}
	}
}

func parseSingleInterface(key string, value any, p *Values) {
	switch v := value.(type) {
	case string:
		p.Add(key, v)
	case int:
		p.Add(key, strconv.Itoa(v))
	case float64:
		p.Add(key, strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		p.Add(key, strconv.FormatBool(v))
	}
}
