package main

import (
	"time"

	stunlib "github.com/lai0xn/stun"
)

func main() {
	// Create a custom logger with JSON format for production
	logger := stunlib.NewLogger(stunlib.LoggerConfig{
		Level:      stunlib.InfoLevel,
		Format:     "json",
		Output:     "stdout",
		ShowCaller: true,
	})

	srv := stunlib.NewServer(stunlib.ServerConfig{
		Addr:    "127.0.0.1",
		Port:    "3478",
		Timeout: 30 * time.Second,
		Logger:  logger,
	})
	srv.Listen()
}
