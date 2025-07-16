package stun

import (
	"encoding/binary"
	"fmt"
	"net"
)

type IPFamily uint16

const IPV4 IPFamily = 0x01
const IPV6 IPFamily = 0x02

//	0                   1                   2                   3
//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |0 0 0 0 0 0 0 0|    Family     |           Port                |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// |                 Address (32 bits or 128 bits)                 |
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
//	Figure 5: Format of MAPPED-ADDRESS Attribute
type XorMappedAddr struct {
	Family IPFamily
	IP     net.IP
	Port   uint16
}

// SerializeAddr takes an ip and Port and encodes into a byte slice
func serializeAddr(addr XorMappedAddr, transactionID [12]byte) ([]byte, error) {
	ipv4 := addr.IP.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("invalid IPv4 address")
	}

	buf := make([]byte, 8)
	buf[0] = 0x00 // Reserved
	buf[1] = byte(IPV4)

	// XOR Port
	xorPort := addr.Port ^ uint16(magicCookie>>16)
	buf[2] = byte(xorPort >> 8)
	buf[3] = byte(xorPort & 0xFF)

	// XOR IP
	magicCookieBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(magicCookieBytes, magicCookie)

	for i := 0; i < 4; i++ {
		buf[4+i] = ipv4[i] ^ magicCookieBytes[i]
	}

	return buf, nil
}

// DecodeAddr takes an ip and Port as bytes and decodes them into XorMappedAddr
func decodeAddr(addr []byte) *XorMappedAddr {

	// Decode IP Family
	// Skip the first reserved byte
	familly := addr[1]

	x := uint16(magicCookie >> 16)

	port := uint16(uint16(addr[2])<<8 | uint16(addr[3]) ^ x)

	ip := make([]byte, 4)

	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(addr[4:8])^magicCookie)

	return &XorMappedAddr{
		Family: IPFamily(familly),
		Port:   port,
		IP:     net.IP(ip),
	}

}
