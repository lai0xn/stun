package stunlib

import (
	"log"
	"net"
	"time"
)


type udpServer struct {
	addr string
	port string
	timeout time.Duration
}

func (s udpServer) Listen() error{
	
	addr := net.JoinHostPort(s.addr,s.port)
	udpAddr,err := net.ResolveUDPAddr("udp",addr)
	
	if err != nil {
		return err
	}
	log.Println("server started listening on ",addr)
	conn,err := net.ListenUDP("udp4",udpAddr)
	if err != nil {
		return err
	}
	HandleUDPConn(*conn)		
	return nil
}

func HandleUDPConn(conn net.UDPConn){
	defer conn.Close()
	buff := make([]byte,1024)
	header := Header{}
	for {
		_,_,err := conn.ReadFromUDP(buff)
		if err != nil {
			return
		}
		header.Type = STUN_MESSAGE_TYPE(uint16(buff[0]) << 8 | uint16(buff[1]))
	}	
}

func (s udpServer) Shutdown() error{
	return nil
}
