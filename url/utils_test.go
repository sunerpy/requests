package url

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildURL(t *testing.T) {
	t.Run("Valid URL with params", func(t *testing.T) {
		params := NewValues()
		params.Set("key1", "value1")
		params.Set("key2", "value2")
		uri, err := BuildURL("https://example.com", params)
		assert.NoError(t, err)
		parsedURL, err := url.Parse(uri)
		assert.NoError(t, err)
		actualParams, err := url.ParseQuery(parsedURL.RawQuery)
		assert.NoError(t, err)
		expectedParams := url.Values{
			"key1": []string{"value1"},
			"key2": []string{"value2"},
		}
		assert.Equal(t, expectedParams, actualParams)
	})
	t.Run("Invalid URL", func(t *testing.T) {
		_, err := BuildURL("://invalid-url", nil)
		assert.Error(t, err)
	})
	t.Run("Valid URL without params", func(t *testing.T) {
		url, err := BuildURL("https://example.com", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", url)
	})
}

func TestParseValues(t *testing.T) {
	t.Run("Parse string", func(t *testing.T) {
		data := "key1=value1&key2=value2"
		values := parseValues(data)
		assert.Equal(t, "value1", values.Get("key1"))
		assert.Equal(t, "value2", values.Get("key2"))
	})
	t.Run("Parse empty string", func(t *testing.T) {
		data := ""
		values := parseValues(data)
		assert.Empty(t, values.Keys())
	})
	t.Run("Parse map[string]string", func(t *testing.T) {
		data := map[string]string{"key1": "value1", "key2": "value2"}
		values := parseValues(data)
		assert.Equal(t, "value1", values.Get("key1"))
		assert.Equal(t, "value2", values.Get("key2"))
	})
	t.Run("Parse map[string][]string", func(t *testing.T) {
		data := map[string][]string{"key1": {"value1", "value2"}}
		values := parseValues(data)
		assert.Equal(t, []string{"value1", "value2"}, values.GetAll("key1"))
	})
	t.Run("Parse map[string]int", func(t *testing.T) {
		data := map[string]int{"key1": 1, "key2": 2}
		values := parseValues(data)
		assert.Equal(t, "1", values.Get("key1"))
		assert.Equal(t, "2", values.Get("key2"))
	})
	t.Run("Parse map[string][]int", func(t *testing.T) {
		data := map[string][]int{"key1": {1, 2}}
		values := parseValues(data)
		assert.Equal(t, []string{"1", "2"}, values.GetAll("key1"))
	})
	t.Run("Parse map[string]float64", func(t *testing.T) {
		data := map[string]float64{"key1": 1.23, "key2": 4.56}
		values := parseValues(data)
		assert.Equal(t, "1.23", values.Get("key1"))
		assert.Equal(t, "4.56", values.Get("key2"))
	})
	t.Run("Parse map[string][]float64", func(t *testing.T) {
		data := map[string][]float64{"key1": {1.23, 4.56}}
		values := parseValues(data)
		assert.Equal(t, []string{"1.23", "4.56"}, values.GetAll("key1"))
	})
	t.Run("Parse map[string]any", func(t *testing.T) {
		data := map[string]any{
			"key1": "value1",
			"key2": []string{"value2", "value3"},
			"key3": 123,
			"key4": []int{4, 5},
			"key5": 1.23,
			"key6": []float64{4.56, 7.89},
			"key7": true,
			"key8": []any{"value4", 6, 7.8, false},
		}
		values := parseValues(data)
		assert.Equal(t, "value1", values.Get("key1"))
		assert.Equal(t, []string{"value2", "value3"}, values.GetAll("key2"))
		assert.Equal(t, "123", values.Get("key3"))
		assert.Equal(t, []string{"4", "5"}, values.GetAll("key4"))
		assert.Equal(t, "1.23", values.Get("key5"))
		assert.Equal(t, []string{"4.56", "7.89"}, values.GetAll("key6"))
		assert.Equal(t, "true", values.Get("key7"))
		assert.Equal(t, []string{"value4", "6", "7.8", "false"}, values.GetAll("key8"))
	})
}

func TestBuildURLWithQuery(t *testing.T) {
	t.Run("Valid URL with query params", func(t *testing.T) {
		query := url.Values{}
		query.Set("key1", "value1")
		query.Set("key2", "value2")
		uri, err := BuildURLWithQuery("https://example.com", query)
		assert.NoError(t, err)
		parsedURL, err := url.Parse(uri)
		assert.NoError(t, err)
		actualParams := parsedURL.Query()
		assert.Equal(t, "value1", actualParams.Get("key1"))
		assert.Equal(t, "value2", actualParams.Get("key2"))
	})
	t.Run("Valid URL without query params", func(t *testing.T) {
		uri, err := BuildURLWithQuery("https://example.com", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", uri)
	})
	t.Run("Valid URL with empty query params", func(t *testing.T) {
		query := url.Values{}
		uri, err := BuildURLWithQuery("https://example.com", query)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", uri)
	})
	t.Run("Merge with existing query params", func(t *testing.T) {
		query := url.Values{}
		query.Set("new_key", "new_value")
		uri, err := BuildURLWithQuery("https://example.com?existing=param", query)
		assert.NoError(t, err)
		parsedURL, err := url.Parse(uri)
		assert.NoError(t, err)
		actualParams := parsedURL.Query()
		assert.Equal(t, "param", actualParams.Get("existing"))
		assert.Equal(t, "new_value", actualParams.Get("new_key"))
	})
	t.Run("Invalid URL", func(t *testing.T) {
		query := url.Values{}
		query.Set("key", "value")
		_, err := BuildURLWithQuery("://invalid-url", query)
		assert.Error(t, err)
	})
	t.Run("Multiple values for same key", func(t *testing.T) {
		query := url.Values{}
		query.Add("key", "value1")
		query.Add("key", "value2")
		uri, err := BuildURLWithQuery("https://example.com", query)
		assert.NoError(t, err)
		parsedURL, err := url.Parse(uri)
		assert.NoError(t, err)
		actualParams := parsedURL.Query()
		assert.Equal(t, []string{"value1", "value2"}, actualParams["key"])
	})
}
