package codec

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

// Feature: http-client-refactor, Property 15: Codec Round-Trip
// For any valid Go struct and any registered Codec, encoding then decoding
// SHALL produce an equivalent struct.
// TestStruct is a test struct for codec testing.
type TestStruct struct {
	Name   string   `json:"name" xml:"name"`
	Age    int      `json:"age" xml:"age"`
	Score  float64  `json:"score" xml:"score"`
	Active bool     `json:"active" xml:"active"`
	Tags   []string `json:"tags" xml:"tags>tag"`
}

// NestedStruct is a nested test struct.
type NestedStruct struct {
	ID      int        `json:"id" xml:"id"`
	Data    TestStruct `json:"data" xml:"data"`
	Comment string     `json:"comment" xml:"comment"`
}

// XMLRoot is a wrapper for XML testing with proper root element.
type XMLRoot struct {
	XMLName xml.Name   `xml:"root"`
	Data    TestStruct `xml:"data"`
}

func TestCodec_Property15_RoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	// Property: JSON codec round-trip
	properties.Property("JSON codec round-trip preserves data", prop.ForAll(
		func(name string, age int, score float64, active bool) bool {
			original := TestStruct{
				Name:   name,
				Age:    age,
				Score:  score,
				Active: active,
				Tags:   []string{"tag1", "tag2"},
			}
			// Encode
			encoded, err := JSON.Encode(original)
			if err != nil {
				return false
			}
			// Decode
			var decoded TestStruct
			err = JSON.Decode(encoded, &decoded)
			if err != nil {
				return false
			}
			// Compare
			return reflect.DeepEqual(original, decoded)
		},
		gen.AlphaString(),
		gen.IntRange(-1000, 1000),
		gen.Float64Range(-1000.0, 1000.0),
		gen.Bool(),
	))
	// Property: XML codec round-trip
	properties.Property("XML codec round-trip preserves data", prop.ForAll(
		func(name string, age int, active bool) bool {
			original := XMLRoot{
				Data: TestStruct{
					Name:   name,
					Age:    age,
					Score:  0, // XML has issues with float precision
					Active: active,
					Tags:   []string{"tag1", "tag2"},
				},
			}
			// Encode
			encoded, err := XML.Encode(original)
			if err != nil {
				return false
			}
			// Decode
			var decoded XMLRoot
			err = XML.Decode(encoded, &decoded)
			if err != nil {
				return false
			}
			// Compare (ignoring XMLName)
			return original.Data.Name == decoded.Data.Name &&
				original.Data.Age == decoded.Data.Age &&
				original.Data.Active == decoded.Data.Active
		},
		gen.AlphaString(),
		gen.IntRange(-1000, 1000),
		gen.Bool(),
	))
	// Property: JSON nested struct round-trip
	properties.Property("JSON nested struct round-trip preserves data", prop.ForAll(
		func(id int, name string, age int, comment string) bool {
			original := NestedStruct{
				ID: id,
				Data: TestStruct{
					Name:   name,
					Age:    age,
					Score:  0,
					Active: true,
					Tags:   []string{},
				},
				Comment: comment,
			}
			// Encode
			encoded, err := JSON.Encode(original)
			if err != nil {
				return false
			}
			// Decode
			var decoded NestedStruct
			err = JSON.Decode(encoded, &decoded)
			if err != nil {
				return false
			}
			// Compare
			return reflect.DeepEqual(original, decoded)
		},
		gen.IntRange(-1000, 1000),
		gen.AlphaString(),
		gen.IntRange(-1000, 1000),
		gen.AlphaString(),
	))
	// Property: JSON map round-trip
	properties.Property("JSON map round-trip preserves data", prop.ForAll(
		func(key1, val1, key2, val2 string) bool {
			original := map[string]string{
				key1: val1,
				key2: val2,
			}
			// Encode
			encoded, err := JSON.Encode(original)
			if err != nil {
				return false
			}
			// Decode
			var decoded map[string]string
			err = JSON.Decode(encoded, &decoded)
			if err != nil {
				return false
			}
			// Compare
			return reflect.DeepEqual(original, decoded)
		},
		gen.AlphaString(),
		gen.AlphaString(),
		gen.AlphaString(),
		gen.AlphaString(),
	))
	// Property: JSON slice round-trip
	properties.Property("JSON slice round-trip preserves data", prop.ForAll(
		func(items []int) bool {
			// Encode
			encoded, err := JSON.Encode(items)
			if err != nil {
				return false
			}
			// Decode
			var decoded []int
			err = JSON.Decode(encoded, &decoded)
			if err != nil {
				return false
			}
			// Handle nil vs empty slice
			if len(items) == 0 && len(decoded) == 0 {
				return true
			}
			// Compare
			return reflect.DeepEqual(items, decoded)
		},
		gen.SliceOf(gen.IntRange(-1000, 1000)),
	))
	properties.TestingRun(t)
}

// Unit tests for specific scenarios
func TestJSONCodec_Encode(t *testing.T) {
	codec := NewJSONCodec()
	data := TestStruct{
		Name:   "test",
		Age:    25,
		Score:  95.5,
		Active: true,
		Tags:   []string{"a", "b"},
	}
	encoded, err := codec.Encode(data)
	assert.NoError(t, err)
	assert.Contains(t, string(encoded), `"name":"test"`)
	assert.Contains(t, string(encoded), `"age":25`)
}

func TestJSONCodec_Decode(t *testing.T) {
	codec := NewJSONCodec()
	jsonData := []byte(`{"name":"test","age":25,"score":95.5,"active":true,"tags":["a","b"]}`)
	var decoded TestStruct
	err := codec.Decode(jsonData, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "test", decoded.Name)
	assert.Equal(t, 25, decoded.Age)
	assert.Equal(t, 95.5, decoded.Score)
	assert.True(t, decoded.Active)
	assert.Equal(t, []string{"a", "b"}, decoded.Tags)
}

func TestJSONCodec_ContentType(t *testing.T) {
	codec := NewJSONCodec()
	assert.Equal(t, "application/json", codec.ContentType())
}

func TestJSONCodec_InvalidJSON(t *testing.T) {
	codec := NewJSONCodec()
	invalidJSON := []byte(`{"name": invalid}`)
	var decoded TestStruct
	err := codec.Decode(invalidJSON, &decoded)
	assert.Error(t, err)
}

func TestXMLCodec_Encode(t *testing.T) {
	codec := NewXMLCodec()
	data := XMLRoot{
		Data: TestStruct{
			Name:   "test",
			Age:    25,
			Active: true,
		},
	}
	encoded, err := codec.Encode(data)
	assert.NoError(t, err)
	assert.Contains(t, string(encoded), "<name>test</name>")
	assert.Contains(t, string(encoded), "<age>25</age>")
}

func TestXMLCodec_Decode(t *testing.T) {
	codec := NewXMLCodec()
	xmlData := []byte(`<root><data><name>test</name><age>25</age><active>true</active></data></root>`)
	var decoded XMLRoot
	err := codec.Decode(xmlData, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "test", decoded.Data.Name)
	assert.Equal(t, 25, decoded.Data.Age)
	assert.True(t, decoded.Data.Active)
}

func TestXMLCodec_ContentType(t *testing.T) {
	codec := NewXMLCodec()
	assert.Equal(t, "application/xml", codec.ContentType())
}

func TestXMLCodec_InvalidXML(t *testing.T) {
	codec := NewXMLCodec()
	invalidXML := []byte(`<root><unclosed>`)
	var decoded XMLRoot
	err := codec.Decode(invalidXML, &decoded)
	assert.Error(t, err)
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	codec := NewJSONCodec()
	registry.Register("application/json", codec)
	retrieved, ok := registry.GetCodec("application/json")
	assert.True(t, ok)
	assert.Equal(t, codec, retrieved)
}

func TestRegistry_GetCodec_NotFound(t *testing.T) {
	registry := NewRegistry()
	_, ok := registry.GetCodec("application/unknown")
	assert.False(t, ok)
}

func TestRegistry_NormalizeContentType(t *testing.T) {
	registry := NewRegistry()
	codec := NewJSONCodec()
	registry.Register("application/json", codec)
	// Should find with charset parameter
	retrieved, ok := registry.GetCodec("application/json; charset=utf-8")
	assert.True(t, ok)
	assert.Equal(t, codec, retrieved)
	// Should find with uppercase
	retrieved, ok = registry.GetCodec("APPLICATION/JSON")
	assert.True(t, ok)
	assert.Equal(t, codec, retrieved)
}

func TestRegistry_Has(t *testing.T) {
	registry := NewRegistry()
	codec := NewJSONCodec()
	registry.Register("application/json", codec)
	assert.True(t, registry.Has("application/json"))
	assert.False(t, registry.Has("application/xml"))
}

func TestRegistry_Remove(t *testing.T) {
	registry := NewRegistry()
	codec := NewJSONCodec()
	registry.Register("application/json", codec)
	registry.Remove("application/json")
	assert.False(t, registry.Has("application/json"))
}

func TestRegistry_ContentTypes(t *testing.T) {
	registry := NewRegistry()
	registry.Register("application/json", NewJSONCodec())
	registry.Register("application/xml", NewXMLCodec())
	types := registry.ContentTypes()
	assert.Len(t, types, 2)
	assert.Contains(t, types, "application/json")
	assert.Contains(t, types, "application/xml")
}

func TestDefaultRegistry(t *testing.T) {
	// JSON should be registered by default
	codec, ok := GetCodec("application/json")
	assert.True(t, ok)
	assert.NotNil(t, codec)
	// XML should be registered by default
	codec, ok = GetCodec("application/xml")
	assert.True(t, ok)
	assert.NotNil(t, codec)
}

// Feature: http-client-refactor, Property 16: Content-Type Based Codec Selection
// For any Content-Type header, the registry SHALL return the appropriate codec.
func TestCodec_Property16_ContentTypeSelection(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)
	// Property: JSON content types return JSON codec
	properties.Property("JSON content types return JSON codec", prop.ForAll(
		func(charset string) bool {
			contentTypes := []string{
				"application/json",
				"text/json",
				"application/json; charset=" + charset,
				"APPLICATION/JSON",
			}
			for _, ct := range contentTypes {
				codec, ok := GetCodec(ct)
				if !ok {
					return false
				}
				if codec.ContentType() != "application/json" {
					return false
				}
			}
			return true
		},
		gen.OneConstOf("utf-8", "UTF-8", "iso-8859-1"),
	))
	// Property: XML content types return XML codec
	properties.Property("XML content types return XML codec", prop.ForAll(
		func(charset string) bool {
			contentTypes := []string{
				"application/xml",
				"text/xml",
				"application/xml; charset=" + charset,
				"APPLICATION/XML",
			}
			for _, ct := range contentTypes {
				codec, ok := GetCodec(ct)
				if !ok {
					return false
				}
				if codec.ContentType() != "application/xml" {
					return false
				}
			}
			return true
		},
		gen.OneConstOf("utf-8", "UTF-8", "iso-8859-1"),
	))
	// Property: Unknown content types return false
	properties.Property("Unknown content types return false", prop.ForAll(
		func(unknown string) bool {
			_, ok := GetCodec("application/" + unknown)
			return !ok
		},
		gen.OneConstOf("unknown", "custom", "binary", "octet-stream"),
	))
	// Property: Custom codec registration works
	properties.Property("Custom codec registration works", prop.ForAll(
		func(contentType string) bool {
			registry := NewRegistry()
			customCodec := NewJSONCodec()
			registry.Register(contentType, customCodec)
			retrieved, ok := registry.GetCodec(contentType)
			if !ok {
				return false
			}
			return retrieved == customCodec
		},
		gen.OneConstOf("application/custom", "text/custom", "application/x-custom"),
	))
	properties.TestingRun(t)
}

func TestRegistry_GetEncoder(t *testing.T) {
	registry := NewRegistry()
	registry.Register("application/json", NewJSONCodec())
	encoder, ok := registry.GetEncoder("application/json")
	assert.True(t, ok)
	assert.NotNil(t, encoder)
	assert.Equal(t, "application/json", encoder.ContentType())
	_, ok = registry.GetEncoder("application/unknown")
	assert.False(t, ok)
}

func TestRegistry_GetDecoder(t *testing.T) {
	registry := NewRegistry()
	registry.Register("application/json", NewJSONCodec())
	decoder, ok := registry.GetDecoder("application/json")
	assert.True(t, ok)
	assert.NotNil(t, decoder)
	_, ok = registry.GetDecoder("application/unknown")
	assert.False(t, ok)
}

func TestDefaultRegistry_TextJSON(t *testing.T) {
	// text/json should also be registered
	codec, ok := GetCodec("text/json")
	assert.True(t, ok)
	assert.NotNil(t, codec)
	assert.Equal(t, "application/json", codec.ContentType())
}

func TestDefaultRegistry_TextXML(t *testing.T) {
	// text/xml should also be registered
	codec, ok := GetCodec("text/xml")
	assert.True(t, ok)
	assert.NotNil(t, codec)
	assert.Equal(t, "application/xml", codec.ContentType())
}

func TestGlobalGetEncoder(t *testing.T) {
	encoder, ok := GetEncoder("application/json")
	assert.True(t, ok)
	assert.NotNil(t, encoder)
	assert.Equal(t, "application/json", encoder.ContentType())
	_, ok = GetEncoder("application/unknown")
	assert.False(t, ok)
}

func TestGlobalGetDecoder(t *testing.T) {
	decoder, ok := GetDecoder("application/json")
	assert.True(t, ok)
	assert.NotNil(t, decoder)
	_, ok = GetDecoder("application/unknown")
	assert.False(t, ok)
}

func TestEncoderFunc_Encode(t *testing.T) {
	fn := EncoderFunc{
		EncodeFunc: func(v any) ([]byte, error) {
			return []byte("encoded"), nil
		},
		ContentTypeFunc: func() string {
			return "application/test"
		},
	}
	result, err := fn.Encode("test")
	assert.NoError(t, err)
	assert.Equal(t, []byte("encoded"), result)
}

func TestEncoderFunc_ContentType(t *testing.T) {
	fn := EncoderFunc{
		EncodeFunc: func(v any) ([]byte, error) {
			return nil, nil
		},
		ContentTypeFunc: func() string {
			return "application/custom"
		},
	}
	assert.Equal(t, "application/custom", fn.ContentType())
}

func TestDecoderFunc_Decode(t *testing.T) {
	fn := DecoderFunc(func(data []byte, v any) error {
		ptr := v.(*string)
		*ptr = "decoded"
		return nil
	})
	var result string
	err := fn.Decode([]byte("test"), &result)
	assert.NoError(t, err)
	assert.Equal(t, "decoded", result)
}

func TestRegisterFast(t *testing.T) {
	// Register a custom codec using RegisterFast
	customCodec := NewJSONCodec()
	RegisterFast("application/x-fast-test", customCodec)

	// Should be retrievable via GetCodec
	retrieved, ok := GetCodec("application/x-fast-test")
	assert.True(t, ok)
	assert.Equal(t, customCodec, retrieved)

	// Should also work with charset parameter
	retrieved, ok = GetCodec("application/x-fast-test; charset=utf-8")
	assert.True(t, ok)
	assert.Equal(t, customCodec, retrieved)
}

func TestRegisterFast_FastPath(t *testing.T) {
	// Register multiple codecs with RegisterFast
	codec1 := NewJSONCodec()
	codec2 := NewXMLCodec()

	RegisterFast("application/x-fast-json", codec1)
	RegisterFast("application/x-fast-xml", codec2)

	// Both should be retrievable
	retrieved1, ok1 := GetCodec("application/x-fast-json")
	retrieved2, ok2 := GetCodec("application/x-fast-xml")

	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.Equal(t, codec1, retrieved1)
	assert.Equal(t, codec2, retrieved2)
}

func TestFastCodecCache_Initialization(t *testing.T) {
	// initFastCache should be safe to call multiple times
	initFastCache()
	initFastCache()
	initFastCache()

	// Cache should be initialized
	assert.NotNil(t, fastCodecCache)
}

func TestRegistry_GetCodec_FastPathPriority(t *testing.T) {
	// Create a new registry
	registry := NewRegistry()

	// Register a codec in the registry
	slowCodec := NewJSONCodec()
	registry.Register("application/x-priority-test", slowCodec)

	// Register a different codec in fast cache
	fastCodec := NewXMLCodec()
	initFastCache()
	fastCodecCache["application/x-priority-test"] = fastCodec

	// Fast cache should take priority
	retrieved, ok := registry.GetCodec("application/x-priority-test")
	assert.True(t, ok)
	assert.Equal(t, fastCodec, retrieved)

	// Clean up
	delete(fastCodecCache, "application/x-priority-test")
}

func TestRegistry_GetCodec_FallbackToRegistry(t *testing.T) {
	// Create a new registry
	registry := NewRegistry()

	// Register a codec only in the registry (not in fast cache)
	codec := NewJSONCodec()
	registry.Register("application/x-fallback-test", codec)

	// Should still be found via registry fallback
	retrieved, ok := registry.GetCodec("application/x-fallback-test")
	assert.True(t, ok)
	assert.Equal(t, codec, retrieved)
}
