package stunlib

// STUN Message Types
type STUN_MESSAGE_TYPE uint16

const (
    // STUN_BINDING_REQUEST represents the Binding Request message type (0x0001),
    // which is used by the client to initiate a STUN transaction.
    // It requests the server to return the client's mapped address and port.
    STUN_BINDING_REQUEST  STUN_MESSAGE_TYPE = 0x0001

    // STUN_BINDING_RESPONSE represents the Binding Response message type (0x0101),
    // which is sent by the STUN server in response to a Binding Request.
    // It contains the client's mapped address and port, allowing NAT traversal.
    STUN_BINDING_RESPONSE STUN_MESSAGE_TYPE = 0x0101

    // STUN_ERROR_RESPONSE represents the Error Response message type  (0x0111),
    // which is sent by the STUN server when there is an error processing the request.
    // It includes an error code and description to notify the client of the issue.
    STUN_ERROR_RESPONSE   STUN_MESSAGE_TYPE = 0x0111
)

// STUN Attributes
type STUN_ATTR uint16

// Magic Cookie used in STUN messages to distinguish it from other protocols
const MAGIC_COOKIE uint32 = 0x2112A442

// STUN Message Attributes
const (
    // ATTR_MAPPED_ADDRESS represents the MAPPED-ADDRESS attribute (0x0001),
    // which indicates the IP address and port used by the client in NAT traversal.
    ATTR_MAPPED_ADDRESS 	STUN_ATTR = 0x0001

    // ATTR_USERNAME represents the USERNAME attribute (0x0006),
    // which is used for authentication purposes in STUN messages.
    ATTR_USERNAME       	STUN_ATTR = 0x0006

    // ATTR_MESSAGE_INTEGRITY represents the MESSAGE-INTEGRITY attribute (0x0008),
    // which provides a way to verify the integrity of the message by hashing its contents.
    ATTR_MESSAGE_INTEGRITY  STUN_ATTR = 0x0008

    // ATTR_ERROR_CODE represents the ERROR-CODE attribute (0x0009),
    // which contains an error code and description when an error occurs in STUN processing.
    ATTR_ERROR_CODE 	STUN_ATTR = 0x0009

    // ATTR_UKNOWN_ATTRS represents the UNKNOWN-ATTRIBUTES attribute (0x000A),
    // which lists any attributes in the message that are not understood by the receiver.
    ATTR_UKNOWN_ATTRS 	STUN_ATTR = 0x000A

    // ATTR_REALM represents the REALM attribute (0x0014),
    // which is used for realm-based authentication (often with the NONCE attribute).
    ATTR_REALM 		STUN_ATTR = 0x0014

    // ATTR_NONCE represents the NONCE attribute (0x0015),
    // which is used for nonce-based authentication and to prevent replay attacks.
    ATTR_NONCE 		STUN_ATTR = 0x0015

    // ATTR_XOR_MAPPED_ADDRESS represents the XOR-MAPPED-ADDRESS attribute (0x0020),
    // which is similar to MAPPED-ADDRESS but uses XOR to obscure the actual IP address for added security.
    ATTR_XOR_MAPPED_ADDRESS STUN_ATTR = 0x0020
)

// Attribute Lengths
const (
    LENGTH_MAPPED_ADDRESS      = 12  // 12 bytes for MAPPED-ADDRESS (IPv4 + port)
    LENGTH_MESSAGE_INTEGRITY   = 20  // 20 bytes for MESSAGE-INTEGRITY (SHA1)
    LENGTH_ERROR_CODE          = 4   // 4 bytes for ERROR-CODE
    LENGTH_UKNOWN_ATTRS        = 0   // Unknown attributes can be of variable length
    LENGTH_REALM               = 0   // REALM can be variable length
    LENGTH_NONCE               = 0   // NONCE can be variable length
    LENGTH_XOR_MAPPED_ADDRESS  = 12  // 12 bytes for XOR-MAPPED-ADDRESS (XOR encoded IPv4 + port)
)




// Header represents the STUN message header.
type Header struct {
    Type           STUN_MESSAGE_TYPE // Type of STUN message (e.g., Binding Request, Binding Response)
    Length         uint16       // Length of the message or attribute data
    MagicCookie    uint32
    TransactionID  [12]byte          // 12-byte Transaction ID to uniquely identify the request/response
}

// Attr represents a STUN message attribute.
type Attr struct {
    Length uint16// Length of the attribute value
    Type   STUN_ATTR // Type of the attribute (e.g., MAPPED-ADDRESS, USERNAME)
    Value  []byte // The value of the attribute (could be IP address, username, etc.)
}

// Message represents a full STUN message, including its header and attributes.
type Message struct {
    Header 
    Attrs  []Attr 
}


