package stunlib

import (
	"fmt"
	"net"
)

// DecodeHeader takes a byte slice (buff) and decodes it into a STUN message header.
func DecodeHeader(buff []byte) *Header {
	// Create a new Header object to store the decoded values
	header := new(Header)

	// Decode the STUN message type (2 bytes) and assign it to header.Type
	// Combine the first byte and second byte into a uint16 value using bitwise shifting
	// The first byte is shifted left by 8 bits, then the second byte is OR-ed with it
	header.Type = STUN_MESSAGE_TYPE(uint16(buff[0])<<8 | uint16(buff[1]))

	// Decode the message length (2 bytes) and assign it to header.Length
	// Combine the third byte and fourth byte into a uint16 value using bitwise shifting
	header.Length = uint16(buff[2])<<8 | uint16(buff[3])

	// Decode the MagicCookie (4 bytes) and assign it to header.MagicCookie
	// MagicCookie is a fixed 4-byte value, so we combine 4 bytes (from index 4 to 7)
	// into a uint32 value using bitwise shifting and OR-ing the individual bytes
	header.MagicCookie = uint32(uint32(buff[4])<<8 | uint32(buff[5])<<8 | uint32(buff[6])<<8 | uint32(buff[7]))

	// Copy the remaining bytes (Transaction ID) into the header.TransactionID field
	// The TransactionID is 12 bytes long, so we copy from index 8 to the end of the buffer
	copy(header.TransactionID[:], buff[8:])

	// Return the decoded header
	return header
}

func EncodeHeader(header Header) []byte {
	// Allocate a 20-byte slice for the STUN header (2 + 2 + 4 + 12 bytes)
	buff := make([]byte, 20)

	// Encode the message type (2 bytes)
	buff[0] = byte(header.Type >> 8)   // High byte
	buff[1] = byte(header.Type & 0xff) // Low byte

	// Encode the message length (2 bytes)
	buff[2] = byte(header.Length >> 8)   // High byte
	buff[3] = byte(header.Length & 0xff) // Low byte

	// Encode the magic cookie (4 bytes)
	buff[4] = byte(header.MagicCookie >> 24) // Highest byte
	buff[5] = byte(header.MagicCookie >> 16)
	buff[6] = byte(header.MagicCookie >> 8)
	buff[7] = byte(header.MagicCookie & 0xff) // Lowest byte

	// Copy the transaction ID (12 bytes)
	copy(buff[8:], header.TransactionID[:])

	return buff
}

// DecodeStunAttr decodes a single STUN attribute from the given byte buffer.
// The STUN attribute format is as follows:
func DecodeStunAttr(buff []byte) Attr {
	// Extract the attribute type (first 2 bytes)
	attrType := STUN_ATTR(uint16(buff[0])<<8 | uint16(buff[1]))

	// Extract the attribute length (next 2 bytes)
	attrLen := uint16(buff[2])<<8 | uint16(buff[3])

	// Calculate the padded length of the attribute value
	// STUN attributes are padded to a multiple of 4 bytes
	paddedLen := int(attrLen)
  fmt.Println(attrType == ATTR_XOR_MAPPED_ADDRESS)
	if paddedLen%4 != 0 {
		paddedLen = paddedLen + 4 - (paddedLen % 4)
	}

	return Attr{
		Type:      attrType,
		Length:    attrLen,
		Value:     buff[4 : 4+paddedLen],
		PaddedLen: paddedLen,
	}
}

// DecodeAttrs decodes multiple STUN attributes from the given byte buffer.
// It iterates through the buffer, decoding each attribute and adding it to a slice.
func DecodeAttrs(buff []byte, length int) []Attr {
	offset := 0
	var attrs []Attr

	// Loop through the buffer until the entire length is processed
	for offset < length {
		// Decode the current STUN attribute starting at the current offset
		attr := DecodeStunAttr(buff[offset:])

		// Append the decoded attribute to the slice
		attrs = append(attrs, attr)

		// Move the offset to the start of the next attribute
		// Each attribute has a 4-byte header (type + length) plus the padded value
		offset += 4 + attr.PaddedLen
	}

	// Return the slice of decoded attributes
	return attrs
}

// SerializeAddr takes an ip and port and encodes into a byte slice
func SerializeAddr(ip net.IP, port uint16) ([]byte, error) {
	// Check if the IP is IPv4
	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("invalid IPv4 address")
	}

	// Allocate the byte slice for the Mapped Address (8 bytes total)
	mappedAddress := make([]byte, 8)

	// Family: IPv4 (0x01)
	mappedAddress[0] = 0x01

	// Copy the 4-byte IPv4 address into the mapped address
	copy(mappedAddress[1:5], ipv4)

	// Serialize the 16-bit port into two bytes (big-endian)
	mappedAddress[5] = uint8(port >> 8)
	mappedAddress[6] = uint8(port & 0xff)

	// Return the serialized mapped address
	return mappedAddress, nil
}
