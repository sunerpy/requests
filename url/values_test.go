package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	t.Run("Add and Get", func(t *testing.T) {
		v := NewValues()
		v.Add("key1", "value1")
		v.Add("key1", "value2")
		v.Add("key2", "value3")
		assert.Equal(t, "value1", v.Get("key1"))
		assert.Equal(t, []string{"value1", "value2"}, v.GetAll("key1"))
		assert.Equal(t, "value3", v.Get("key2"))
	})
	t.Run("Set and Get", func(t *testing.T) {
		v := NewValues()
		v.Set("key1", "value1")
		v.Set("key1", "value2")
		assert.Equal(t, "value2", v.Get("key1"))
		assert.Equal(t, []string{"value2"}, v.GetAll("key1"))
	})
	t.Run("Del", func(t *testing.T) {
		v := NewValues()
		v.Add("key1", "value1")
		v.Add("key2", "value2")
		v.Del("key1")
		assert.Equal(t, "", v.Get("key1"))
		assert.Nil(t, v.GetAll("key1"))
		assert.NotContains(t, v.Keys(), "key1")
		assert.Equal(t, "value2", v.Get("key2"))
	})
	t.Run("Encode", func(t *testing.T) {
		v := NewValues()
		v.Add("key1", "value1")
		v.Add("key1", "value2")
		v.Add("key2", "value3")
		encoded := v.Encode()
		assert.Contains(t, encoded, "key1=value1")
		assert.Contains(t, encoded, "key1=value2")
		assert.Contains(t, encoded, "key2=value3")
	})
	t.Run("Keys", func(t *testing.T) {
		v := NewValues()
		v.Add("key1", "value1")
		v.Add("key2", "value2")
		keys := v.Keys()
		assert.Equal(t, []string{"key1", "key2"}, keys)
	})
	t.Run("Values", func(t *testing.T) {
		v := NewValues()
		v.Add("key1", "value1")
		v.Add("key1", "value2")
		expected := map[string][]string{
			"key1": {"value1", "value2"},
		}
		assert.Equal(t, expected, v.Values())
	})
	t.Run("Merge", func(t *testing.T) {
		v1 := NewValues()
		v1.Add("key1", "value1")
		v1.Add("key2", "value2")
		v2 := NewValues()
		v2.Add("key2", "value3")
		v2.Add("key3", "value4")
		v1.Merge(v2)
		assert.Equal(t, []string{"value1"}, v1.GetAll("key1"))
		assert.Equal(t, []string{"value2", "value3"}, v1.GetAll("key2"))
		assert.Equal(t, []string{"value4"}, v1.GetAll("key3"))
	})
	t.Run("Contains helper function", func(t *testing.T) {
		assert.True(t, contains([]string{"a", "b", "c"}, "b"))
		assert.False(t, contains([]string{"a", "b", "c"}, "d"))
	})
	t.Run("SearchStrings helper function", func(t *testing.T) {
		assert.Equal(t, 1, searchStrings([]string{"a", "b", "c"}, "b"))
		assert.Equal(t, -1, searchStrings([]string{"a", "b", "c"}, "d"))
	})
}
