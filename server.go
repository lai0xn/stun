package stunlib

import "time"

type ServerConfig struct {
	Addr    string
	Port    string
	Proto   string
	Timeout time.Duration
}

type Server interface {
	Listen() error
	Shutdown() error
}

func NewServer(cfg ServerConfig) Server {
	return udpServer{
		addr: cfg.Addr,
		port: cfg.Port,
	}
}
