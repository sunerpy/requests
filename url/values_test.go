package url

import (
	"sync"
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

func TestNewForm(t *testing.T) {
	v := NewForm()
	assert.NotNil(t, v)
	assert.NotNil(t, v.data)
}

func TestNewURLParams(t *testing.T) {
	v := NewURLParams()
	assert.NotNil(t, v)
	assert.NotNil(t, v.data)
}

func TestParseParams(t *testing.T) {
	// Test with string
	v := ParseParams("key1=value1&key2=value2")
	assert.Equal(t, "value1", v.Get("key1"))
	assert.Equal(t, "value2", v.Get("key2"))
}

func TestValuesReset(t *testing.T) {
	v := NewValues()
	v.Add("key1", "value1")
	v.Add("key2", "value2")
	assert.Equal(t, "value1", v.Get("key1"))
	assert.Equal(t, "value2", v.Get("key2"))

	v.Reset()
	assert.Equal(t, "", v.Get("key1"))
	assert.Equal(t, "", v.Get("key2"))
	assert.Empty(t, v.Keys())
}

func TestAcquireReleaseValues(t *testing.T) {
	v := AcquireValues()
	assert.NotNil(t, v)

	v.Add("key1", "value1")
	assert.Equal(t, "value1", v.Get("key1"))

	ReleaseValues(v)

	// After release, we should be able to acquire again
	v2 := AcquireValues()
	assert.NotNil(t, v2)
	// The values should be reset
	assert.Equal(t, "", v2.Get("key1"))
	ReleaseValues(v2)
}

func TestReleaseValuesNil(t *testing.T) {
	// Should not panic
	ReleaseValues(nil)
}

// ============================================================================
// FastValues Tests
// ============================================================================

func TestNewFastValues(t *testing.T) {
	v := NewFastValues()
	assert.NotNil(t, v)
	assert.NotNil(t, v.data)
	assert.NotNil(t, v.keyIndex)
}

func TestFastValuesAdd(t *testing.T) {
	v := NewFastValues()
	v.Add("key1", "value1")
	v.Add("key1", "value2")
	v.Add("key2", "value3")

	assert.Equal(t, "value1", v.Get("key1"))
	assert.Equal(t, []string{"value1", "value2"}, v.GetAll("key1"))
	assert.Equal(t, "value3", v.Get("key2"))
}

func TestFastValuesSet(t *testing.T) {
	v := NewFastValues()
	v.Set("key1", "value1")
	v.Set("key1", "value2")

	assert.Equal(t, "value2", v.Get("key1"))
	assert.Equal(t, []string{"value2"}, v.GetAll("key1"))
}

func TestFastValuesGet(t *testing.T) {
	v := NewFastValues()
	v.Set("key1", "value1")

	assert.Equal(t, "value1", v.Get("key1"))
	assert.Equal(t, "", v.Get("nonexistent"))
}

func TestFastValuesGetAll(t *testing.T) {
	v := NewFastValues()
	v.Add("key1", "value1")
	v.Add("key1", "value2")

	assert.Equal(t, []string{"value1", "value2"}, v.GetAll("key1"))
	assert.Nil(t, v.GetAll("nonexistent"))
}

func TestFastValuesDel(t *testing.T) {
	v := NewFastValues()
	v.Add("key1", "value1")
	v.Add("key2", "value2")

	v.Del("key1")

	assert.Equal(t, "", v.Get("key1"))
	assert.False(t, v.Has("key1"))
	assert.Equal(t, "value2", v.Get("key2"))
}

func TestFastValuesHas(t *testing.T) {
	v := NewFastValues()
	v.Set("key1", "value1")

	assert.True(t, v.Has("key1"))
	assert.False(t, v.Has("nonexistent"))
}

func TestFastValuesLen(t *testing.T) {
	v := NewFastValues()
	assert.Equal(t, 0, v.Len())

	v.Add("key1", "value1")
	assert.Equal(t, 1, v.Len())

	v.Add("key1", "value2")
	assert.Equal(t, 1, v.Len()) // Same key, len should remain 1

	v.Add("key2", "value3")
	assert.Equal(t, 2, v.Len())
}

func TestFastValuesEncode(t *testing.T) {
	t.Run("encode with values", func(t *testing.T) {
		v := NewFastValues()
		v.Add("key1", "value1")
		v.Add("key2", "value2")

		encoded := v.Encode()
		assert.Contains(t, encoded, "key1=value1")
		assert.Contains(t, encoded, "key2=value2")
	})

	t.Run("encode empty values", func(t *testing.T) {
		v := NewFastValues()
		encoded := v.Encode()
		assert.Equal(t, "", encoded)
	})

	t.Run("encode with special characters", func(t *testing.T) {
		v := NewFastValues()
		v.Set("key", "value with spaces")

		encoded := v.Encode()
		assert.Contains(t, encoded, "key=value+with+spaces")
	})
}

func TestFastValuesKeys(t *testing.T) {
	v := NewFastValues()
	v.Add("key1", "value1")
	v.Add("key2", "value2")

	keys := v.Keys()
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
}

func TestFastValuesValues(t *testing.T) {
	v := NewFastValues()
	v.Add("key1", "value1")
	v.Add("key1", "value2")

	values := v.Values()
	assert.Equal(t, []string{"value1", "value2"}, values["key1"])
}

func TestFastValuesToURLValues(t *testing.T) {
	v := NewFastValues()
	v.Add("key1", "value1")
	v.Add("key1", "value2")
	v.Set("key2", "value3")

	urlValues := v.ToURLValues()
	assert.Equal(t, []string{"value1", "value2"}, urlValues["key1"])
	assert.Equal(t, []string{"value3"}, urlValues["key2"])
}

func TestFastValuesReset(t *testing.T) {
	v := NewFastValues()
	v.Add("key1", "value1")
	v.Add("key2", "value2")

	v.Reset()

	assert.Equal(t, 0, v.Len())
	assert.Equal(t, "", v.Get("key1"))
	assert.Equal(t, "", v.Get("key2"))
	assert.False(t, v.Has("key1"))
}

func TestAcquireReleaseFastValues(t *testing.T) {
	v := AcquireFastValues()
	assert.NotNil(t, v)

	v.Add("key1", "value1")
	assert.Equal(t, "value1", v.Get("key1"))

	ReleaseFastValues(v)

	// After release, we should be able to acquire again
	v2 := AcquireFastValues()
	assert.NotNil(t, v2)
	// The values should be reset
	assert.Equal(t, "", v2.Get("key1"))
	ReleaseFastValues(v2)
}

func TestReleaseFastValuesNil(t *testing.T) {
	// Should not panic
	ReleaseFastValues(nil)
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestValuesConcurrency(t *testing.T) {
	v := NewValues()
	var wg sync.WaitGroup
	iterations := 100

	// Test concurrent writes
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := "key" + string(rune('0'+idx%10))
			v.Add(key, "value")
			v.Set(key, "newvalue")
			_ = v.Get(key)
			_ = v.GetAll(key)
			_ = v.Keys()
			_ = v.Encode()
		}(i)
	}
	wg.Wait()
}

func TestValuesPoolConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	iterations := 100

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v := AcquireValues()
			v.Add("key", "value")
			_ = v.Get("key")
			ReleaseValues(v)
		}()
	}
	wg.Wait()
}
