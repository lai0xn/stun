package stun

import (
	"net"
	"time"
)

// Server represents a STUN server that listens for binding requests from clients
// and responds with the client's public IP address and port information.
//
// The server implements the core STUN protocol functionality:
//   - Listening for UDP connections
//   - Parsing incoming STUN messages
//   - Generating XOR-MAPPED-ADDRESS responses
//   - Handling multiple concurrent clients
//   - Comprehensive logging and error handling
//
// Example:
//
//	server := stun.NewServer(stun.ServerConfig{
//		Addr: "127.0.0.1",
//		Port: "3478",
//		Logger: stun.NewDefaultLogger(),
//	})
//	if err := server.Listen(); err != nil {
//		log.Fatal(err)
//	}
type Server struct {
	addr    string
	port    string
	timeout time.Duration
	logger  *Logger
}

// ServerConfig holds configuration options for creating a STUN server.
type ServerConfig struct {
	// Addr is the IP address to bind to (e.g., "127.0.0.1", "0.0.0.0")
	Addr string
	// Port is the port number to listen on (e.g., "3478")
	Port string
	// Timeout is the connection timeout duration
	Timeout time.Duration
	// Logger is the logger instance to use for logging
	Logger *Logger
}

// NewServer creates a new STUN server with the specified configuration.
// If no logger is provided, a default logger will be used.
//
// Example:
//
//	server := stun.NewServer(stun.ServerConfig{
//		Addr:    "127.0.0.1",
//		Port:    "3478",
//		Timeout: 30 * time.Second,
//		Logger:  stun.NewDefaultLogger(),
//	})
func NewServer(cfg ServerConfig) *Server {
	logger := cfg.Logger
	if logger == nil {
		logger = NewDefaultLogger()
	}

	return &Server{
		addr:    cfg.Addr,
		port:    cfg.Port,
		timeout: cfg.Timeout,
		logger:  logger,
	}
}

// Listen starts the STUN server and begins listening for incoming connections.
// This method blocks indefinitely until the server is stopped or an error occurs.
//
// The server will:
//   - Bind to the specified address and port
//   - Accept incoming UDP connections
//   - Process STUN binding requests
//   - Send appropriate responses with XOR-MAPPED-ADDRESS
//   - Log all activities and errors
//
// Returns:
//   - error: Any error that occurred during server startup or operation
//
// Example:
//
//	server := stun.NewServer(stun.ServerConfig{
//		Addr: "127.0.0.1",
//		Port: "3478",
//	})
//	if err := server.Listen(); err != nil {
//		log.Fatal(err)
//	}
func (s *Server) Listen() error {
	addr := net.JoinHostPort(s.addr, s.port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		s.logger.LogError("Failed to resolve UDP address", err, map[string]interface{}{
			"address": addr,
		})
		return err
	}

	s.logger.Info("STUN server starting", map[string]interface{}{
		"address": addr,
		"timeout": s.timeout.String(),
	})

	conn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		s.logger.LogError("Failed to listen on UDP address", err, map[string]interface{}{
			"address": addr,
		})
		return err
	}

	defer conn.Close()

	s.logger.LogConnection(conn.LocalAddr().String(), "", "stun_server")

	for {
		s.HandleUDPConn(conn)
	}
}

// HandleUDPConn processes a single UDP connection and handles STUN requests.
// This method is called for each incoming UDP packet and performs:
//   - Reading the UDP packet
//   - Parsing the STUN message
//   - Validating the message format
//   - Generating the XOR-MAPPED-ADDRESS response
//   - Sending the response back to the client
//
// The method includes comprehensive error handling and logging for debugging
// and monitoring purposes.
func (s *Server) HandleUDPConn(con *net.UDPConn) {
	buff := make([]byte, 1024)
	n, remoteAddr, err := con.ReadFromUDP(buff)
	if err != nil {
		s.logger.LogError("Failed to read from UDP connection", err, map[string]interface{}{
			"remote_addr": remoteAddr.String(),
		})
		return
	}

	s.logger.Debug("Received UDP packet", map[string]interface{}{
		"remote_addr": remoteAddr.String(),
		"bytes_read":  n,
		"local_addr":  con.LocalAddr().String(),
	})

	packet, err := NewPacket(con, buff[:n], remoteAddr)
	if err != nil {
		s.logger.LogError("Failed to create packet from UDP data", err, map[string]interface{}{
			"remote_addr": remoteAddr.String(),
			"bytes_read":  n,
		})
		return
	}

	// Log the incoming request
	s.logger.LogRequest(remoteAddr.String(), packet.message.Header.Type, packet.message.Header.TransactionID)

	trID := packet.message.Header.TransactionID

	xorAddr, err := serializeAddr(XorMappedAddr{
		Family: IPV4,
		IP:     packet.remoteIP,
		Port:   packet.remotePort,
	}, trID)
	if err != nil {
		s.logger.LogError("Failed to serialize XOR mapped address", err, map[string]interface{}{
			"remote_addr":    remoteAddr.String(),
			"transaction_id": trID,
		})
		return
	}

	xorAttr := Attribute{
		Length:       XORMappedAddressLength,
		Type:         XORMappedAddress,
		PaddedLength: XORMappedAddressLength,
		Value:        xorAddr,
	}

	msg := Message{
		Header: Header{
			Type:          BindingResponse,
			Length:        XORMappedAddressLength + 4,
			TransactionID: trID,
			MagicCookie:   magicCookie,
		},
		Attributes: []Attribute{xorAttr},
	}
	content := msg.Encode()

	// Create XOR mapped address for logging
	xorMappedAddr := &XorMappedAddr{
		Family: IPV4,
		IP:     packet.remoteIP,
		Port:   packet.remotePort,
	}

	// Log the response being sent
	s.logger.LogResponse(remoteAddr.String(), msg.Header.Type, trID, xorMappedAddr)

	n, err = packet.Write(content, remoteAddr)
	if err != nil {
		s.logger.LogError("Failed to write response", err, map[string]interface{}{
			"remote_addr":    remoteAddr.String(),
			"transaction_id": trID,
			"bytes_written":  n,
		})
		return
	}

	s.logger.Debug("Response sent successfully", map[string]interface{}{
		"remote_addr":   remoteAddr.String(),
		"bytes_written": n,
	})
}

// Shutdown gracefully shuts down the STUN server.
// This method logs the shutdown event and can be extended to perform
// cleanup operations if needed.
//
// Returns:
//   - error: Any error that occurred during shutdown
func (s *Server) Shutdown() error {
	s.logger.LogShutdown("stun_server", 0)
	return nil
}
