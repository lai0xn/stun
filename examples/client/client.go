package main

import (
	"fmt"

	"github.com/lai0xn/stun"
)

func main() {
	// Create a client with debug logging
	logger := stun.NewLogger(stun.LoggerConfig{
		Level:      stun.DebugLevel,
		Format:     "text",
		Output:     "stdout",
		ShowCaller: false,
	})

	client := stun.NewClientWithLogger("stun.l.google.com:19302", logger)

	msg, err := client.Dial(&stun.Message{
		Header: stun.Header{
			Type: stun.BindingRequest,
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	xorAddr, err := msg.GetXorAddr()
	if err != nil {
		fmt.Printf("Error getting XOR address: %v\n", err)
		return
	}

	fmt.Printf("XOR Mapped Address: %s:%d\n", xorAddr.IP, xorAddr.Port)
}
