package stunlib

import (
	"fmt"
	"net"
)

type IPFamily uint16

const IPV4 IPFamily = 0x01
const IPV6 IPFamily = 0x02

type XorMappedAddr struct {
  Family IPFamily
  IP net.IP
  port uint16
}


// SerializeAddr takes an ip and port and encodes into a byte slice
func SerializeAddr(addr XorMappedAddr) ([]byte, error) {
	// Check if the IP is IPv4
	ipv4 := addr.IP.To4()
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
	mappedAddress[5] = uint8(addr.port >> 8)
  
  mappedAddress[6] = uint8(addr.port & 0xff)

	// Return the serialized mapped address
	return mappedAddress, nil
}

// DecodeAddr takes an ip and port as bytes and decodes them into XorMappedAddr
func DecodeAddr(addr []byte) *XorMappedAddr {
  // Decode IP Family
  familly := addr[0]

  ip := make([]byte,4)
  net.IPv4(ip[0],ip[1],ip[2],ip[3])
  
  port := uint16(uint16(addr[5]) << 8 | uint16(addr[6]))
  
  return &XorMappedAddr{
    Family: IPFamily(familly),
    IP: ip,
    port: port,
  }

} 

