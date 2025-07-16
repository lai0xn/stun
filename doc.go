// Package stun implements the Session Traversal Utilities for NAT (STUN) protocol
// as defined in RFC 5389. This library provides both client and server implementations
// for NAT traversal, which is essential for peer-to-peer networking applications.
//
// STUN is a protocol that allows clients to discover their public IP address and
// the type of NAT they are behind. This information is crucial for establishing
// peer-to-peer connections in applications like WebRTC, VoIP, and online gaming.
//
// Key Features:
//   - Full STUN protocol implementation (RFC 5389)
//   - Client and server implementations
//   - Support for XOR-MAPPED-ADDRESS attribute
//   - Structured logging with configurable levels
//   - Comprehensive error handling
//   - Easy-to-use API
//
// Basic Usage:
//
//	// Create a STUN client
//	client := stun.NewClient("stun.l.google.com:19302")
//
//	// Send a binding request
//	msg, err := client.Dial(&stun.Message{
//		Header: stun.Header{
//			Type: stun.BindingRequest,
//		},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Get the XOR mapped address
//	xorAddr, err := msg.GetXorAddr()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Public IP: %s:%d\n", xorAddr.IP, xorAddr.Port)
//
// Server Usage:
//
//	// Create a STUN server
//	server := stun.NewServer(stun.ServerConfig{
//		Addr: "127.0.0.1",
//		Port: "3478",
//		Logger: stun.NewDefaultLogger(),
//	})
//
//	// Start listening
//	if err := server.Listen(); err != nil {
//		log.Fatal(err)
//	}
//
// Logging:
//
// The library includes a comprehensive logging system that can be configured
// for different environments:
//
//	// Development logging (text format)
//	logger := stun.NewLogger(stun.LoggerConfig{
//		Level:      stun.DebugLevel,
//		Format:     "text",
//		ShowCaller: true,
//	})
//
//	// Production logging (JSON format)
//	logger := stun.NewLogger(stun.LoggerConfig{
//		Level:      stun.InfoLevel,
//		Format:     "json",
//		ShowCaller: false,
//	})
//
// Protocol Details:
//
// This implementation supports the core STUN protocol features:
//   - Binding Request/Response messages
//   - XOR-MAPPED-ADDRESS attribute
//   - Transaction ID generation and validation
//   - Magic cookie validation
//   - IPv4 address support
//
// The library follows RFC 5389 specifications and includes proper error handling
// for malformed messages, network issues, and protocol violations.
//
// For more information about the STUN protocol, see RFC 5389:
// https://tools.ietf.org/html/rfc5389
package stun
