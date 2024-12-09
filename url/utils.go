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

func ParseParams(data interface{}) *Values {
	return parseValues(data)
}

func parseValues(data interface{}) *Values {
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
	case map[string]interface{}:
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

func parseMapInterface(data map[string]interface{}, p *Values) {
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
		case []interface{}:
			for _, item := range v {
				parseSingleInterface(key, item, p)
			}
		}
	}
}

func parseSingleInterface(key string, value interface{}, p *Values) {
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
