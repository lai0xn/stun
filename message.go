package stunlib

// Message represents a full STUN message, including its header and attributes.
type Message struct {
  header Header
	Attributes Attributes
}

func NewMessage(buff []byte) *Message {
  header := decodeHeader(buff)
  attributes := decodeAttrs(buff[20:],int(header.Length))
  return &Message {
    header: *header,
    Attributes: attributes,
  }
}


// DecodeAttrs decodes multiple STUN attributes from the given byte buffer.
// It iterates through the buffer, decoding each attribute and adding it to a slice.
func decodeAttrs(buff []byte, length int) []Attribute {
	offset := 0
	var attrs []Attribute

	// Loop through the buffer until the entire length is processed
	for offset < length {
		// Decode the current STUN attribute starting at the current offset
		attr := DecodeStunAttr(buff[offset:])

		// Append the decoded attribute to the slice
		attrs = append(attrs, attr)

		// Move the offset to the start of the next attribute
		// Each attribute has a 4-byte header (type + length) plus the padded value
		offset += 4 + attr.PaddedLength
	}

	// Return the slice of decoded attributes
	return attrs
}


func (m *Message) Encode() {

}
