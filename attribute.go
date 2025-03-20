package stunlib

// Attribute represents a STUN message attribute.
type Attribute struct {
	Length      uint16   // Length of the attribute value
	Type        StunAttribute // Type of the attribute (e.g., MAPPED-ADDRESS, USERNAME)
	PaddedLength int     // Length of the attribute value after padding (must be a multiple of 4)
	Value       []byte   // The value of the attribute (could be IP address, username, etc.)
}

// DecodeStunAttr decodes a single STUN attribute from the given byte buffer.
// The STUN attribute format is as follows:
func DecodeStunAttr(buff []byte) Attribute {
  // Extract the attribute type (first 2 bytes)
  attrType := StunAttribute(uint16(buff[0])<<8 | uint16(buff[1]))

  // Extract the attribute length (next 2 bytes)
  attrLen := uint16(buff[2])<<8 | uint16(buff[3])

  // Calculate the padded length of the attribute value
  // STUN attributes are padded to a multiple of 4 bytes
  paddedLen := int(attrLen)
  if paddedLen%4 != 0 {
    paddedLen = paddedLen + 4 - (paddedLen % 4)
  }

  return Attribute{
    Type:      attrType,
    Length:    attrLen,
    Value:     buff[4 : 4+paddedLen],
    PaddedLength: paddedLen,
  }
}

