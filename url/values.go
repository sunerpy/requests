package url

import (
	"net/url"
	"strings"
	"sync"
)

// ============================================================================
// Thread-safe Values (original implementation)
// ============================================================================
type Values struct {
	data     map[string][]string
	mu       sync.RWMutex
	indexKey []string
}

// valuesPool is a pool of thread-safe Values objects for reuse
var valuesPool = sync.Pool{
	New: func() any {
		return &Values{
			data:     make(map[string][]string, 8),
			indexKey: make([]string, 0, 8),
		}
	},
}

func NewValues() *Values {
	return &Values{
		data:     make(map[string][]string, 8),
		indexKey: make([]string, 0, 8),
	}
}

func NewForm() *Values {
	return &Values{
		data:     make(map[string][]string, 8),
		indexKey: make([]string, 0, 8),
	}
}

// AcquireValues gets a thread-safe Values from the pool.
// Remember to call ReleaseValues when done.
func AcquireValues() *Values {
	return valuesPool.Get().(*Values)
}

// ReleaseValues returns a thread-safe Values to the pool.
func ReleaseValues(v *Values) {
	if v == nil {
		return
	}
	v.Reset()
	valuesPool.Put(v)
}

// Reset clears the Values for reuse.
func (v *Values) Reset() {
	v.mu.Lock()
	defer v.mu.Unlock()
	for k := range v.data {
		delete(v.data, k)
	}
	v.indexKey = v.indexKey[:0]
}

// ============================================================================
// FastValues - Lock-free version for single-threaded scenarios
// ============================================================================
// FastValues is a high-performance, non-thread-safe URL values container.
// Use this when you don't need concurrent access for better performance.
type FastValues struct {
	data     map[string][]string
	keyIndex map[string]struct{} // O(1) key existence check
}

// fastValuesPool is a pool of FastValues objects for reuse
var fastValuesPool = sync.Pool{
	New: func() any {
		return &FastValues{
			data:     make(map[string][]string, 8),
			keyIndex: make(map[string]struct{}, 8),
		}
	},
}

// NewFastValues creates a new FastValues instance.
func NewFastValues() *FastValues {
	return &FastValues{
		data:     make(map[string][]string, 8),
		keyIndex: make(map[string]struct{}, 8),
	}
}

// AcquireFastValues gets a FastValues from the pool.
func AcquireFastValues() *FastValues {
	return fastValuesPool.Get().(*FastValues)
}

// ReleaseFastValues returns a FastValues to the pool.
func ReleaseFastValues(v *FastValues) {
	if v == nil {
		return
	}
	v.Reset()
	fastValuesPool.Put(v)
}

// Reset clears the FastValues for reuse.
func (v *FastValues) Reset() {
	for k := range v.data {
		delete(v.data, k)
	}
	for k := range v.keyIndex {
		delete(v.keyIndex, k)
	}
}

// Add adds a value to the key.
func (v *FastValues) Add(key, value string) {
	v.data[key] = append(v.data[key], value)
	v.keyIndex[key] = struct{}{}
}

// Set sets the key to a single value.
func (v *FastValues) Set(key, value string) {
	v.data[key] = []string{value}
	v.keyIndex[key] = struct{}{}
}

// Get returns the first value for the key.
func (v *FastValues) Get(key string) string {
	if vals, ok := v.data[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// GetAll returns all values for the key.
func (v *FastValues) GetAll(key string) []string {
	return v.data[key]
}

// Del deletes the key.
func (v *FastValues) Del(key string) {
	delete(v.data, key)
	delete(v.keyIndex, key)
}

// Has returns true if the key exists.
func (v *FastValues) Has(key string) bool {
	_, ok := v.keyIndex[key]
	return ok
}

// Len returns the number of keys.
func (v *FastValues) Len() int {
	return len(v.data)
}

// Encode encodes the values to a URL query string.
// Uses pre-allocated buffer for better performance.
func (v *FastValues) Encode() string {
	if len(v.data) == 0 {
		return ""
	}
	// Pre-allocate buffer: estimate ~32 bytes per key-value pair
	var sb strings.Builder
	sb.Grow(len(v.data) * 32)
	first := true
	for key, values := range v.data {
		escapedKey := url.QueryEscape(key)
		for _, value := range values {
			if !first {
				sb.WriteByte('&')
			}
			sb.WriteString(escapedKey)
			sb.WriteByte('=')
			sb.WriteString(url.QueryEscape(value))
			first = false
		}
	}
	return sb.String()
}

// Keys returns all keys.
func (v *FastValues) Keys() []string {
	keys := make([]string, 0, len(v.data))
	for k := range v.data {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a copy of the underlying data.
func (v *FastValues) Values() map[string][]string {
	result := make(map[string][]string, len(v.data))
	for k, vals := range v.data {
		result[k] = append([]string(nil), vals...)
	}
	return result
}

// ToURLValues converts to standard url.Values.
func (v *FastValues) ToURLValues() url.Values {
	result := make(url.Values, len(v.data))
	for k, vals := range v.data {
		result[k] = append([]string(nil), vals...)
	}
	return result
}

// ============================================================================
// Thread-safe Values methods (original)
// ============================================================================
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
