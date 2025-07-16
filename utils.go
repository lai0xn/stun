package stun

import (
	"crypto/rand"
	"fmt"
	"net"
)

func randomTransactionID() []byte {
	transactionID := make([]byte, 12)
	_, err := rand.Read(transactionID)
	if err != nil {
		// In case of error, return a zero-filled slice (fallback)
		return transactionID
	}
	return transactionID
}

func GetPortFromAddr(addr net.Addr) (int, error) {
	switch a := addr.(type) {
	case *net.TCPAddr:
		return a.Port, nil
	case *net.UDPAddr:
		return a.Port, nil
	case *net.UnixAddr:
		// Unix domain sockets don't have ports
		return 0, nil
	default:
		// Unknown address type
		return 0, fmt.Errorf("unsupported address type: %T", addr)
	}
}

// GetPortAndIPFromAddr extracts both port and IP from a net.Addr interface
func GetPortAndIPFromAddr(addr net.Addr) (int, net.IP, error) {
	switch a := addr.(type) {
	case *net.TCPAddr:
		return a.Port, a.IP, nil
	case *net.UDPAddr:
		return a.Port, a.IP, nil
	case *net.UnixAddr:
		// Unix domain sockets don't have ports or IPs in the same way
		return 0, nil, nil
	default:
		return 0, nil, fmt.Errorf("unsupported address type: %T", addr)
	}
}
