package stun

import "errors"

// STUN Message Types
type MessageType uint16

type Attributes []Attribute

const headrLength = 20

const (
	// BindingRequest represents the Binding Request message type (0x0001),
	// which is used by the client to initiate a STUN transaction.
	// It requests the server to return the client's mapped address and port.
	BindingRequest MessageType = 0x0001

	// BindingResponse represents the Binding Response message type (0x0101),
	// which is sent by the STUN server in response to a Binding Request.
	// It contains the client's mapped address and port, allowing NAT traversal.
	BindingResponse MessageType = 0x0101

	// ErrorResponse represents the Error Response message type (0x0111),
	// which is sent by the STUN server when there is an error processing the request.
	// It includes an error code and description to notify the client of the issue.
	ErrorResponse MessageType = 0x0111
)

// STUN StunAttributes
type StunAttribute uint16

// MagicCookie used in STUN messages to distinguish it from other protocols
const magicCookie uint32 = 0x2112A442

// STUN Message StunAttributes
const (
	// MappedAddress represents the MAPPED-ADDRESS attribute (0x0001),
	// which indicates the IP address and port used by the client in NAT traversal.
	MappedAddress StunAttribute = 0x0001

	// Username represents the USERNAME attribute (0x0006),
	// which is used for authentication purposes in STUN messages.
	Username StunAttribute = 0x0006

	// MessageIntegrity represents the MESSAGE-INTEGRITY attribute (0x0008),
	// which provides a way to verify the integrity of the message by hashing its contents.
	MessageIntegrity StunAttribute = 0x0008

	// ErrorCode represents the ERROR-CODE attribute (0x0009),
	// which contains an error code and description when an error occurs in STUN processing.
	ErrorCode StunAttribute = 0x0009

	// UnknownStunAttributes represents the UNKNOWN-ATTRIBUTES attribute (0x000A),
	// which lists any attributes in the message that are not understood by the receiver.
	UnknownStunAttributes StunAttribute = 0x000A

	// Realm represents the REALM attribute (0x0014),
	// which is used for realm-based authentication (often with the NONCE attribute).
	Realm StunAttribute = 0x0014

	// Nonce represents the NONCE attribute (0x0015),
	// which is used for nonce-based authentication and to prevent replay attacks.
	Nonce StunAttribute = 0x0015

	// XORMappedAddress represents the XOR-MAPPED-ADDRESS attribute (0x0020),
	// which is similar to MAPPED-ADDRESS but uses XOR to obscure the actual IP address for added security.
	XORMappedAddress StunAttribute = 0x0020
)

var (
	ErrAttrNotFound  = errors.New("attribute not found")
	ErrShortBuffer   = errors.New("buffer too short for reading")
	ErrInvalidCookie = errors.New("invalid magic cookie")
	ErrShortWrite    = errors.New("short byte write")
)

// StunAttribute Lengths, attributes with 0 as value have variable lengths
const (
	MappedAddressLength         = 8  // 8 bytes for MAPPED-ADDRESS (IPv4 Value only)
	MessageIntegrityLength      = 20 // 20 bytes for MESSAGE-INTEGRITY (SHA1 HMAC digest)
	ErrorCodeLength             = 4  // 4 bytes minimal for ERROR-CODE (not including reason phrase)
	UnknownStunAttributesLength = 0  // Unknown attributes are variable length
	RealmLength                 = 0  // REALM is variable length
	NonceLength                 = 0  // NONCE is variable length
	XORMappedAddressLength      = 8  // 8 bytes for XOR-MAPPED-ADDRESS (IPv4 Value only)
)

// String returns the string representation of the MessageType
func (mt MessageType) String() string {
	switch mt {
	case BindingRequest:
		return "BindingRequest"
	case BindingResponse:
		return "BindingResponse"
	case ErrorResponse:
		return "ErrorResponse"
	default:
		return "Unknown"
	}
}
