package codec

import (
	"encoding/xml"
)

const (
	// ContentTypeXML is the MIME type for XML.
	ContentTypeXML = "application/xml"
	// ContentTypeTextXML is the text MIME type for XML.
	ContentTypeTextXML = "text/xml"
)

// XMLCodec implements Codec for XML encoding/decoding.
type XMLCodec struct{}

// NewXMLCodec creates a new XML codec.
func NewXMLCodec() *XMLCodec {
	return &XMLCodec{}
}

// Encode encodes the value to XML bytes.
func (c *XMLCodec) Encode(v any) ([]byte, error) {
	return xml.Marshal(v)
}

// Decode decodes XML bytes into the destination.
func (c *XMLCodec) Decode(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

// ContentType returns the XML MIME type.
func (c *XMLCodec) ContentType() string {
	return ContentTypeXML
}

// XML is the default XML codec instance.
var XML = NewXMLCodec()

func init() {
	// Register XML codec in the default registry
	Register(ContentTypeXML, XML)
	Register(ContentTypeTextXML, XML)
	// Also register common variations
	Register("application/xml; charset=utf-8", XML)
	Register("text/xml; charset=utf-8", XML)
}
