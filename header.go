package stunlib

// Header represents the STUN message header.
type Header struct {
	Type          MessageType // Type of STUN message (e.g., Binding Request, Binding Response)
	Length        uint16      // Length of the message or attribute data
	MagicCookie   uint32
	TransactionID [12]byte // 12-byte Transaction ID to uniquely identify the request/response
}

// DecodeHeader takes a byte slice (buff) and decodes it into a STUN message header.
func decodeHeader(buff []byte) *Header {
	// Create a new Header object to store the decoded values
	header := new(Header)

	// Decode the STUN message type (2 bytes) and assign it to header.Type
	// Combine the first byte and second byte into a uint16 value using bitwise shifting
	// The first byte is shifted left by 8 bits, then the second byte is OR-ed with it
	header.Type = MessageType(uint16(buff[0])<<8 | uint16(buff[1]))

	// Decode the message length (2 bytes) and assign it to header.Length
	// Combine the third byte and fourth byte into a uint16 value using bitwise shifting
	header.Length = uint16(buff[2])<<8 | uint16(buff[3])

	// Decode the MagicCookie (4 bytes) and assign it to header.MagicCookie
	// MagicCookie is a fixed 4-byte value, so we combine 4 bytes (from index 4 to 7)
	// into a uint32 value using bitwise shifting and OR-ing the individual bytes
	header.MagicCookie = uint32(uint32(buff[4])<<24 | uint32(buff[5])<<16 | uint32(buff[6])<<8 | uint32(buff[7]))

	// Copy the remaining bytes (Transaction ID) into the header.TransactionID field
	// The TransactionID is 12 bytes long, so we copy from index 8 to the end of the buffer
	copy(header.TransactionID[:], buff[8:])

	// Return the decoded header
	return header
}

func encodeHeader(header Header) []byte {
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

func (h *Header) Encode() []byte {
  return encodeHeader(*h)
}
