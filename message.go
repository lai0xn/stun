package stun

// Message represents a complete STUN message, including its header and attributes.
// A STUN message consists of a 20-byte header followed by zero or more attributes.
//
// The Message structure follows RFC 5389 specifications:
//   - Header: Contains message type, length, magic cookie, and transaction ID
//   - Attributes: Variable-length list of STUN attributes
//
// Example:
//
//	msg := &stun.Message{
//		Header: stun.Header{
//			Type: stun.BindingRequest,
//		},
//		Attributes: []stun.Attribute{},
//	}
type Message struct {
	Header     Header
	Attributes Attributes
}

// NewMessage creates a new Message by parsing the provided byte buffer.
// The buffer should contain a complete STUN message starting with the header.
//
// The function performs the following operations:
//   - Decodes the 20-byte header
//   - Validates the magic cookie
//   - Parses all attributes based on the message length
//   - Returns a fully populated Message structure
//
// Returns:
//   - *Message: The parsed STUN message
//   - error: Any error that occurred during parsing
//
// Example:
//
//	buff := []byte{...} // STUN message bytes
//	msg, err := stun.NewMessage(buff)
//	if err != nil {
//		log.Fatal(err)
//	}
func NewMessage(buff []byte) (*Message, error) {
	header, err := decodeHeader(buff)
	if err != nil {
		return nil, err
	}
	attributes := decodeAttrs(buff[20:], int(header.Length))
	return &Message{
		Header:     *header,
		Attributes: attributes,
	}, nil
}

// GetAttr searches for a specific attribute type in the message and returns it if found.
// This method iterates through all attributes in the message to find a match.
//
// Parameters:
//   - t: The StunAttribute type to search for
//
// Returns:
//   - *Attribute: The found attribute, or nil if not found
//   - bool: True if the attribute was found, false otherwise
//
// Example:
//
//	if attr, found := msg.GetAttr(stun.XORMappedAddress); found {
//		// Process the XOR-MAPPED-ADDRESS attribute
//		xorAddr := decodeAddr(attr.Value)
//		fmt.Printf("XOR Address: %s:%d\n", xorAddr.IP, xorAddr.Port)
//	}
func (m Message) GetAttr(t StunAttribute) (*Attribute, bool) {
	for _, attr := range m.Attributes {
		if attr.Type == t {
			return &attr, true
		}
	}
	return nil, false
}

// GetXorAddr extracts the XOR-MAPPED-ADDRESS attribute from the message.
// This method is specifically designed for handling binding responses and
// provides a convenient way to access the client's public IP address and port.
//
// The method checks if the message is a binding response and then looks for
// the XOR-MAPPED-ADDRESS attribute. If found, it decodes the address information.
//
// Returns:
//   - *XorMappedAddr: The decoded XOR mapped address, or nil if not found
//   - error: Any error that occurred during decoding
//
// Example:
//
//	msg, err := client.Dial(&stun.Message{
//		Header: stun.Header{Type: stun.BindingRequest},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	xorAddr, err := msg.GetXorAddr()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Public IP: %s:%d\n", xorAddr.IP, xorAddr.Port)
func (m Message) GetXorAddr() (*XorMappedAddr, error) {
	if m.Header.Type != BindingResponse {
		return nil, nil
	}
	if attr, ok := m.GetAttr(XORMappedAddress); ok {
		return decodeAddr(attr.Value), nil
	}
	return nil, ErrAttrNotFound
}

// decodeAttrs decodes multiple STUN attributes from the given byte buffer.
// It iterates through the buffer, decoding each attribute and adding it to a slice.
//
// The function processes attributes sequentially, using the length information
// in each attribute header to determine where the next attribute begins.
// Each attribute has a 4-byte header (type + length) followed by the attribute value.
//
// Parameters:
//   - buff: The byte buffer containing attribute data
//   - length: The total length of attribute data to process
//
// Returns:
//   - []Attribute: A slice of decoded STUN attributes
func decodeAttrs(buff []byte, length int) []Attribute {
	offset := 0
	var attrs []Attribute

	// Loop through the buffer until the entire length is processed
	for offset < length {
		// Decode the current STUN attribute starting at the current offset
		attr := DecodeAttr(buff[offset:])

		// Append the decoded attribute to the slice
		attrs = append(attrs, attr)

		// Move the offset to the start of the next attribute
		// Each attribute has a 4-byte header (type + length) plus the padded value
		offset += 4 + attr.PaddedLength
	}

	// Return the slice of decoded attributes
	return attrs
}

// Encode converts the Message to its binary representation.
// This method serializes the complete STUN message including header and all attributes.
//
// The encoding process:
//   - Encodes the 20-byte header
//   - Encodes each attribute in sequence
//   - Returns the complete binary message
//
// Returns:
//   - []byte: The encoded STUN message as a byte slice
//
// Example:
//
//	msg := &stun.Message{
//		Header: stun.Header{
//			Type: stun.BindingRequest,
//		},
//		Attributes: []stun.Attribute{},
//	}
//	encoded := msg.Encode()
//	// Send encoded message over network
func (m *Message) Encode() []byte {
	buff := make([]byte, m.Header.Length+20)
	copy(buff[0:20], m.Header.Encode())
	offset := 20
	for _, attr := range m.Attributes {
		encodedAttr := attr.Encode()
		attrSize := attr.PaddedLength + 4 // 4 bytes header + padded value length
		copy(buff[offset:offset+attrSize], encodedAttr)
		offset += attrSize
	}
	return buff
}
