package main

import stunlib "github.com/lai0xn/stun"

func main() {
	srv := stunlib.NewServer(stunlib.ServerConfig{
		Addr:  "127.0.0.1",
		Port:  "3478",
		Proto: "udp",
	})
	srv.Listen()
}
