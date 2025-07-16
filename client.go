package stun

import (
	"net"
)

// Client represents a STUN client that can send binding requests to STUN servers
// and receive responses containing the client's public IP address and port.
//
// The client handles the complete STUN transaction lifecycle:
//   - Resolving server addresses
//   - Creating and sending binding requests
//   - Receiving and parsing responses
//   - Extracting XOR-MAPPED-ADDRESS information
//
// Example:
//
//	client := stun.NewClient("stun.l.google.com:19302")
//	msg, err := client.Dial(&stun.Message{
//		Header: stun.Header{Type: stun.BindingRequest},
//	})
type Client struct {
	ServerAddr string
	logger     *Logger
}

// NewClient creates a new STUN client with the specified server address.
// The server address should be in the format "host:port".
//
// Example:
//
//	client := stun.NewClient("stun.l.google.com:19302")
func NewClient(addr string) *Client {
	return &Client{
		ServerAddr: addr,
		logger:     NewDefaultLogger(),
	}
}

// NewClientWithLogger creates a new STUN client with a custom logger.
// This allows for fine-grained control over logging behavior.
//
// Example:
//
//	logger := stun.NewLogger(stun.LoggerConfig{
//		Level:  stun.DebugLevel,
//		Format: "json",
//	})
//	client := stun.NewClientWithLogger("stun.l.google.com:19302", logger)
func NewClientWithLogger(addr string, logger *Logger) *Client {
	return &Client{
		ServerAddr: addr,
		logger:     logger,
	}
}

// Dial sends a STUN binding request to the server and returns the response.
// The method performs the complete STUN transaction:
//   - Resolves the server address
//   - Creates a UDP connection
//   - Sends the binding request
//   - Receives and parses the response
//   - Returns the parsed message
//
// The input message should have at least the Header.Type field set to BindingRequest.
// The method will automatically set the MagicCookie, Length, and TransactionID fields.
//
// Returns:
//   - *Message: The parsed STUN response message
//   - error: Any error that occurred during the transaction
//
// Example:
//
//	msg, err := client.Dial(&stun.Message{
//		Header: stun.Header{Type: stun.BindingRequest},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	xorAddr, err := msg.GetXorAddr()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Public IP: %s:%d\n", xorAddr.IP, xorAddr.Port)
func (client *Client) Dial(m *Message) (*Message, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", client.ServerAddr)
	if err != nil {
		client.logger.LogError("Failed to resolve server address", err, map[string]interface{}{
			"server_addr": client.ServerAddr,
		})
		return nil, err
	}

	m.Header.MagicCookie = magicCookie
	m.Header.Length = uint16(len(m.Attributes))
	m.Header.TransactionID = [12]byte(randomTransactionID())

	// Log the request being sent
	client.logger.LogClientRequest(client.ServerAddr, m.Header.Type, m.Header.TransactionID)

	encodedHeader := m.Header.Encode()

	c, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		client.logger.LogError("Failed to dial UDP connection", err, map[string]interface{}{
			"server_addr": client.ServerAddr,
		})
		return nil, err
	}
	defer c.Close()

	client.logger.LogConnection(c.LocalAddr().String(), udpAddr.String(), "stun_client")

	_, err = c.Write(encodedHeader)
	if err != nil {
		client.logger.LogError("Failed to write request to server", err, map[string]interface{}{
			"server_addr":    client.ServerAddr,
			"transaction_id": m.Header.TransactionID,
		})
		return nil, err
	}

	buff := make([]byte, 2048)
	_, _, err = c.ReadFromUDP(buff)
	if err != nil {
		client.logger.LogError("Failed to read response from server", err, map[string]interface{}{
			"server_addr":    client.ServerAddr,
			"transaction_id": m.Header.TransactionID,
		})
		return nil, err
	}

	msg, err := NewMessage(buff)
	if err != nil {
		client.logger.LogError("Failed to parse response message", err, map[string]interface{}{
			"server_addr":    client.ServerAddr,
			"transaction_id": m.Header.TransactionID,
		})
		return nil, err
	}

	// Get XOR mapped address for logging
	xorAddr, _ := msg.GetXorAddr()
	client.logger.LogClientResponse(client.ServerAddr, msg.Header.Type, xorAddr)

	return msg, nil
}
