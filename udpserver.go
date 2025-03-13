package stunlib

import (
	"log"
	"net"
	"time"
)

type udpServer struct {
	addr    string
	port    string
	timeout time.Duration
}

func (s udpServer) Listen() error {

	addr := net.JoinHostPort(s.addr, s.port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		return err
	}
	log.Println("server started listening on ", addr)
	conn, err := net.ListenUDP("udp4", udpAddr)

	if err != nil {
		return err
	}

	defer conn.Close()

	for {
		go HandleUDPConn(*conn)

	}
}

func HandleUDPConn(conn net.UDPConn) {
	buff := make([]byte, 1024)
	_, _, err := conn.ReadFromUDP(buff)
	if err != nil {
		return
	}
	_ = DecodeHeader(buff)
}

func (s udpServer) Shutdown() error {
	return nil
}
