package codec

import (
	"strings"
	"sync"
)

// fastCodecCache stores commonly used codecs for lock-free access.
// This is populated lazily on first access to avoid init order issues.
var (
	fastCodecCache     map[string]Codec
	fastCodecCacheOnce sync.Once
)

// initFastCache initializes the fast codec cache lazily.
func initFastCache() {
	fastCodecCacheOnce.Do(func() {
		fastCodecCache = make(map[string]Codec, 8)
	})
}

// Registry manages codec registration by content type.
type Registry struct {
	mu     sync.RWMutex
	codecs map[string]Codec
}

// NewRegistry creates a new codec registry.
func NewRegistry() *Registry {
	return &Registry{
		codecs: make(map[string]Codec),
	}
}

// Register registers a codec for the given content type.
func (r *Registry) Register(contentType string, codec Codec) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.codecs[normalizeContentType(contentType)] = codec
}

// GetCodec returns the codec for the given content type.
// Uses fast path for commonly registered codecs to avoid lock contention.
func (r *Registry) GetCodec(contentType string) (Codec, bool) {
	normalized := normalizeContentType(contentType)

	// Fast path: check cache first (lock-free read after init)
	if fastCodecCache != nil {
		if codec, ok := fastCodecCache[normalized]; ok {
			return codec, true
		}
	}

	// Slow path: check registry with lock
	r.mu.RLock()
	defer r.mu.RUnlock()
	codec, ok := r.codecs[normalized]
	return codec, ok
}

// GetEncoder returns the encoder for the given content type.
func (r *Registry) GetEncoder(contentType string) (Encoder, bool) {
	codec, ok := r.GetCodec(contentType)
	if !ok {
		return nil, false
	}
	return codec, true
}

// GetDecoder returns the decoder for the given content type.
func (r *Registry) GetDecoder(contentType string) (Decoder, bool) {
	codec, ok := r.GetCodec(contentType)
	if !ok {
		return nil, false
	}
	return codec, true
}

// Has checks if a codec is registered for the given content type.
func (r *Registry) Has(contentType string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.codecs[normalizeContentType(contentType)]
	return ok
}

// Remove removes the codec for the given content type.
func (r *Registry) Remove(contentType string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.codecs, normalizeContentType(contentType))
}

// ContentTypes returns all registered content types.
func (r *Registry) ContentTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	types := make([]string, 0, len(r.codecs))
	for ct := range r.codecs {
		types = append(types, ct)
	}
	return types
}

// normalizeContentType normalizes a content type string.
// It extracts the base MIME type without parameters.
func normalizeContentType(contentType string) string {
	// Remove parameters (e.g., "application/json; charset=utf-8" -> "application/json")
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = contentType[:idx]
	}
	return strings.TrimSpace(strings.ToLower(contentType))
}

// DefaultRegistry is the global default codec registry.
var DefaultRegistry = NewRegistry()

func init() {
	// Register XML codecs (JSON is registered in sonic.go or sonic_fallback.go)
	DefaultRegistry.Register("application/xml", NewXMLCodec())
	DefaultRegistry.Register("text/xml", NewXMLCodec())
}

// RegisterFast registers a codec in both the default registry and fast cache.
// Use this for frequently accessed content types to avoid lock contention.
func RegisterFast(contentType string, c Codec) {
	initFastCache()
	normalized := normalizeContentType(contentType)
	fastCodecCache[normalized] = c
	DefaultRegistry.Register(contentType, c)
}

// Register registers a codec in the default registry.
func Register(contentType string, codec Codec) {
	DefaultRegistry.Register(contentType, codec)
}

// GetCodec returns a codec from the default registry.
func GetCodec(contentType string) (Codec, bool) {
	return DefaultRegistry.GetCodec(contentType)
}

// GetEncoder returns an encoder from the default registry.
func GetEncoder(contentType string) (Encoder, bool) {
	return DefaultRegistry.GetEncoder(contentType)
}

// GetDecoder returns a decoder from the default registry.
func GetDecoder(contentType string) (Decoder, bool) {
	return DefaultRegistry.GetDecoder(contentType)
}
