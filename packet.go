package stun

import (
	"fmt"
	"net"
)

type Packet struct {
	con        *net.UDPConn
	sourceIP   net.IP
	message    *Message
	sourcePort uint16
  remoteIP net.IP
  remotePort uint16
}
func NewPacket(con *net.UDPConn, buff []byte,remoteAddr *net.UDPAddr) (*Packet, error) {
	msg, err := NewMessage(buff) // Assuming NewMessage is defined elsewhere
	if err != nil {
		return nil, err
	}

	// Get the local address from the connection
	addr := con.LocalAddr()
	if addr == nil {
		return nil, fmt.Errorf("failed to get local address from connection")
	}

	// Extract port and IP from the address
	localPort, localIP, err := GetPortAndIPFromAddr(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get port and IP: %v", err)
	}
  if remoteAddr == nil {
    return nil, fmt.Errorf("failed to get remote adress from connectoin")
  }
  port, ip ,err := GetPortAndIPFromAddr(remoteAddr)
  if err != nil {
    return nil, fmt.Errorf("faield to get port and ip :%v",err)
  }
	return &Packet{
		con:        con,
		sourceIP:   localIP,
		sourcePort: uint16(localPort),
		message:        msg,
	  remoteIP: ip,
    remotePort: uint16(port),
  }, nil
}


func (p *Packet) Write(buff []byte,remoteAddr *net.UDPAddr) (int, error) {
	msg, err := NewMessage(buff)
	if err != nil {
		return 0, err
	}
	n, err := p.con.WriteTo(msg.Encode(),remoteAddr)

	if err != nil {
		return 0, err
	}

	if n < (int(msg.Header.Length) + headrLength) {
		return n, ErrShortWrite
	}
	return n, nil
}

